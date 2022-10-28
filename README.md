# moshack_2022

go mod init crudapp
# go mod init github.com/rvasily/crudapp
go build
go mod download
go mod verify
go mod tidy

go build  -o ./bin/crudapp ./cmd/crudapp
go test -v -coverpkg=./... ./...

go mod vendor
go build -mod=vendor -o ./bin/myapp ./cmd/myapp
go test -v -mod=vendor -coverpkg=./... ./...


ЧТОБЫ ПОСМОТРЕТЬ ПАРСИНГ ЕКСЕЛЬ ФАЙЛА: файл ексель в корневом каталоге с именем "test.xls" ОБЯЗАТЕЛЬНО XLS
 go run cmd/moshack_2022/*.go excel
ЧТОБЫ ЗАПУСТЬ СЕРВЕР:
 go run cmd/moshack_2022/*.go
