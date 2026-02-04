package handler_test

import (
	"app/internal/department/handler"
	"app/internal/department/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildTree_SimpleHierarchy(t *testing.T) {
	departments := []models.Department{
		{Id: 1, Name: "Root", RootId: 0},
		{Id: 2, Name: "Child1", RootId: 1},
		{Id: 3, Name: "Child2", RootId: 1},
	}

	result := handler.BuildTree(departments)

	require.Len(t, result, 1)
	assert.Equal(t, "Root", result[0].Name)
	assert.Len(t, result[0].Children, 2)

	childNames := []string{result[0].Children[0].Name, result[0].Children[1].Name}
	assert.ElementsMatch(t, []string{"Child1", "Child2"}, childNames)
}

func TestBuildTree_MultipleRoots(t *testing.T) {
	departments := []models.Department{
		{Id: 1, Name: "Root1", RootId: 0},
		{Id: 2, Name: "Root2", RootId: 0},
		{Id: 3, Name: "Child", RootId: 1},
	}

	result := handler.BuildTree(departments)

	require.Len(t, result, 2)
	assert.ElementsMatch(t,
		[]string{"Root1", "Root2"},
		[]string{result[0].Name, result[1].Name},
	)

	var root1 *handler.DepartmentTree
	for _, r := range result {
		if r.Name == "Root1" {
			root1 = r
			break
		}
		require.NotNil(t, root1)
		assert.Len(t, root1.Children, 1)
		assert.Equal(t, "Child", root1.Children[0].Name)
	}
}

func TestBuildTree_SingleRoot(t *testing.T) {
	departments := []models.Department{{Id: 1, Name: "Solo", RootId: 0}}

	result := handler.BuildTree(departments)

	require.Len(t, result, 1)
	assert.Equal(t, "Solo", result[0].Name)
	assert.Empty(t, result[0].Children)
}

func TestBuildTree_DeepHierarchy(t *testing.T) {
	departments := []models.Department{
		{Id: 1, Name: "A", RootId: 0},
		{Id: 2, Name: "B", RootId: 1},
		{Id: 3, Name: "C", RootId: 2},
		{Id: 4, Name: "D", RootId: 3},
	}

	result := handler.BuildTree(departments)

	require.Len(t, result, 1)
	assert.Equal(t, "A", result[0].Name)

	a := result[0]
	b := a.Children[0]
	c := b.Children[0]
	d := c.Children[0]

	assert.Equal(t, "B", b.Name)
	assert.Equal(t, "C", c.Name)
	assert.Equal(t, "D", d.Name)
}
