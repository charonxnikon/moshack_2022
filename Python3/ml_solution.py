import pandas as pd
import numpy as np

from sklearn.neighbors import NearestNeighbors
from corrections import PullFlats
from corrections import adjustments_all as adjustments_default
import psycopg2
import geopy.distance
import argparse
import typing as tp

columns_to_compare = ["type", "material", "rooms", "height"]

columns_db_apartments = ["id", "user_id", "address", "rooms", "type",
                         "height", "material", "floor", "area",
                         "kitchen", "balcony", "metro", "condition",
                         "latitude", "longitude", "total_price", "price_m2"]

columns_user_apartments = ["id", "user_id", "address", "rooms", "type",
                           "height", "material", "floor", "area",
                           "kitchen", "balcony", "metro", "condition",
                           "latitude", "longitude", "total_price", "price_m2"]


def get_coords(address):
    from geopy.geocoders import Nominatim

    geolocator = Nominatim(user_agent="smy-application")
    location = geolocator.geocode(address)
    latitude_x = location.latitude
    longitude_y = location.longitude
    return latitude_x, longitude_y


def get_coords(address: str) -> tp.Tuple[int, int]:
    from geopy.geocoders import Nominatim

    geolocator = Nominatim(user_agent="smy-application")
    location = geolocator.geocode(address)
    print('location: ', location)
    if location is None:
        return 45, 42
    latitude_x = location.latitude
    longitude_y = location.longitude
    return latitude_x, longitude_y


def update_coords_user_apartments(user_id: int) -> None:
    columns = columns_user_apartments
    sql3 = \
    f"""
    SELECT * 
    FROM user_apartments 
    WHERE user_id = {user_id} 
    """
    cur.execute(sql3)
    pull = cur.fetchall()
    print('pull ', pull)
    pull = pd.DataFrame(np.array(pull).reshape(-1, len(columns)),
                        columns=columns)
    update_rows = []
    for index, row in pull.iterrows():
        if int(float(row.latitude)) == 0:
            lat, lon = get_coords(row.address)
            row.latitude = lat
            row.longitude = lon
            update_rows.append(row)

    for row in update_rows:
        print(f'update {row.id} row! \n')
        sql4 = \
            f"""
        UPDATE user_apartments
        SET latitude = {row.latitude},
            longitude = {row.longitude}
        WHERE id = {row.id}
        """
        cur.execute(sql4)
    conn.commit()


type2number = {
    'новостройка': 0,
    'современное жилье': 1,
    'старый жилой фонд': 2
}

condition2number = {
    'без отделки': 0,
    'муниципальный ремонт': 1,
    'современная отделка': 2
}

material2number = {
    'кирпич': 0,
    'панель': 1,
    'монолит': 2
}

pull_flats = PullFlats(adjustments_default)


def rooms2number(rooms):
    if rooms == 'студия':
        return 1.5
    else:
        return str(rooms)


def get_args():
    parser = argparse.ArgumentParser()

    parser.add_argument('-p', '--port', default=5432)
    parser.add_argument('-db', '--dbname', default="moshack")
    parser.add_argument('-u', '--user', default="postgres")
    parser.add_argument('-pw', '--password', default="777777")

    args = parser.parse_args()

    return args


args = get_args()

conn = psycopg2.connect(host="localhost", port=args.port,
                        dbname=args.dbname,
                        user=args.user, password=args.password)
cur = conn.cursor()


def calc_distance(coords_1, coords_2):
    return geopy.distance.geodesic(coords_1, coords_2).km


def get_flats_from_radius(df, max_dist, expert_flat):
    lst_nearest = []
    lst_idxs = []
    lst_dists = []
    for idx, row in df.iterrows():
        dist = calc_distance((row["latitude"], row["longitude"]),
                             (expert_flat.latitude, expert_flat.longitude))
        if 0 <= dist < max_dist:
            lst_nearest.append(row)
            lst_idxs.append(idx)
            lst_dists.append(dist)

    return lst_idxs, lst_dists, lst_nearest


def my_dist(first_flat, second_flat):
    """

    Parameters
    ----------
    expert_flat with None price

    Returns
    -------

    """
    first_flat = pd.DataFrame(first_flat.reshape(1, 5),
                              columns=columns_to_compare + ["dist"])
    second_flat = pd.DataFrame(second_flat.reshape(1, 5),
                               columns=columns_to_compare + ["dist"])

    #    total_price, expert_price_sq_meter, percent_corrects, counts_carefully = pull_flats.calculate_weights(first_flat, second_flat)

    return 1000 * np.abs(
        first_flat.loc[0, "type"] - first_flat.loc[0, "type"]) + \
           1000 * np.abs(
        first_flat.loc[0, "rooms"] - second_flat.loc[0, "rooms"]) + \
           1000 * np.abs(
        first_flat.loc[0, "material"] - second_flat.loc[0, "material"]) + \
           1000 * np.abs(
        first_flat.loc[0, "height"] - second_flat.loc[0, "height"])


def get_neighbors(id_expert_flat: int):
    # Open a cursor to perform database operations

    sql3 = f"""
    SELECT *
    FROM user_apartments
    WHERE id = {id_expert_flat}
    """
    cur.execute(sql3)
    fetchall = cur.fetchall()
    if not fetchall:
        print("don't id in db")
        raise ValueError("Don't Id in db")
    tuple_flat = fetchall[0]
    expert_flat = pd.DataFrame(
        np.array(tuple_flat).reshape(1, len(columns_db_apartments)),
        columns=columns_db_apartments)

    print('expert_flat', expert_flat)
    sql4 = f"""
    SELECT * 
    FROM db_apartments
    WHERE (latitude BETWEEN {float(expert_flat.latitude) - 0.01} AND {float(expert_flat.latitude) + 0.01}) AND 
        (longitude BETWEEN {float(expert_flat.longitude) - 0.02} AND {float(expert_flat.longitude) + 0.02})
    """

    cur.execute(sql4)
    record = cur.fetchall()
    print('record', record)

    df_nearest = pd.DataFrame(record, columns=columns_db_apartments)
    df_nearest[["latitude", "longitude"]] = df_nearest[
        ["latitude", "longitude"]].astype('float')
    expert_flat[["latitude", "longitude"]] = expert_flat[
        ["latitude", "longitude"]].astype('float')

    return df_nearest, expert_flat.iloc[0, :]


def get_price(df, idxs):
    df = df.set_index("id")
    price = np.mean(df.loc[idxs, "total_price"].values)

    return price


def get_analogs_flat_idxs(id_expert_flat: int):
    try:
        df, expert_flat = get_neighbors(id_expert_flat)
    except:
        print('error')
        return [], None, None
    print('df', df)
    print('exp', expert_flat)
    max_dist = 1
    lst_idxs, lst_dists, lst_nearest = get_flats_from_radius(df, max_dist,
                                                             expert_flat)
    print('lst_idxs', lst_idxs)
    if not lst_nearest:
        idxs = []
        print('not lst nearest')
        return idxs, None, None

    df2 = pd.DataFrame(lst_nearest)[columns_to_compare]
    #    df2['area'] = df2['area'].apply(lambda x: float(x))
    #    df2["balcony"] = df2["balcony"].apply(lambda x: 1 if x == "Да" else 0)
    df2["type"] = df2["type"].apply(lambda x: type2number[x.lower()])
    df2["material"] = df2["material"].apply(
        lambda x: material2number[x.lower()])
    df2["rooms"] = df2["rooms"].apply(lambda x: float(rooms2number(x.lower())))
    df2["height"] = df2["height"].apply(lambda x: float(x))

    expert_flat["type"] = type2number[expert_flat["type"].lower()]
    expert_flat["material"] = material2number[expert_flat["material"].lower()]
    expert_flat["rooms"] = float(rooms2number(expert_flat["rooms"].lower()))
    expert_flat["height"] = float(expert_flat["height"])

    expert_flat.balcony = 1 if expert_flat.balcony == "Да" else 0
    expert_flat.area = float(expert_flat.area)
    df2["dist"] = lst_dists
    vals = df2.values
    print('vals: ', vals)
    print('expert_flat: ', expert_flat.values)
    if len(vals) < 5:
        idxs = df.iloc[lst_idxs, 0].values
        price = get_price(df, idxs)

        print('len(vals < 5')
        return idxs, price / expert_flat.area, price

    expert_value = np.array([
        float(expert_flat.type), float(expert_flat.material),
        float(expert_flat.rooms), float(expert_flat.height), 0
    ]).reshape(1, -1)
    print('expert_value: ', expert_value)
    print('vals.shape: ', vals.shape)
    nbrs = NearestNeighbors(n_neighbors=5, algorithm='brute', metric=my_dist)
    nbrs.fit(vals)
    dists, idxs = nbrs.kneighbors(expert_value)
    idxs_lst = idxs.ravel().tolist()

    idxs = df.loc[idxs_lst, "id"].values
    price = get_price(df, idxs)
    print('idxs', idxs)

    return idxs, price / expert_flat.area, price,
