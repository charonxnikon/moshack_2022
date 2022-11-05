from corrections import recalculate_price_expert
from jsonrpcserver import Success, method, serve, Result
import json
import typing as tp
import numpy as np

from ml_solution import get_analogs_flat_idxs
from corrections import type_adjustments, data_adjustments
from corrections import adjustments as adjustments_not_inits
from copy import deepcopy
@method
def get_analogs(id_flat: int) -> Result:
#    try:
    idxs, price, total_price = get_analogs_flat_idxs(id_flat)
    idxs_new = list(map(float, idxs))

    return Success(json.dumps({"Analogs": idxs_new,
                               "PriceM2": price, "TotalPrice": total_price}))
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
    return Success(json.dumps({"Price": price, "TotalPrice": total_price}))

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
    return {"Price": price, "TotalPrice": total_price}

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

    return Success(json.dumps({"AllPrices": all_data.tolist(),
                   "FinalPrice": final_result.tolist()}))


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
#    print("IT'S HERE\n", json.dumps({"all_prices": all_data.tolist(),
#                   "final_price": final_result.tolist()}))

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



