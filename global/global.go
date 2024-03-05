package global

var (
	GenPackage = "ctx"
	StructName = "Ctx"

	CurrentDirectory string

	FlagAll       = true
	FlagSingleton = false
	FlagMultiple  = false
	FlagFunc      = false
	FlagWeb       = false

	GenCtxFilename         = "gen_ctx.go"
	GenConstructorFilename = "gen_constructor.go"
	GenFuncFilename        = "gen_func.go"
	GenMethodFilename      = "gen_method.go"
	GenWebFilename         = "gen_web.go"

	ScanDirectories []string
)
