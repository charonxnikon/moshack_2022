package main

import (
	"html/template"
	"moshack_2022/pkg/apartments"
	"moshack_2022/pkg/handlers"
	"moshack_2022/pkg/items"
	"moshack_2022/pkg/middleware"
	"moshack_2022/pkg/session"
	"moshack_2022/pkg/user"
	"net/http"
	// "os"

	"github.com/gorilla/mux"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go.uber.org/zap"
)

func main() {

	templates := template.Must(template.ParseGlob("./templates/*"))

	dsn := "host=localhost user=postgres password=3546"
	dsn += " dbname=gusev port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err) // TODO
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err) // TODO
	}
	err = sqlDB.Ping()
	if err != nil {
		panic(err) // TODO
	}

	// if len(os.Args) > 1 {
	// 	excel := ExcelParser{fileName: "test.xls"}
	// 	jsonData := excel.parse(db).marshalExcel()
	// 	println(string(jsonData))
	// 	return
	// }

	sm := session.NewSessionsManager()
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	userRepo := user.NewMemoryRepo(db)
	itemsRepo := items.NewMemoryRepo()
	apartmentRepo := apartments.NewApartmentRepo(db)

	userHandler := &handlers.UserHandler{
		Tmpl:     templates,
		UserRepo: userRepo,
		Logger:   logger,
		Sessions: sm,
	}

	apartmentHandler := &handlers.ApartmentHandler{
		Tmpl:          templates,
		Logger:        logger,
		ApartmentRepo: apartmentRepo,
	}

	handlers := &handlers.ItemsHandler{
		Tmpl:      templates,
		Logger:    logger,
		ItemsRepo: itemsRepo,
	}

	r := mux.NewRouter()

	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(fileServer)

	r.HandleFunc("/", userHandler.Index).Methods("GET")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/logout", userHandler.Logout).Methods("POST")
	r.HandleFunc("/registration", userHandler.Registration).Methods("GET")
	r.HandleFunc("/registration", userHandler.Register).Methods("POST")

	r.HandleFunc("/loadxls", apartmentHandler.Load).Methods("GET")

	r.HandleFunc("/items", handlers.List).Methods("GET")
	r.HandleFunc("/items/new", handlers.AddForm).Methods("GET")
	r.HandleFunc("/items/new", handlers.Add).Methods("POST")
	r.HandleFunc("/items/{id}", handlers.Edit).Methods("GET")
	r.HandleFunc("/items/{id}", handlers.Update).Methods("POST")
	r.HandleFunc("/items/{id}", handlers.Delete).Methods("DELETE")

	mux := middleware.Auth(sm, r)
	mux = middleware.AccessLog(logger, mux)
	mux = middleware.Panic(mux)

	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	http.ListenAndServe(addr, mux)
}
