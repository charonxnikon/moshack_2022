# moshack_2022

go mod init moshack_2022
go mod download
go mod verify
go mod tidy
go mod vendor

Чтобы запустить golang (из корневой директории проекта):
go run cmd/moshack_2022/main.go

Чтобы запустить питон (из Python3):
pip install -r requirements.txt
python3 create_dumy_table_pull.py -pw пароль
python3 server.py -pw пароль

Если нет python3 - тогда просто python:
python create_dumy_table_pull.py -pw пароль
python server.py -pw пароль

При необходимости серверу через командную строку можно не только пароль от бд передать,
но и порт БД (-p), имя пользователя (-u), название БД (-db).

При запуске golang пытается подключиться к серверу python,
если не получается - логгирует ошибку.

# go build  -o ./bin/crudapp ./cmd/crudapp
# go test -v -coverpkg=./... ./...

# go build -mod=vendor -o ./bin/myapp ./cmd/myapp
# go test -v -mod=vendor -coverpkg=./... ./...


ЧТОБЫ ПОСМОТРЕТЬ ПАРСИНГ ЕКСЕЛЬ ФАЙЛА: файл ексель в корневом каталоге с именем "test.xls" ОБЯЗАТЕЛЬНО XLS
 go run cmd/moshack_2022/*.go excel
ЧТОБЫ ЗАПУСТЬ СЕРВЕР:
 go run cmd/moshack_2022/*.go
 
 Я ЗАПУСКАЛ SQL СКРИПТ ТАК:
 	sudo cp DataBaseScripts/<имя скрипта(можно * чтобы скопировать сразу все)>.sql /var/lib/postgresql/SetDataBase.sql
	sudo -i -u postgres
	psql -U postgres -f <имя скрипта>.sql
