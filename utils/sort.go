package utils

import (
	"github.com/ellisez/inject-golang/model"
	"go/ast"
	"sort"
)

type importSorter []*ast.ImportSpec

func (s importSorter) Len() int {
	return len(s)
}

func (s importSorter) Less(x int, y int) bool {
	return s[x].Path.Value < s[y].Path.Value
}

func (s importSorter) Swap(x int, y int) {
	old := s[x]
	s[x] = s[y]
	s[y] = old
}

func SortImports(importSpecs []*ast.ImportSpec) {
	var sorter importSorter = importSpecs
	sort.Sort(sorter)
}

type sortableInstance []model.Instance

func (s sortableInstance) Len() int {
	return len(s)
}

func (s sortableInstance) Less(x int, y int) bool {
	if s[x].GetOrder() != "" && s[y].GetOrder() == "" {
		return true
	}
	return s[x].GetOrder() < s[y].GetOrder()
}

func (s sortableInstance) Swap(x int, y int) {
	old := s[x]
	s[x] = s[y]
	s[y] = old
}

func SortInstance(instances []model.Instance) {
	var sorter sortableInstance = instances
	sort.Sort(sorter)
}
