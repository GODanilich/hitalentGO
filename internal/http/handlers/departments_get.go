package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"GODanilich/hitalentGO/internal/apperr"
	"GODanilich/hitalentGO/internal/http/router"
)

func (h *DepartmentsHandler) GetDepartment(w http.ResponseWriter, r *http.Request) {

	idStr := router.Param(r, "id")
	depID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || depID <= 0 {
		writeAppError(w, apperr.Validation("invalid department id (int64 expected)"))
		return
	}

	q := r.URL.Query()

	depth := 1
	if v := strings.TrimSpace(q.Get("depth")); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			writeAppError(w, apperr.Validation("depth must be a number"))
			return
		}
		depth = n
	}

	includeEmployees := true
	if v := strings.TrimSpace(q.Get("include_employees")); v != "" {
		b, err := strconv.ParseBool(v) // accepts: 1,t,T,TRUE,true,True / 0,f,F,FALSE,false,False
		if err != nil {
			writeAppError(w, apperr.Validation("include_employees must be boolean"))
			return
		}
		includeEmployees = b
	}

	resp, err := h.svc.GetDepartmentTree(r.Context(), depID, depth, includeEmployees)
	if err != nil {
		writeAppError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
