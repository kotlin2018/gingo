package main

import (
	"gingo/command/api"
	"github.com/spf13/pflag"
	"log"
	"os"
)
// 接收命令行参数的变量
// var cliName = pflag.StringP("name","n","ginApp","input -n=$value(your program name)")
// var cliApi = pflag.BoolP("api","a",true,"input --api=$val(true or false) and you will get api template")

// 终端输入 gingo -n test 或者 gingo new test 即可创建test项目
var cliName = pflag.StringP("name","n","ginApp","input new=$value(your program name)")

func main(){

	pflag.Parse()
	//获取当前文件路径
	currentPath, _ := os.Getwd()
	if !api.IsExist(currentPath) {
		log.Printf("Application '%s' already exists", currentPath)
		os.Exit(0)
	}
	//if *cliApi {
	//	api.CreateApp(currentPath,*cliName)
	//	return
	//}
	api.CreateApp(currentPath,*cliName)
}
