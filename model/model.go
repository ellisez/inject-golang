package model

type Module struct {
	Path    string            // the dir of go.mod
	Package string            // go.mod mod
	Version string            // go.mod version
	Require map[string]string // go.mod require
	Work    map[string]string // go.work
}

type Import struct {
	Alias   string
	Package string
	Path    string
}

type CommonFunc struct {
	Imports []*Import

	*Func

	Override bool
	Order    string
	Comment  string
}

func NewCommonFunc() *CommonFunc {
	return &CommonFunc{
		Func: &Func{},
	}
}

type Comment struct {
	Text string
	Args []string
}

type Method struct {
	From     *Func
	FuncName string
	Params   []*Field
	Results  []*Field
}
