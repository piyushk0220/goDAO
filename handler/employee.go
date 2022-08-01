package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"src/DaoInterface/driver"

	"github.com/go-chi/chi"

	models "src/DaoInterface/model"
	repository "src/DaoInterface/repository"
	employee "src/DaoInterface/repository/emp"
)

// NewPostHandler ...
func NewEmpHandler(db *driver.DB) *Employee {
	return &Employee{
		repo: employee.NewSQLEmpRepo(db.SQL),
	}
}

// Post ...
type Employee struct {
	repo repository.EmpRepo
}

// Fetch all emp data
func (p *Employee) Fetch(w http.ResponseWriter, r *http.Request) {
	payload, _ := p.repo.Fetch(r.Context(), 5)

	respondwithJSON(w, http.StatusOK, payload)
}

// Create a new emp
func (p *Employee) Create(w http.ResponseWriter, r *http.Request) {
	emp := models.Employee{}
	json.NewDecoder(r.Body).Decode(&emp)

	newID, err := p.repo.Create(r.Context(), &emp)
	fmt.Println(newID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
	}

	respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Successfully Created"})
}

// Update a emp by id
func (p *Employee) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	data := models.Employee{Id: int64(id)}
	json.NewDecoder(r.Body).Decode(&data)
	payload, err := p.repo.Update(r.Context(), &data)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
	}

	respondwithJSON(w, http.StatusOK, payload)
}

// GetByID returns a emp details
func (p *Employee) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	payload, err := p.repo.GetByID(r.Context(), int64(id))

	if err != nil {
		respondWithError(w, http.StatusNoContent, "Content not found")
	}

	respondwithJSON(w, http.StatusOK, payload)
}

// Delete a emp
func (p *Employee) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, err := p.repo.Delete(r.Context(), int64(id))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
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
