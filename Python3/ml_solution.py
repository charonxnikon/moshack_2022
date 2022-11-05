import pandas as pd
import numpy as np

from sklearn.neighbors import NearestNeighbors
import psycopg2
import geopy.distance
import argparse


columns_db_apartments = ["id", "user_id", "address", "rooms", "type",
                         "height", "material", "floor", "area",
                         "kitchen", "balcony", "metro", "condition",
                         "latitude", "longitude", "total_price", "price_m2"]


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
        if 0 <= dist < max_dist and row.id != int(expert_flat.id):
            lst_nearest.append(row)
            lst_idxs.append(idx)
            lst_dists.append(dist)

    return lst_idxs, lst_dists, lst_nearest

def my_dist(x, y):
    return abs(x[0] - y[0]) + 15 * abs(x[1] - y[1]) + 25 * abs(x[2] - y[2])


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
        raise ValueError("Don't Id in db")
    tuple_flat = fetchall[0]
    expert_flat = pd.DataFrame(np.array(tuple_flat).reshape(1, len(columns_db_apartments)), columns=columns_db_apartments)

    sql4 = f"""
    SELECT * 
    FROM db_apartments
    WHERE (latitude BETWEEN {float(expert_flat.latitude) - 0.01} AND {float(expert_flat.latitude) + 0.01}) AND 
        (longitude BETWEEN {float(expert_flat.longitude) - 0.02} AND {float(expert_flat.longitude) + 0.02})
    """

    cur.execute(sql4)
    record = cur.fetchall()

    df_nearest = pd.DataFrame(record, columns=columns_db_apartments)
    df_nearest[["latitude", "longitude"]] = df_nearest[["latitude", "longitude"]].astype('float')
    expert_flat[["latitude", "longitude"]] = expert_flat[["latitude", "longitude"]].astype('float')

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
        return
    max_dist = 1
    lst_idxs, lst_dists, lst_nearest = get_flats_from_radius(df, max_dist,
                                                             expert_flat)
    if not lst_nearest:
        idxs = []
        return idxs, None

    df2 = pd.DataFrame(lst_nearest)[["area", 'balcony']]
    df2['area'] = df2['area'].apply(lambda x: float(x))
    df2["balcony"] = df2["balcony"].apply(lambda x: 1 if x == "Да" else 0)
    expert_flat.balcony = 1 if expert_flat.balcony == "Да" else 0
    expert_flat.area = float(expert_flat.area)
    df2["dist"] = lst_dists
    vals = df2.values
    if len(vals) < 5:
        idxs = df.iloc[lst_idxs, 0].values
        price = get_price(df, idxs)

        return idxs, price

    expert_value = np.array(
        [float(expert_flat.area), float(expert_flat.balcony), 0]
    ).reshape(1, -1)
    nbrs = NearestNeighbors(n_neighbors=5, algorithm='brute', metric=my_dist)
    nbrs.fit(vals)
    #    print(expert_value)
    #    print(vals)
    dists, idxs = nbrs.kneighbors(expert_value)
    idxs_lst = idxs.ravel().tolist()

    idxs = df.loc[idxs_lst, "id"].values
    price = get_price(df, idxs)

    return idxs, price / expert_flat.area, price,


