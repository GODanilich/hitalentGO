package handlers

import (
	"net/http"
	"strconv"

	"GODanilich/hitalentGO/internal/apperr"
	"GODanilich/hitalentGO/internal/dto"
	"GODanilich/hitalentGO/internal/http/router"
)

func (h *DepartmentsHandler) PatchDepartment(w http.ResponseWriter, r *http.Request) {
	idStr := router.Param(r, "id")
	depID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || depID <= 0 {
		writeAppError(w, apperr.Validation("invalid department id (int64 expected)"))
		return
	}

	var req dto.PatchDepartmentRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeAppError(w, err)
		return
	}

	resp, err := h.svc.UpdateDepartment(r.Context(), depID, req)
	if err != nil {
		writeAppError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
