sudo -u postgres psql < ../DataBaseScripts/SetDataBase.sql 
python3 create_dumy_table_pull.py
python3 server.py
curl -X POST http://localhost:5000 -d '{"jsonrpc": "2.0", "method": "get_analogs", "params": [1], "id": 1}'
curl -X POST http://localhost:5000 -d '{"jsonrpc": "2.0", "method": "recalculate_price_expert_flat", "params": {"expert_flat_id": 1, "analog_idxs": [2, 3, 4], "needed_adjustments":{"tender":[-0.06], "id": 1}'

