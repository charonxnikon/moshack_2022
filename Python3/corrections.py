import numpy as np
import pandas as pd
import psycopg2
import argparse
import typing as tp

# add enums
# покрыть тестами все
# изменил bin's для metro_dist
# kitchen_area was changed
# metro_dist was changed

columns = ["id", "user_id", "adress", "rooms", "type",
           "height", "material", "floor", "area",
           "kitchen", "balcony", "metro", "condition",
           "latitude", "longitude", "price", "price_m2"]

condition2code = {
    'Без отделки': 0,
    'Муниципальный ремонт': 1,
    'Улучшенный': 2
}

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


def get_idxs_from_table(table, idxs):
    flats = []
    for idx in idxs:
        sql3 = f"""
        SELECT *
        FROM {table} 
        WHERE Id = {idx}
        """
        cur.execute(sql3)
        fetchall = cur.fetchall()
        if not fetchall:
            raise ValueError("Don't Id in db")
        tuple_flat = fetchall[0]
        expert_flat = pd.DataFrame(
            np.array(tuple_flat).reshape(1, len(columns)), columns=columns)
        flats.append(expert_flat)

    df = pd.concat(flats).reset_index().set_index('id')
    df.balcony = df.balcony.apply(lambda x: 1 if x == 'Да' else 0)

    return df


class AdjustmentTender:
    def __init__(self, adjustment):
        self.adjustment = adjustment

    def calculate(self, expert, analog):
        return self.adjustment


class AdjustmentFloorFlat:
    def __init__(self, adjustment_matr):
        # suppose like in appendix 2
        self.adjustment_matr = adjustment_matr

    @staticmethod
    def _calculate_index(flat):
        if float(flat.height) == float(flat.floor):
            return 2

        if 1 < float(flat.floor) < float(flat.height):
            return 1

        if float(flat.floor) == 1:
            return 0

    def calculate(self, expert, analog):
        index_expert = self._calculate_index(expert)
        index_analog = self._calculate_index(analog)

        return self.adjustment_matr[index_expert][index_analog]


# FlatArea, KitchenArea, Balcony, MetroDist,

class AdjustmentGeneral:
    def __init__(self, adjustment_matr, bins, attr):
        self.bins = bins
        self.adjustment_matr = adjustment_matr
        self.attr = attr

    def _calculate_index(self, flat):
        for i in range(len(self.bins) - 1):
            if self.bins[i] <= float(flat[self.attr]) < self.bins[i + 1]:
                return i

        assert False, "Unreachable Code"

    def calculate(self, expert, analog):
        index_expert = self._calculate_index(expert)
        index_analog = self._calculate_index(analog)

        return self.adjustment_matr[index_expert][index_analog]


class AdjustmentRepair:
    def __init__(self, adjustment_matr):
        # suppose like in appendix 2
        self.adjustment_matr = adjustment_matr

    @staticmethod
    def _calculate_index(flat):
        return int(condition2code[flat.condition])

    def calculate(self, expert, analog):
        index_expert = self._calculate_index(expert)
        index_analog = self._calculate_index(analog)

        return self.adjustment_matr[index_expert][index_analog]\
               / float(analog.price_m2)


data_adjustments = {
    "tender": [-0.045],
    "floor": [
        np.array([
            [-0.000, -0.070, -0.031],
            [0.075, 0.000, 0.042],
            [0.032, -0.040, 0.000]
        ])
    ],
    "area": [
        np.array([
            [-0.00, 0.06, 0.14, 0.21, 0.28, 0.31],
            [-0.06, 0.00, 0.07, 0.14, 0.21, 0.24],
            [-0.12, -0.07, 0.00, 0.06, 0.13, 0.16],
            [-0.17, -0.12, -0.06, 0.00, 0.06, 0.09],
            [-0.22, -0.17, -0.11, -0.06, 0.00, 0.03],
            [-0.24, -0.19, -0.13, -0.08, -0.03, 0.00]]),
        np.array([
            0, 30, 50, 65, 90, 120, float('inf')
        ]),
        "area"
    ],
    "kitchen": [
        np.array([
            [+0.000, -0.029, -0.083],
            [+0.030, +0.000, -0.055],
            [+0.090, +0.058, +0.000]
        ]),
        np.array([
            #            0, 7, 10, 15, float('inf')
            0, 7, 10, float('inf')
        ]),
        "kitchen"
    ],
    "balcony": [
        np.array([
            [0.00, -0.05],
            [0.053, 0.00]
        ]),
        np.array([0, float('inf')]),
        "balcony"
    ],
    "metro": [
        np.array([
            [0.00, 0.07, 0.12, 0.17, 0.24, 0.29],
            [-0.07, 0.00, 0.04, 0.09, 0.15, 0.20],
            [-0.11, -0.04, 0.00, 0.05, 0.11, 0.15],
            [-0.15, -0.08, -0.05, 0.00, 0.06, 0.10],
            [-0.19, -0.13, -0.10, -0.06, 0.00, 0.04],
            [-0.22, -0.17, -0.13, -0.09, -0.04, 0.00]
        ]),
        #        np.array([0, 5, 10, 15, 30, 60, 90, float('inf')]),
        np.array([0, 5, 10, 15, 30, 60, float('inf')]),
        "metro"
    ],
    "condition": [
        np.array([
            [0, -13400, -20100],
            [13399, 0, -6700],
            [20099, 6700, 0]
        ]),
    ]
}

adjustments = {
    "tender": AdjustmentTender,
    "floor": AdjustmentFloorFlat,
    "area": AdjustmentGeneral,
    "kitchen": AdjustmentGeneral,
    "balcony": AdjustmentGeneral,
    "metro": AdjustmentGeneral,
    "condition": AdjustmentRepair
}

type_adjustments = [
    "tender",
    "floor",
    "area",
    "kitchen",
    "balcony",
    "metro",
    "condition"
]

adjustments_all = []
for type_adj in type_adjustments:
    #    print('type_adjustment initialized: ', type_adj)
    adjustments_all.append(adjustments[type_adj](*data_adjustments[type_adj]))


class PullFlats:
    def __init__(self, needed_adjustments):
        self.weights = None
        self.needed_adjustments = needed_adjustments

    def calculate_analog_price(self, expert, analog):
        """
        expert and analog is
        describe of flats
        :::
        returns
        price of expert by analog
        and
        """
        analog_price_sq_meter = float(analog.price_m2)
        expert_price_sq_meter = float(analog_price_sq_meter)
        percent_corrects = 0
        for i, adjustment in enumerate(self.needed_adjustments):
            print(type_adjustments[i], end=': ')
            percent = float(adjustment.calculate(expert, analog))
            print('percent: ', percent)
            expert_price_sq_meter *= (1 + percent)
            print(' : new_price: ', (1 + percent))
            print(' : new_price: ', expert_price_sq_meter)
            percent_corrects += abs(percent)

        total_price = expert_price_sq_meter * float(analog.area)

        return total_price, expert_price_sq_meter, percent_corrects

    def calculate_weights(self, percent_corrects):
        """
        :param percent_corrects:  is numpy array
        :return:
        """
        min_correct = np.min(percent_corrects)
        if min_correct == 0:
            min_idxs = np.argmin(percent_corrects)
            weights = np.zeros_like(percent_corrects)
            weights[min_idxs] = 1 / len(min_idxs)
            return weights
        inv_percent_corrects = min_correct / percent_corrects
        inv_sum = 1 / np.sum(inv_percent_corrects)
        return inv_percent_corrects * inv_sum

    def calculate_pull(self, id_expert_flat, pull_analogs):
        """
        calculate pull of pandas dataframe
        :param id_expert_flat: id expert_flat
        :param pull_analogs:  pd.DataFrame
        :return:
        """
        prices_total_analogs = []
        prices_sq_meter_analogs = []
        percent_corrects = []
        expert_flat = pull_analogs.iloc[str(id_expert_flat) == pull_analogs.index, :]
#        expert_flat = pull_analogs.loc[int(id_expert_flat), :]
        for index, row in pull_analogs.iterrows():
            if index != id_expert_flat:
                analog_flat = pull_analogs.loc[index, :]
                total_price, price, percent_correct = \
                    self.calculate_analog_price(expert_flat, analog_flat)
                prices_total_analogs.append(total_price)
                prices_sq_meter_analogs.append(price)
                percent_corrects.append(percent_correct)

        weights = self.calculate_weights(np.array(percent_corrects))
        print(weights)
        return np.dot(weights, prices_sq_meter_analogs), \
               np.dot(weights, prices_total_analogs)


data_flats = {
    "floor": [7, 3, 1, 4],
    "total_area": [85.0, 77.4, 84.0, 64.0],
    "kitchen_area": [15, 14, 12, 11.5],
    "balcony": [1, 1, 1, 1],
    "metrodist": [11, 10, 14, 11],
    "condition": [1, 2, 2, 2],
    #  other parameters without correction
    "max_floor": [22, 24, 18, 18],
    "NumberRooms": [2, 2, 2, 2],
    "price_all": [None, 28_750_000, 30_650_000, 26_500_000],
    "price_sq_meter": [None, 371_447, 364_881, 414_063]
}



def recalculate_price_expert(expert_idx: int,
                             analogs_idxs: tp.List[int],
                             adjustments):
    """
    :param expert_idx: id_expert  from table of apartments
    :param analogs_idxs: analogs list idxs from table of user_analogs
    :param adjustments:
    :return:
    """
    expert_flat = get_idxs_from_table('db_apartments', [expert_idx])
    analogs_flat = get_idxs_from_table('user_apartments', analogs_idxs)
    all_flats = pd.concat([expert_flat, analogs_flat])
    pull_flats_df = pd.DataFrame.from_dict(all_flats)
    pull_flats = PullFlats(needed_adjustments=adjustments)
    pull_flats_df = pull_flats_df.reset_index().set_index('id')

    return pull_flats.calculate_pull(1, pull_flats_df)


# pull_flats_df = pd.DataFrame.from_dict(data_flats)
# pull_flats = PullFlats(needed_adjustments=adjustments_all)
# print(pull_flats.calculate_pull(0, pull_flats_df))
# print(recalculate_price_expert(1, [3, 4, 5, 6], adjustments_all))
