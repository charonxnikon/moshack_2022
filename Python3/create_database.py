import psycopg2
import argparse

def get_args():
    parser = argparse.ArgumentParser()

    path_default = '/home/pedashenko/Desktop/' \
                   'Hackathon/CIAN_project/apartments.csv'
    parser.add_argument('-p', '--port', default=5432)
    parser.add_argument('-db', '--dbname', default="test_erp")
    parser.add_argument('-u', '--user', default="postgres")
    parser.add_argument('-pw', '--password', default="777777")
    parser.add_argument('-ph', '--path', default=path_default)

    args = parser.parse_args()

    return args

args = get_args()

conn = psycopg2.connect(host="localhost", port=args.port,
                        dbname=args.dbname,
                        user=args.user, password=args.password)
cur = conn.cursor()
sql = """
CREATE TABLE apartments(Id int, Lat float, \
Lon float, District int, Small_district text, Rooms int, Price float, \
Totsp float, Livesp float, Kitsp float, Dist float, Stname text, \
Metrdist float, Walk float, Brick float, Tel int, Bal int, \
Floor float, Nfloors float, New int, Link text);
"""
cur.execute(sql)

sql2 = f"""
COPY apartments 
FROM {args.path} 
DELIMITER ',' 
CSV HEADER;
"""

cur.execute(sql2)

cur.close()
conn.commit()
conn.close()
