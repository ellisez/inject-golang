package global

import (
	"github.com/ellisez/inject-golang/model"
)

var (
	FlagAll       = true
	FlagSingleton = false
	FlagMultiple  = false
	FlagFunc      = false
	FlagWeb       = false
)

var (
	GenPackage     = "ctx"
	CtxType        = "Ctx"
	ArgumentVar    = "__args__"
	InternalVar    = "__internal__"
	GenCtxFilename = "gen_ctx.go"

	GenFactoryPackage  = "factory"
	GenFactoryFilename = "gen_factory.go"

	GenInternalPackage   = "internal"
	GenSingletonFilename = "gen_singleton.go"
	GenMultipleFilename  = "gen_multiple.go"
	GenFuncFilename      = "gen_func.go"
	GenMethodFilename    = "gen_method.go"
	GenWebFilename       = "gen_web.go"

	GenUtilsPackage     = "utils"
	GenWebUtilsFilename = "gen_web_utils.go"
)

var (
	ScanDirectories []string

	CacheModMap = make(map[string]*model.Module)

	Mod *model.Module

	CtxFieldMap    = map[string][]*model.Field{}
	ImportAliasMap = map[string]*model.Import{}
	ImportPathMap  = map[string]*model.Import{}
)
