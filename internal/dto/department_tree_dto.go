package dto

type DepartmentNodeResponse struct {
	Department DepartmentResponse       `json:"department"`
	Employees  []EmployeeResponse       `json:"employees,omitempty"`
	Children   []DepartmentNodeResponse `json:"children"`
}
