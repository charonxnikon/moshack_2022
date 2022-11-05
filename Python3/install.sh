sudo -u postgres psql < ../DataBaseScripts/SetDataBase.sql
source ./venv/bin/activate
python3 create_dumy_table_pull.py
python3 server.py
source ./venv/bin/deactivate
curl -X POST http://localhost:5000 -d '{"jsonrpc": "2.0", "method": "get_analogs", "id": 1, "id_flat":1}'

