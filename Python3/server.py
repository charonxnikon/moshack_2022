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


def recalculate_price_expert_flat_my(expert_flat_id: int,
                                     analog_idxs: tp.List[int],
                                     needed_adjustments,
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

    price, total_price = \
        recalculate_price_expert(expert_flat_id, analog_idxs,
                                 my_adjustments_all,
                                 table_expert, table_analogs)
    return {"Price": price, "TotalPrice": total_price}

@method
def update_pull(user_id: int) -> Result:
    update_coords_user_apartments(user_id)

    return Success('done')


@method
def get_analogs(id_flat: int) -> Result:
    result = get_analogs_tmp(id_flat)
    return Success(json.dumps(result))
    

def get_analogs_tmp(id_flat: int) -> tp.Any:
    idxs, price, total_price = get_analogs_flat_idxs(id_flat)
    if idxs == []:
        return {"Analogs": [], "PriceM2": -1.0, "TotalPrice": -1.0}
    idxs_new = list(map(int, idxs))
    print('id_flat: ', id_flat)
    print('idxs : ', idxs)
    res = recalculate_price_expert_flat_my(id_flat, idxs, adjustments_default,
                                           'user_apartments', 'db_apartments')
    price_m2 = res["Price"]
    total_price = res["TotalPrice"]

    return {"Analogs": idxs_new, "PriceM2": price_m2, "TotalPrice": total_price}


@method
def recalculate_price_expert_flat(expert_flat_id: int,
                                  analog_idxs: tp.List[int],
                                  needed_adjustments) -> Result:
    result = recalculate_price_expert_flat_my(expert_flat_id, analog_idxs,
                                              needed_adjustments,
                                              'user_apartments',
                                              'db_apartments')
    return Success(json.dumps(result))


def calculate_pull_one_expert(id_expert_flat: int,
                              idx_analogs: tp.List[int],
                              needed_adjustments):
    lst_price_total_price = []
    for idx_analog in idx_analogs:
        dct = recalculate_price_expert_flat_my(idx_analog,
                                               [id_expert_flat],
                                               needed_adjustments,
                                               'user_apartments',
                                               'user_apartments')
        lst_price_total_price.append(list(dct.values()))

    return lst_price_total_price


def calculate_pull_my(idxs_expert_flat: tp.List[int],
                      idx_analogs: tp.List[int],
                      needed_adjustments):
    lst_all = []
    for idx_expert_flat in idxs_expert_flat:
        result = calculate_pull_one_expert(idx_expert_flat, idx_analogs,
                                           needed_adjustments)
        lst_all.append(result)

    all_data = np.array(lst_all)
    final_result = np.mean(all_data, axis=0)

    return json.dumps({"AllPrices": all_data.tolist(),
                       "FinalPrice": final_result.tolist()})


@method
def calculate_pull(idxs_expert_flat: tp.List[int],
                   idx_analogs: tp.List[int],
                   needed_adjustments):
    result = calculate_pull_my(idxs_expert_flat, idx_analogs,
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
    print(get_analogs_tmp(7))
    print(calculate_pull_my([1, 2, 8], [3, 4, 5, 6, 7], {"tender": [-0.06]}))
    print(recalculate_price_expert_flat(1, [2, 3, 4], {"tender": [-0.06]}))
    #print(update_coords_user_apartments(1))


    try:
        serve()
    except BaseException:
        cur.close()
        conn.close()
        print('connection is closed')
