package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"GODanilich/hitalentGO/internal/apperr"
	"GODanilich/hitalentGO/internal/http/router"
)

func (h *DepartmentsHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {

	idStr := router.Param(r, "id")
	depID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || depID <= 0 {
		writeAppError(w, apperr.Validation("invalid department id (int64 expected)"))
		return
	}

	q := r.URL.Query()
	mode := strings.TrimSpace(q.Get("mode"))

	var reassignTo *int64
	if strings.TrimSpace(q.Get("reassign_to_department_id")) != "" {
		v := strings.TrimSpace(q.Get("reassign_to_department_id"))
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			writeAppError(w, apperr.Validation("reassign_to_department_id must be int64"))
			return
		}
		reassignTo = &n
	}

	if err := h.svc.DeleteDepartment(r.Context(), depID, mode, reassignTo); err != nil {
		writeAppError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
