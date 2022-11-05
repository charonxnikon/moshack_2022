import argparse

parser = argparse.ArgumentParser()

parser.add_argument('-p', '--port')           # positional argument
parser.add_argument('-db', '--dbname')      # option that takes a value
parser.add_argument('-u', '--user')  # on/off flag
parser.add_argument('-pw', '--password')  # on/off flag

args = parser.parse_args()
print(args.port, args.dbname, args.user, args.password)

