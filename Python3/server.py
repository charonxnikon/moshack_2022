from corrections import recalculate_price_expert
from jsonrpcserver import Success, method, serve, Result
import json
import typing as tp
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
        if 0 <= dist < max_dist:
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



@method
def get_analogs(id_flat: int) -> Result:
#    try:
    idxs, price, total_price = get_analogs_flat_idxs(id_flat)
    idxs_new = list(map(float, idxs))

    return Success(json.dumps({"Analogs": idxs_new,
                               "PriceM2": price, "TotalPrice": total_price}))

from corrections import type_adjustments, data_adjustments
from corrections import adjustments as adjustments_not_inits
from copy import deepcopy
@method
def recalculate_price_expert_flat(expert_flat_id: int,
                                  analog_idxs: tp.List[int],
                                  needed_adjustments) -> Result:
    my_data_adjustments = deepcopy(data_adjustments)
    for adjustment in type_adjustments:
        if adjustment in needed_adjustments:
            my_data_adjustments[adjustment][0] = needed_adjustments[adjustment][0]
        else:
            del my_data_adjustments[adjustment]

    my_adjustments_all = []
    for idx, type_adj in enumerate(type_adjustments):
        if type_adj in needed_adjustments:
            my_adjustments_all.append(adjustments_not_inits[type_adj](*my_data_adjustments[type_adj]))

    price, total_price = \
        recalculate_price_expert(expert_flat_id, analogs_idxs=analog_idxs,
                                 adjustments=my_adjustments_all)
    return Success(json.dumps({"price": price, "totalPrice": total_price}))

def recalculate_price_expert_flat_my(expert_flat_id: int,
                                  analog_idxs: tp.List[int],
                                  needed_adjustments) -> tp.Any:
    my_data_adjustments = deepcopy(data_adjustments)
    for adjustment in type_adjustments:
        if adjustment in needed_adjustments:
            my_data_adjustments[adjustment][0] = needed_adjustments[adjustment][0]
        else:
            del my_data_adjustments[adjustment]

    my_adjustments_all = []
    for idx, type_adj in enumerate(type_adjustments):
        if type_adj in needed_adjustments:
            my_adjustments_all.append(adjustments_not_inits[type_adj](*my_data_adjustments[type_adj]))

    price, total_price = \
        recalculate_price_expert(expert_flat_id, analogs_idxs=analog_idxs,
                                 adjustments=my_adjustments_all)
    return {"price": price, "totalPrice": total_price}

def calculate_pull_one_expert(id_expert_flat: int,
                   idx_analogs: tp.List[int],
                   needed_adjustments):
    lst_price_total_price =[]
    for idx_analog in idx_analogs:
        dct = recalculate_price_expert_flat_my(idx_analog,
                                          [id_expert_flat],
                                          needed_adjustments)
        lst_price_total_price.append(list(dct.values()))

    return lst_price_total_price

@method
def calculate_pull(idxs_expert_flat: tp.List[int],
                   idx_analogs: tp.List[int],
                   needed_adjustments):
    lst_all = []
    for idx_expert_flat in idxs_expert_flat:
        result = calculate_pull_one_expert(idx_expert_flat, idx_analogs,
                                 needed_adjustments)
        lst_all.append(result)

    all_data = np.array(lst_all)
    final_result = np.mean(all_data, axis=0)
#    all_data = list(map())

    return Success(json.dumps({"all_prices": all_data.tolist(),
                   "final_price": final_result.tolist()}))


def tmp(idxs_expert_flat: tp.List[int],
                   idx_analogs: tp.List[int],
                   needed_adjustments):
    lst_all =[]
    for idx_expert_flat in idxs_expert_flat:
        result = calculate_pull_one_expert(idx_expert_flat, idx_analogs,
                                          needed_adjustments)
        lst_all.append(result)

    all_data = np.array(lst_all)
    final_result = np.mean(all_data, axis=0)
    print("IT'S HERE\n", json.dumps({"all_prices": all_data.tolist(),
                   "final_price": final_result.tolist()}))

    return lst_all


@method
def ping(expert_flat_id: int, dct):
    print(dct)
    print(expert_flat_id)

    return Success(1)


if __name__ == "__main__":
#    result = tmp([1, 2, 8], [3, 4, 5, 6, 7], {"tender": [-0.06]})
#    print(np.array(result))
#    print(np.array(result).shape)
#    print(np.mean(np.array(result), axis=0))
    #print('final: ', tmp([1, 2], [3, 4, 5, 6, 7], {"tender": [-0.06]}))

#    idxs, price, total_price = get_analogs_flat_idxs(1)
    try:
        serve()
    except:
        cur.close()
        conn.close()
        print('connection is closed')



