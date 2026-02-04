package handler

import "app/internal/department/models"

type DepartmentTree struct {
	Name     string            `json:"name"`
	Children []*DepartmentTree `json:"children,omitempty"`
}

func BuildTree(departments []models.Department) []*DepartmentTree {
	nodeMap := make(map[int]*DepartmentTree)

	for _, d := range departments {
		nodeMap[d.Id] = &DepartmentTree{Name: d.Name}
	}

	var tree []*DepartmentTree
	for _, d := range departments {
		node := nodeMap[d.Id]
		if d.RootId == 0 {
			tree = append(tree, node)
		} else {
			parent, ok := nodeMap[d.RootId]
			if ok {
				if parent.Children == nil {
					parent.Children = []*DepartmentTree{}
				}
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return tree
}
