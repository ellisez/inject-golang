package main

import (
	"generate/gen"
	. "generate/global"
	"generate/scan"
	"os"
)

func init() {
	var err error
	RootDirectory, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}
func main() {
	// 步骤： 创建实例到容器 > 互相注入属性
	// 1. 扫描目录下所有代码文件，
	// 2. 解析ast
	// 2.1 结构体上的注解：
	// @provide <实例名，默认同类名> <singleton默认|multiple>
	// @constructor <构造函数名，默认New+类名>
	// 2.2 结构体的属性注解：
	// @inject <实例名，默认同类名>
	// 2.3 方法上的注解:
	// @proxy <代理方法名，默认同方法名>
	// 2.4 路由方法上的注解（参照swag注解）：
	// @route <path> <httpMethod: get|post>
	// @param <参数名> <参数类型:query|path|header|body|formData|inject> <接收类型> <必填与否> <参数说明>
	// 2.5 方法参数上的注解:
	// @inject <实例名，默认同类名>
	annotateInfo, err := scan.DoScan()
	if err != nil {
		panic(err)
	}
	// 3. 单例生成到容器结构体中，
	// 4. 多例则生成容器New方法片段中，
	// 5. 方法注入则生成容器注入方法片段中，
	err = gen.DoGen(annotateInfo)
	if err != nil {
		panic(err)
	}

}
