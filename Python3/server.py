from corrections import recalculate_price_expert
from jsonrpcserver import Success, method, serve, Result
import json
import typing as tp
import numpy as np
from ml_solution import get_analogs_flat_idxs, conn, cur
from ml_solution import update_coords_user_apartments
from corrections import type_adjustments, data_adjustments
from corrections import adjustments as adjustments_not_inits
from corrections import adjustments_all as adjustments_default
from copy import deepcopy
from ml_solution import columns_user_apartments
import pandas as pd
import logging
logger = logging.getLogger('jsonrpcserver').setLevel(logging.DEBUG)
print(logger)


user2expert_flats = {}

def recalculate_price_expert_flat_my(expert_flat_id: int,
                                     analog_idxs: tp.List[int],
                                     needed_adjustments: tp.Dict[str, tp.List[tp.Any]],
                                     table_expert: str,
                                     table_analogs: str) -> tp.Any:
    my_data_adjustments = deepcopy(data_adjustments)
    for adjustment in type_adjustments:
        if adjustment in needed_adjustments:
            my_data_adjustments[adjustment][0] = \
                needed_adjustments[adjustment][0]
        else:
            del my_data_adjustments[adjustment]

    my_adjustments_all = []
    for idx, type_adj in enumerate(type_adjustments):
        if type_adj in needed_adjustments:
            my_adjustments_all.append(
                adjustments_not_inits[type_adj](*my_data_adjustments[type_adj])
            )

    print('my_adjustments_all: ', my_adjustments_all)
    print('type_adjustments: ', type_adjustments)
    print('needed_adjustments', needed_adjustments)
    price, total_price = \
        recalculate_price_expert(expert_flat_id, analog_idxs,
                                 my_adjustments_all,
                                 table_expert, table_analogs)
    return {"Price": price, "TotalPrice": total_price}

@method
def update_pull(user_id: int) -> Result:
    update_coords_user_apartments(user_id)

    return Success('done')

def update_price(id_flat: int, price_m2: float, total_price: float):
    sql_update = \
    f"""
    UPDATE user_apartments
    SET price_m2 = {price_m2},
        total_price = {total_price}
    WHERE id = {id_flat}
    """
    cur.execute(sql_update)
    conn.commit()

@method
def get_analogs(id_flat: int) -> Result:
    result = get_analogs_tmp(id_flat)
    print('write in user_apartments', id_flat, result["PriceM2"], result["TotalPrice"])
    update_price(id_flat, result["PriceM2"], result["TotalPrice"])

    return Success(json.dumps(result))
    

def get_analogs_tmp(id_flat: int) -> tp.Any:
#    update_coords_user_apartments(1)
    idxs, price, total_price = get_analogs_flat_idxs(id_flat)
    if idxs == []:
        return {"Analogs": [], "PriceM2": -1.0, "TotalPrice": -1.0}
    idxs_new = list(map(int, idxs))
    print('id_flat: ', id_flat)
    print('idxs : ', idxs)
    res = recalculate_price_expert_flat_my(id_flat, idxs, data_adjustments,
                                           'user_apartments', 'db_apartments')
    print('get_analogs_tmp: ', )
    price_m2 = res["Price"]
    total_price = res["TotalPrice"]

    return {"Analogs": idxs_new, "PriceM2": price_m2, "TotalPrice": total_price}


@method
def recalculate_price_expert_flat(id: int,
                                  analogs: tp.List[int],
                                  tender: tp.List[tp.List[int]],
                                  floor: tp.List[tp.List[int]],
                                  area: tp.List[tp.List[int]],
                                  kitchen: tp.List[tp.List[int]],
                                  balcony: tp.List[tp.List[int]],
                                  metro: tp.List[tp.List[int]],
                                  condition: tp.List[tp.List[int]]) -> Result:

    expert_flat_id = id
    analog_idxs = analogs

    needed_adjustments = {
        "tender": [np.array((np.array(tender) / 100).tolist())],
        "floor": [np.array((np.array(floor) / 100).tolist())],
        "area": [np.array((np.array(area) / 100).tolist())],
        "kitchen": [np.array((np.array(kitchen) / 100).tolist())],
        "balcony": [np.array((np.array(balcony) / 100).tolist())],
        "metro": [np.array((np.array(metro) / 100).tolist())],
        "condition": [np.array((np.array(condition)).tolist())]
    }
    result = recalculate_price_expert_flat_my(expert_flat_id, analog_idxs,
                                              needed_adjustments,
                                              'user_apartments',
                                              'db_apartments')
    if result["Price"] < 0:
        exit(1)

    update_price(expert_flat_id, result["Price"], result["TotalPrice"])

    return Success(json.dumps(result))


def calculate_pull_one_expert(id_expert_flat: int,
                              idx_analogs: tp.List[int],
                              needed_adjustments) -> Result:

    lst_price_total_price = []
    for idx_analog in idx_analogs:
        dct = recalculate_price_expert_flat_my(idx_analog,
                                               [id_expert_flat],
                                               needed_adjustments,
                                               'user_apartments',
                                               'user_apartments')
        lst_price_total_price.append(list(dct.values()))

    return lst_price_total_price

def get_idxs_pull(user_id: int, table: str) -> tp.List[int]:
    columns = columns_user_apartments
    sql_get_idxs = f"""
    SELECT *
    FROM {table}
    WHERE user_id = {user_id}
    """
    cur.execute(sql_get_idxs)
    pull = cur.fetchall()
    pull = pd.DataFrame(np.array(pull).reshape(-1, len(columns)),
                        columns=columns)
    return pull.id.values.tolist()

def update_prices(idx_analogs: tp.List[int], final_prices: tp.List[tp.List[int]],
                  table_name: str) -> None:
    for idx, idx_analog in enumerate(idx_analogs):
        sql_update = \
        f"""
        UPDATE {table_name}
        
        SET price_m2 = {final_prices[idx][0]},
            total_price = {final_prices[idx][1]}
        WHERE id = {idx_analog};
        """
        cur.execute(sql_update)

    conn.commit()



def calculate_pull_my(idxs_expert_flat: tp.List[int],
                      user_id: int,
                      needed_adjustments) -> Result:
    lst_all = []
    all_analogs = get_idxs_pull(user_id, 'user_apartments')
    print('all_analogs: ', all_analogs)
    idx_analogs = []
    for idx in all_analogs:
        print(idx, idxs_expert_flat)
        if int(idx) not in idxs_expert_flat:
            idx_analogs.append(idx)
#    idx_analogs = list(set(all_analogs) - set(idxs_expert_flat))
    print('idx_analogs: ', idx_analogs)

    for idx_expert_flat in idxs_expert_flat:
        result = calculate_pull_one_expert(idx_expert_flat, idx_analogs,
                                           needed_adjustments)
        print(result)
        lst_all.append(result)

    all_data = np.array(lst_all)
    final_result = np.mean(all_data, axis=0)
    update_prices(idx_analogs, final_result.tolist(), 'user_apartments')

    return 1


@method
def calculate_pull(Samples: tp.List[int],
                   user_id: int,
                      tender: tp.List[tp.List[int]],
                      floor: tp.List[tp.List[int]],
                      area: tp.List[tp.List[int]],
                      kitchen: tp.List[tp.List[int]],
                      balcony: tp.List[tp.List[int]],
                      metro: tp.List[tp.List[int]],
                      condition: tp.List[tp.List[int]]) -> Result:
    needed_adjustments = {
        "tender": [np.array((np.array(tender) / 100).tolist())],
        "floor": [np.array((np.array(floor) / 100).tolist())],
        "area": [np.array((np.array(area) / 100).tolist())],
        "kitchen": [np.array((np.array(kitchen) / 100).tolist())],
        "balcony": [np.array((np.array(balcony) / 100).tolist())],
        "metro": [np.array((np.array(metro) / 100).tolist())],
        "condition": [np.array((np.array(condition)).tolist())]
    }
    result = calculate_pull_my(Samples, user_id,
                               needed_adjustments)
    
    return Success(json.dumps(result))


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
    # print('final: ', tmp([1, 2], [3, 4, 5, 6, 7], {"tender": [-0.06]}))

    #    idxs, price, total_price = get_analogs_flat_idxs(1)
    print(get_analogs_tmp(1))
#    print(calculate_pull_my([1, 2, 8], 1, {"tender": [-0.06]}))
#    print(recalculate_price_expert_flat(1, [2, 3, 4], {"tender": [-0.06]}))
#    #print(update_coords_user_apartments(1))


    try:
        serve()
    except BaseException:
        cur.close()
        conn.close()
        print('connection is closed')
