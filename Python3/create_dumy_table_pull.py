import numpy as np
import pandas as pd
import psycopg2
import argparse

# add enums
# покрыть тестами все
# изменил bin's для metro_dist
# kitchen_area was changed
# metro_dist was changed

columns_db_apartments = ["id", "user_id", "address", "rooms", "type",
           "height", "material", "floor", "area",
           "kitchen", "balcony", "metro", "condition",
           "latitude", "longitude", "total_price", "price_m2"]

columns_user_apartments = ["id", "user_id", "address", "rooms", "type",
           "height", "material", "floor", "area",
           "kitchen", "balcony", "metro", "condition",
           "latitude", "longitude", "total_price", "price_m2"]


def get_args():
    parser = argparse.ArgumentParser()

    parser.add_argument('-p', '--port', default=5432)
    parser.add_argument('-db', '--dbname', default="moshack")
    parser.add_argument('-u', '--user', default="postgres")
    parser.add_argument('-pw', '--password', default="777777")

    args_res = parser.parse_args()

    return args_res

args = get_args()

conn = psycopg2.connect(host="localhost", port=args.port,
                        dbname=args.dbname,
                        user=args.user, password=args.password)
cur = conn.cursor()

def fill_table(table, file):
    import os
    sql2 = f"""
    COPY {table} ({", ".join(columns_user_apartments[1:])})
    FROM '{os.getcwd() + '/' + file}'
    DELIMITER ','
    CSV HEADER;
    """
    print(os.getcwd() + '/' + file)
    cur.execute(sql2)

def get_idxs_from_table(table, idxs, columns):
    flats = []
    for idx in idxs:
        sql3 = f"""
        SELECT *
        FROM {table} 
        WHERE id = {idx}
        """
        cur.execute(sql3)
        fetchall = cur.fetchall()
        if not fetchall:
            raise ValueError("Don't Id in db")
        tuple_flat = fetchall[0]
        expert_flat = pd.DataFrame(
            np.array(tuple_flat).reshape(1, len(columns)), columns=columns)
        flats.append(expert_flat)

    return pd.concat(flats).reset_index().set_index('id')

fill_table(table='user_apartments',
           file='example_date.csv')
fill_table(table='db_apartments',
           file='example_db.csv')
conn.commit()
#df = get_idxs_from_table(table='user_apartments',
#                         idxs=[1, 2],
#                         columns=columns_user_apartments)
#
#print(df)
#
#df = get_idxs_from_table(table='db_apartments',
#                         idxs=[1, 2],
#                         columns=columns_user_apartments)
#print(df)

cur.close()
conn.commit()
conn.close()
