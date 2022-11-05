sudo -u postgres psql < ../DataBaseScripts/SetDataBase.sql 
python3 create_dumy_table_pull.py
python3 server.py
curl -X POST http://localhost:5000 -d '{"jsonrpc": "2.0", "method": "get_analogs", "params": [1], "id": 1}'

