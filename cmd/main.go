package main

import (
	"github.com/NVTer/rest-api-example/internal/handler"
	"net/http"

	"github.com/NVTer/rest-api-example/internal/middleware"
	"github.com/NVTer/rest-api-example/internal/repository"
	"github.com/NVTer/rest-api-example/internal/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Handler interface { // nolint: deadcode
	GetPositions(w http.ResponseWriter, r *http.Request)
	GetEmployees(w http.ResponseWriter, r *http.Request)
	GetPosition(w http.ResponseWriter, r *http.Request)
	GetEmployee(w http.ResponseWriter, r *http.Request)
	CreatePosition(w http.ResponseWriter, r *http.Request)
	CreateEmployee(w http.ResponseWriter, r *http.Request)
	DeletePosition(w http.ResponseWriter, r *http.Request)
	DeleteEmployee(w http.ResponseWriter, r *http.Request)
	UpdatePosition(w http.ResponseWriter, r *http.Request)
	UpdateEmployee(w http.ResponseWriter, r *http.Request)
}

const (
	pathPositions  = "/positions"
	pathEmployees  = "/employees"
	pathPosition   = "/position"
	pathEmployee   = "/employee"
	pathPositionID = "/position/{id:\\S+}"
	pathEmployeeID = "/employee/{id:\\S+}"
)

func Run() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	r := mux.NewRouter()
	myData := repository.NewDataBase()
	myRepo := repository.NewRepo(myData)
	myServ := service.NewServ(myRepo)
	myH := handler.NewHandler(myServ)
	pathLimit := "{limit:\\S+}"
	pathOffset := "{offset:\\S+}"
	r.HandleFunc(pathPositions, myH.GetPositions).Queries("limit", pathLimit, "offset", pathOffset).Methods("GET")
	r.HandleFunc(pathEmployees, myH.GetEmployees).Queries("limit", pathLimit, "offset", pathOffset).Methods("GET")
	r.HandleFunc(pathPositionID, myH.GetPosition).Methods("GET")
	r.HandleFunc(pathEmployeeID, myH.GetEmployee).Methods("GET")
	r.HandleFunc(pathPositionID, myH.DeletePosition).Methods("DELETE")
	r.HandleFunc(pathEmployeeID, myH.DeleteEmployee).Methods("DELETE")
	r.HandleFunc(pathPosition, myH.UpdatePosition).Methods("PUT")
	r.HandleFunc(pathEmployee, myH.UpdateEmployee).Methods("PUT")
	r.HandleFunc(pathPosition, myH.CreatePosition).Methods("POST")
	r.HandleFunc(pathEmployee, myH.CreateEmployee).Methods("POST")
	log := logrus.New()
	r.Use(middleware.IDMiddleware(log), middleware.TimeLogMiddleware(log), middleware.AccessLogMiddleware(log))
	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	Run()
}
