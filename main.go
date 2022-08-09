package main

import (
	"fmt"
	"net/http"
	"os"
	"src/DaoInterface/driver"
	ph "src/DaoInterface/handler"
	"src/DaoInterface/repository"
	"src/DaoInterface/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	dbName := "godb"
	dbPass := "root"
	dbHost := "localhost"
	dbPort := "3306"

	connection, err := driver.ConnectSQL(dbHost, dbPort, "root", dbPass, dbName)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	mySql := repository.NewSQLEmpRepo(connection.SQL)
	serv := service.NewEmpService(mySql)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	pHandler := ph.NewEmpHandler(serv)
	r.Route("/", func(rt chi.Router) {
		rt.Mount("/employees", employeeRouter(pHandler))
	})

	fmt.Println("Server listen at :8080")
	http.ListenAndServe(":8080", r)
}

// A completely separate router for employees routes
func employeeRouter(pHandler *ph.Employee) http.Handler {
	r := chi.NewRouter()
	r.Get("/", pHandler.Fetch)
	r.Get("/{id:[0-9]+}", pHandler.GetByID)
	r.Post("/", pHandler.Create)
	r.Put("/{id:[0-9]+}", pHandler.Update)
	r.Delete("/{id:[0-9]+}", pHandler.Delete)

	return r
}
