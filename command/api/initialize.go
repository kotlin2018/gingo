package api

var dbTable = `package initialize

import (
	"{{.Appname}}/global"
	"{{.Appname}}/model"
)

// 注册数据库表专用
func DBTables() {
	db := global.GVA_DB
	db.AutoMigrate(model.SysUser{},
		model.SysAuthority{},
		//model.SysApi{},
		//model.SysBaseMenu{},
		//model.SysBaseMenuParameter{},
		//model.JwtBlacklist{},
		//model.SysWorkflow{},
		//model.SysWorkflowStepInfo{},
		//model.SysDictionary{},
		//model.SysDictionaryDetail{},
		//model.ExaFileUploadAndDownload{},
		//model.ExaFile{},
		//model.ExaFileChunk{},
		//model.ExaSimpleUploader{},
		//model.ExaCustomer{},
		//model.SysOperationRecord{},
	)
	global.GVA_LOG.Debug("register table success")
}
`

var mysql = `package initialize

import (
	"{{.Appname}}/global"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

// 初始化数据库并产生数据库全局变量
func Mysql() {
	admin := global.GVA_CONFIG.Mysql
	if db, err := gorm.Open("mysql", admin.Username+":"+admin.Password+"@("+admin.Path+")/"+admin.Dbname+"?"+admin.Config); err != nil {
		global.GVA_LOG.Error("MySQL启动异常", err)
		os.Exit(0)
	} else {
		global.GVA_DB = db
		global.GVA_DB.DB().SetMaxIdleConns(admin.MaxIdleConns)
		global.GVA_DB.DB().SetMaxOpenConns(admin.MaxOpenConns)
		global.GVA_DB.LogMode(admin.LogMode)
	}
}
`

var redis = `package initialize

import (
	"{{.Appname}}/global"
	"github.com/go-redis/redis"
)

func Redis() {
	redisCfg := global.GVA_CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	pong, err := client.Ping().Result()
	if err != nil {
		global.GVA_LOG.Error(err)
	} else {
		global.GVA_LOG.Info("redis connect ping response:", pong)
		global.GVA_REDIS = client
	}
}
`

var router = `package initialize

import (
	_ "{{.Appname}}/docs"
	"{{.Appname}}/global"
	"{{.Appname}}/middleware"
	"{{.Appname}}/router"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

// 初始化总路由

func Routers() *gin.Engine {
	var Router = gin.Default()
	// 为用户头像和文件提供静态地址
	Router.StaticFS(global.GVA_CONFIG.LocalUpload.AvatarPath, http.Dir(global.GVA_CONFIG.LocalUpload.AvatarPath))
	Router.StaticFS(global.GVA_CONFIG.LocalUpload.FilePath, http.Dir(global.GVA_CONFIG.LocalUpload.FilePath))
	// Router.Use(middleware.LoadTls())  // 打开就能玩https了
	global.GVA_LOG.Debug("use middleware logger")
	// 跨域
	Router.Use(middleware.Cors())
	global.GVA_LOG.Debug("use middleware cors")
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.GVA_LOG.Debug("register swagger handler")
	// 方便统一添加路由组前缀 多服务器上线使用
	ApiGroup := Router.Group("")
	router.InitUserRouter(ApiGroup)                  // 注册用户路由
	router.InitBaseRouter(ApiGroup)                  // 注册基础功能路由 不做鉴权
	router.InitMenuRouter(ApiGroup)                  // 注册menu路由
	router.InitAuthorityRouter(ApiGroup)             // 注册角色路由
	router.InitApiRouter(ApiGroup)                   // 注册功能api路由
	router.InitFileUploadAndDownloadRouter(ApiGroup) // 文件上传下载功能路由
	router.InitSimpleUploaderRouter(ApiGroup)        // 断点续传（插件版）
	router.InitWorkflowRouter(ApiGroup)              // 工作流相关路由
	router.InitCasbinRouter(ApiGroup)                // 权限相关路由
	router.InitJwtRouter(ApiGroup)                   // jwt相关路由
	router.InitSystemRouter(ApiGroup)                // system相关路由
	router.InitCustomerRouter(ApiGroup)              // 客户路由
	router.InitAutoCodeRouter(ApiGroup)              // 创建自动化代码
	router.InitSysDictionaryDetailRouter(ApiGroup)   // 字典详情管理
	router.InitSysDictionaryRouter(ApiGroup)         // 字典管理
	router.InitSysOperationRecordRouter(ApiGroup)    // 操作记录

	global.GVA_LOG.Info("router register success")
	return Router
}

`

var sqLite =`package initialize

// sqlite需要gcc支持 windows用户需要自行安装gcc 如需使用打开注释即可

// 感谢 sqlitet提供者 [rikugun] 作者github： https://github.com/rikugun

// import (
// 	"fmt"
// 	"{{.Appname}}/global"
// 	"github.com/jinzhu/gorm"
// 	_ "github.com/jinzhu/gorm/dialects/sqlite"
// )
//
// // 初始化数据库并产生数据库全局变量
// func Sqlite() {
// 	admin := global.GVA_CONFIG.Sqlite
// 	if db, err := gorm.Open("sqlite3", fmt.Sprintf("%s?%s", admin.Path,admin.Config)); err != nil {
// 		global.GVA_LOG.Error("DEFAULTDB数据库启动异常", err)
// 	} else {
// 		global.GVA_DB = db
// 		global.GVA_DB.LogMode(admin.LogMode)
// 	}
// }
`

var validator = `package initialize

import "{{.Appname}}/utils"

func init() {
	_ = utils.RegisterRule("PageVerify",
		utils.Rules{
			"Page":     {utils.NotEmpty()},
			"PageSize": {utils.NotEmpty()},
		},
	)
	_ = utils.RegisterRule("IdVerify",
		utils.Rules{
			"Id": {utils.NotEmpty()},
		},
	)
	_ = utils.RegisterRule("AuthorityIdVerify",
		utils.Rules{
			"AuthorityId": {utils.NotEmpty()},
		},
	)
}
`
