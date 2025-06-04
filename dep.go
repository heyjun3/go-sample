package gosample

import (
	"fmt"
)

type Department struct {
	ID                         string
	Name                       string
	ParentID                   string
	NumberOfEmployee           int
	NumberOfEmployeeInChildren int
	Children                   []*Department
}

func NewDepartment(id, name, parentID string, numberOfEmployee, numberOfEmployeeInChildren int) *Department {
	return &Department{
		ID:                         id,
		Name:                       name,
		ParentID:                   parentID,
		NumberOfEmployee:           numberOfEmployee,
		NumberOfEmployeeInChildren: numberOfEmployeeInChildren,
		Children:                   []*Department{},
	}
}

func DepartmentTree(departments []*Department) []*Department {
	deptMap := make(map[string]*Department)
	for _, dept := range departments {
		dept.Children = []*Department{} // clear children in case of reuse
		deptMap[dept.ID] = dept
	}
	for _, dept := range departments {
		if dept.ParentID != "" {
			if parent, ok := deptMap[dept.ParentID]; ok {
				parent.Children = append(parent.Children, dept)
			}
		}
	}
	// ルートノードのみ返す
	var roots []*Department
	for _, dept := range departments {
		if dept.ParentID == "" {
			roots = append(roots, dept)
		}
	}
	return roots
}

type GroupDepartment struct {
	departments   []*Department
	status        string
	employeeCount int
	results       []Result
}

func GroupingDepartmentsByNumberOfEmployee(department *Department, group *GroupDepartment) []*GroupDepartment {
	result := []*GroupDepartment{}
	if department.NumberOfEmployee > 0 {
		result = append(result, &GroupDepartment{
			departments:   []*Department{department},
			status:        "hasEmployee",
			employeeCount: department.NumberOfEmployee,
		})
		if group != nil {
			group.employeeCount -= department.NumberOfEmployee
		}
	} else if department.NumberOfEmployeeInChildren > 0 {
		if group != nil {
			group.employeeCount -= department.NumberOfEmployeeInChildren
		}
		group = &GroupDepartment{
			departments:   []*Department{department},
			status:        "hasEmployeeInChildren",
			employeeCount: department.NumberOfEmployeeInChildren,
		}
		result = append(result, group)
	} else if group != nil {
		group.departments = append(group.departments, department)
	} else {
		result = append(result, &GroupDepartment{departments: []*Department{department}, status: "noEmployee"})
	}
	for _, child := range department.Children {
		childResult := GroupingDepartmentsByNumberOfEmployee(child, group)
		result = append(result, childResult...)
	}
	return result
}

type Result struct {
	DepartmentID string
	EmployeeID   string
}

func Do() {
	// Example usage of DepartmentTree
	results := []Result{
		{DepartmentID: "1", EmployeeID: "E1"},
		{DepartmentID: "1", EmployeeID: "E2"},
	}
	resultMap := make(map[string][]Result)
	for _, result := range results {
		exist, ok := resultMap[result.DepartmentID]
		if ok {
			exist = append(exist, result)
			resultMap[result.DepartmentID] = exist
		} else {
			resultMap[result.DepartmentID] = []Result{result}
		}
	}
	departments := []*Department{
		NewDepartment("1", "HR", "", 0, 10),
		NewDepartment("2", "Engineering", "1", 1, 0),
		NewDepartment("3", "Sales", "2", 0, 5),
		NewDepartment("4", "Recruitment", "3", 0, 0),
		NewDepartment("5", "Development", "4", 0, 0),
	}

	tree := DepartmentTree(departments)
	r := GroupingDepartmentsByNumberOfEmployee(tree[0], nil)
	for _, group := range r {
		for _, dep := range group.departments {
			if result, ok := resultMap[dep.ID]; ok {
				group.results = append(group.results, result...)
			}
		}
		if len(group.results) >= group.employeeCount {
			fmt.Println("Group:", group.status, "Employee Count:", group.employeeCount)
		}
	}
}
