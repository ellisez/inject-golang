package utils

import (
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
