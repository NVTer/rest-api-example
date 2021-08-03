package handler

import (
	"encoding/json"
	errs "errors"
	"net/http"
	"strconv"

	"github.com/NVTer/rest-api-example/internal"
	"github.com/NVTer/rest-api-example/internal/errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

type Hand struct {
	service Service
}

func NewHandler(service Service) *Hand {
	return &Hand{service: service}
}

func (h *Hand) GetPositions(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	positions, err := h.service.GetPositions(r.Context(), limit, offset)
	if err != nil {
		if errs.Is(err, errors.NotFound()) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errs.Is(err, errors.BadRequest()) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(positions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
}

func (h *Hand) GetEmployees(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	employees, err := h.service.GetEmployees(r.Context(), limit, offset)
	if err != nil {
		if errs.Is(err, errors.NotFound()) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errs.Is(err, errors.BadRequest()) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(employees)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
}

func (h *Hand) GetPosition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	p, err := h.service.GetPosition(r.Context(), vars["id"])
	if err != nil {
		if errs.Is(err, errors.NotFound()) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
}

func (h *Hand) GetEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	e, err := h.service.GetEmployee(r.Context(), vars["id"])
	if err != nil {
		if errs.Is(err, errors.NotFound()) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
}

func (h *Hand) CreatePosition(w http.ResponseWriter, r *http.Request) {
	var p internal.Position
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if p.Salary == decimal.Zero || p.Name == "" {
		http.Error(w, errors.BadRequest().Error(), http.StatusBadRequest)
		return
	}
	id, err := h.service.CreatePosition(r.Context(), &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := map[string]string{
		"id": id,
	}
	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(201)
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
}

func (h *Hand) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var e internal.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, errors.ParseError().Error(), http.StatusInternalServerError)
		return
	}
	if e.LasName == "" || e.FirstName == "" || e.PositionID == uuid.Nil {
		http.Error(w, errors.BadRequest().Error(), http.StatusBadRequest)
		return
	}
	id, err := h.service.CreateEmployee(r.Context(), &e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := map[string]string{
		"id": id,
	}
	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(201)
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
}

func (h *Hand) UpdatePosition(w http.ResponseWriter, r *http.Request) {
	var p internal.Position
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err := h.service.UpdatePosition(r.Context(), &p)
	if err != nil {
		if errs.Is(err, errors.BadRequest()) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
}

func (h *Hand) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	var e internal.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err := h.service.UpdateEmployee(r.Context(), &e)
	if err != nil {
		if errs.Is(err, errors.BadRequest()) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errs.Is(err, errors.PositionIsNotExists()) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
}

func (h *Hand) DeletePosition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	err := h.service.DeletePosition(r.Context(), vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(internal.Position{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
}

func (h *Hand) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	err := h.service.DeleteEmployee(r.Context(), vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(internal.Employee{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, er := w.Write(jsonBytes)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
	}
}
