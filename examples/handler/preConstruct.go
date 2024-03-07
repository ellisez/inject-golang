package handler

import (
	"fmt"
	"github.com/ellisez/inject-golang/examples/model"
)

func PrepareWebCtxAlias() *model.WebApp {
	fmt.Println("call WebApp.preConstruct")
	return &model.WebApp{}
}
