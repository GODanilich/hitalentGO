package handlers

import (
	"net/http"

	"github.com/google/uuid"

	"GODanilich/hitalentGO/internal/apperr"
	"GODanilich/hitalentGO/internal/dto"
	"GODanilich/hitalentGO/internal/http/router"
	"GODanilich/hitalentGO/internal/service"
)

type DepartmentsHandler struct {
	svc *service.DepartmentService
}

func NewDepartmentsHandler(svc *service.DepartmentService) *DepartmentsHandler {
	return &DepartmentsHandler{svc: svc}
}

// POST /departments
func (h *DepartmentsHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateDepartmentRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeAppError(w, err)
		return
	}

	resp, err := h.svc.CreateDepartment(r.Context(), req)
	if err != nil {
		writeAppError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// POST /departments/{id}/employees
func (h *DepartmentsHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	idStr := router.Param(r, "id")

	depID, err := uuid.Parse(idStr)
	if err != nil {
		writeAppError(w, apperr.Validation("invalid department id (uuid expected)"))
		return
	}

	var req dto.CreateEmployeeRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeAppError(w, err)
		return
	}

	resp, err := h.svc.CreateEmployee(r.Context(), depID, req)
	if err != nil {
		writeAppError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}
