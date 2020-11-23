package api

import (
	"log"
	"os"
	"path"
	"strings"
)



func CreateApp(appPath,appName string){
	log.Println("Creating application...")

	os.MkdirAll(appName,0755) // 创建多级目录 (根据项目名,创建项目)
	os.Mkdir(path.Join(appName,"config"),0755) // 创建目录 (包)
	os.Mkdir(path.Join(appName,"api"),0755)
	os.Mkdir(path.Join(appName,"core"),0755)
	os.Mkdir(path.Join(appName,"global"),0755)
	os.Mkdir(path.Join(appName,"initialize"),0755)
	os.Mkdir(path.Join(appName,"middleware"),0755)
	os.Mkdir(path.Join(appName,"model"),0755)
	os.Mkdir(path.Join(appName,"router"),0755)
	os.Mkdir(path.Join(appName,"service"),0755)
	//os.Mkdir(path.Join(appName,"resource"),0755)
	os.Mkdir(path.Join(appName,"packfile"),0755)
	os.Mkdir(path.Join(appName,"log"),0755)
	os.Mkdir(path.Join(appName,"db"),0755)
	os.Mkdir(path.Join(appName,"utils"),0755)

	os.MkdirAll(path.Join(appName, "/api/v1"), 0755)
	os.MkdirAll(path.Join(appName, "/global/response"), 0755)
	os.MkdirAll(path.Join(appName, "/model/request"), 0755)
	os.MkdirAll(path.Join(appName, "/model/response"), 0755)
	//os.MkdirAll(path.Join(appName, "/models/repository"), 0755)
	//os.MkdirAll(path.Join(appName, "/resource/page"), 0755)
	//os.MkdirAll(path.Join(appName, "/resource/page/css"), 0755)
	//os.MkdirAll(path.Join(appName, "/resource/page/img"), 0755)
	//os.MkdirAll(path.Join(appName, "/resource/page/js"), 0755)
	//os.MkdirAll(path.Join(appName, "/resource/template"), 0755)
	//os.MkdirAll(path.Join(appName, "/resource/template/fe"), 0755)
	//os.MkdirAll(path.Join(appName, "/resource/template/te"), 0755)

	// 生成模版代码
	WriteToFile(path.Join(appName, "config", "config.yaml"), yaml)  // yaml配置文件
	WriteToFile(path.Join(appName, "config", "config.go"), config) // config.go结构体

	WriteToFile(path.Join(appName, "core", "server_other.go"), coreServerOther)
	WriteToFile(path.Join(appName, "core", "server_win.go"), coreServerWin)

	WriteToFile(path.Join(appName, "/global/response", "response.go"), globalResp)

	WriteToFile(path.Join(appName, "middleware", "cors.go"), cors)
	WriteToFile(path.Join(appName, "middleware", "loadtls.go"), loadTls)

	WriteToFile(path.Join(appName, "/model/request", "user.go"), reqUser)
	WriteToFile(path.Join(appName, "model", "user.go"), modUser)

	//WriteToFile(path.Join(appName, "/resource/template/fe", "api.js.tpl"), feApi)
	//WriteToFile(path.Join(appName, "/resource/template/fe", "table.vue.tpl"), feTable)

	//WriteToFile(path.Join(appName, "/resource/template/te", "model.go.tpl"), modelGo)

	WriteToFile(path.Join(appName, "packfile", "usePackFile.go"),UsePack)
	WriteToFile(path.Join(appName, "packfile", "notUsePackFile.go"),notUser)

	WriteToFile(path.Join(appName, "utils", "md5.go"), md5)
	WriteToFile(path.Join(appName, "utils", "validator.go"), utilsValidator)
	WriteToFile(path.Join(appName, "utils", "zipfiles.go"), zipFiles)
	WriteToFile(path.Join(appName, "utils", "struct_to_map.go"), structToMap)
	WriteToFile(path.Join(appName, "utils", "breakpoint_continue.go"), breakpoint)
	WriteToFile(path.Join(appName, "utils", "des.go"), des)
	WriteToFile(path.Join(appName, "utils", "array_to_string.go"), array)

	// 模版嵌套
	WriteToFile(path.Join(appName, "/api/v1", "user.go"), strings.Replace(api, "{{.Appname}}", appName, -1))

	WriteToFile(path.Join(appName, "core", "config.go"), strings.Replace(coreConfig, "{{.Appname}}", appName, -1))
	WriteToFile(path.Join(appName, "core", "config.go"), strings.Replace(coreLog, "{{.Appname}}", appName, -1))
	WriteToFile(path.Join(appName, "core", "config.go"), strings.Replace(coreServe, "{{.Appname}}", appName, -1))

	WriteToFile(path.Join(appName, "global", "global.go"), strings.Replace(global, "{{.Appname}}", appName, -1))

	WriteToFile(path.Join(appName, "initialize", "db_table.go"), strings.Replace(dbTable, "{{.Appname}}", appName, -1))
	WriteToFile(path.Join(appName, "initialize", "mysql.go"), strings.Replace(mysql, "{{.Appname}}", appName, -1))
	WriteToFile(path.Join(appName, "initialize", "redis.go"), strings.Replace(redis, "{{.Appname}}", appName, -1))
	WriteToFile(path.Join(appName, "initialize", "router.go"), strings.Replace(router, "{{.Appname}}", appName, -1))
	WriteToFile(path.Join(appName, "initialize", "sqlite.go"), strings.Replace(sqLite, "{{.Appname}}", appName, -1))
	WriteToFile(path.Join(appName, "initialize", "validator.go"), strings.Replace(validator, "{{.Appname}}", appName, -1))

	WriteToFile(path.Join(appName, "middleware", "jwt.go"), strings.Replace(jwt, "{{.Appname}}", appName, -1))
	WriteToFile(path.Join(appName, "middleware", "operation.go"), strings.Replace(operation, "{{.Appname}}", appName, -1))
	WriteToFile(path.Join(appName, "middleware", "casbin_rcba.go"), strings.Replace(casBin, "{{.Appname}}", appName, -1))

	WriteToFile(path.Join(appName, "/model/response", "user.go"), strings.Replace(respUser, "{{.Appname}}", appName, -1))

	WriteToFile(path.Join(appName, "router", "user.go"), strings.Replace(routerUser, "{{.Appname}}", appName, -1))
	WriteToFile(path.Join(appName, "router", "menu.go"), strings.Replace(routerMenu, "{{.Appname}}", appName, -1))

	WriteToFile(path.Join(appName, "service", "user.go"), strings.Replace(ServiceUser, "{{.Appname}}", appName, -1))

	//WriteToFile(path.Join(appName, "/resource/template/te", "api.go.tpl"), strings.Replace(apiGo, "{{.Appname}}", appName, -1))
	//WriteToFile(path.Join(appName, "/resource/template/te", "request.go.tpl"), strings.Replace(requestGo, "{{.Appname}}", appName, -1))
	//WriteToFile(path.Join(appName, "/resource/template/te", "router.go.tpl"), strings.Replace(routerGo, "{{.Appname}}", appName, -1))
	//WriteToFile(path.Join(appName, "/resource/template/te", "service.go.tpl"), strings.Replace(serviceGo, "{{.Appname}}", appName, -1))

	WriteToFile(path.Join(appName, "utils", "upload_avatar_local.go"), strings.Replace(uploadAvatar, "{{.Appname}}", appName, -1))
	WriteToFile(path.Join(appName, "utils", "upload_file_local.go"), strings.Replace(uploadFile, "{{.Appname}}", appName, -1))
	WriteToFile(path.Join(appName, "utils", "upload_remote.go"), strings.Replace(uploadRemote, "{{.Appname}}", appName, -1))
	WriteToFile(path.Join(appName, "utils", "directory.go"), strings.Replace(directory, "{{.Appname}}", appName, -1))

	WriteToFile(path.Join(appName, "main.go"), strings.Replace(main, "{{.Appname}}", appName, -1))
	log.Println("new application successfully created!")
}
// 将content写入fileName 文件中
// 将文本内容写入文件中
func WriteToFile(fileName,content string){
	f, err := os.Create(fileName)
	if err !=nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err !=nil {
		panic(err)
	}
}

// 一个文件或者目录是否存在
func IsExist(path string)bool{
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
