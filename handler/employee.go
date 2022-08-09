package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"src/DaoInterface/service"

	"github.com/go-chi/chi"

	models "src/DaoInterface/model"
)

// NewPostHandler ...
func NewEmpHandler(ser service.EmpService) *Employee {
	return &Employee{
		service: ser,
	}
}

// Post ...
type Employee struct {
	service service.EmpService
}

// Fetch all emp data
func (p *Employee) Fetch(w http.ResponseWriter, r *http.Request) {
	payload, err := p.service.Fetch()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}
	respondwithJSON(w, http.StatusOK, payload)
}

// Create a new emp
func (p *Employee) Create(w http.ResponseWriter, r *http.Request) {
	emp := models.Employee{}
	json.NewDecoder(r.Body).Decode(&emp)

	newID, err := p.service.Create(&emp)
	fmt.Println(newID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Successfully Created"})
}

// Update a emp by id
func (p *Employee) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	data := models.Employee{Id: int64(id)}
	json.NewDecoder(r.Body).Decode(&data)
	payload, err := p.service.Update(&data)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	respondwithJSON(w, http.StatusOK, payload)
}

// GetByID returns a emp details
func (p *Employee) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	payload, err := p.service.GetByID(int64(id))

	if err != nil {
		respondWithError(w, http.StatusNoContent, "Content not found")
		return
	}

	respondwithJSON(w, http.StatusOK, payload)
}

// Delete a emp
func (p *Employee) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, err := p.service.Delete(int64(id))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	respondwithJSON(w, http.StatusMovedPermanently, map[string]string{"message": "Delete Successfully"})
}

// respondwithJSON write json response format
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondwithJSON(w, code, map[string]string{"message": msg})
}
