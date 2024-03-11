package model

import "strings"

type Value struct {
	Name    string
	Type    string
	Default string
	Comment string
}

func (v *Value) PrivateName() string {
	return strings.ToLower(v.Name[0:1]) + v.Name[1:]
}
func (v *Value) Getter() string {
	return v.Name
}

func (v *Value) Setter() string {
	return "Set" + v.Name
}

type CtxInfo struct {
	*FuncInfo
	Values []*Value
}

func NewCtxInfoFromFunc(funcInfo *FuncInfo) *CtxInfo {
	funcInfo.Proxy = ""
	funcInfo.ProxyComment = ""
	return &CtxInfo{
		FuncInfo: funcInfo,
	}
}
