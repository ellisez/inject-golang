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

type sortStructInfo []*model.StructInfo

func (s sortStructInfo) Len() int {
	return len(s)
}

func (s sortStructInfo) Less(x int, y int) bool {
	if s[x].Order != "" && s[y].Order == "" {
		return true
	}
	return s[x].Order < s[y].Order
}

func (s sortStructInfo) Swap(x int, y int) {
	old := s[x]
	s[x] = s[y]
	s[y] = old
}

func SortStructInfo(structInfos []*model.StructInfo) {
	var sorter sortStructInfo = structInfos
	sort.Sort(sorter)
}
