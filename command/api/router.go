package api

var routerUser = `package router

import (
	"{{.Appname}}/api/v1"
	"{{.Appname}}/middleware"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user").
		Use(middleware.JWTAuth()).
		Use(middleware.CasbinHandler()).
		Use(middleware.OperationRecord())
	{
		UserRouter.POST("changePassword", v1.ChangePassword)     // 修改密码
		UserRouter.POST("uploadHeaderImg", v1.UploadHeaderImg)   // 上传头像
		UserRouter.POST("getUserList", v1.GetUserList)           // 分页获取用户列表
		UserRouter.POST("setUserAuthority", v1.SetUserAuthority) // 设置用户权限
		UserRouter.DELETE("deleteUser", v1.DeleteUser)           // 删除用户
	}
}
`

var routerMenu = `package router

import (
	"{{.Appname}}/api/v1"
	"{{.Appname}}/middleware"
	"github.com/gin-gonic/gin"
)

func InitMenuRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	MenuRouter := Router.Group("menu").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		MenuRouter.POST("getMenu", v1.GetMenu)                   // 获取菜单树
		MenuRouter.POST("getMenuList", v1.GetMenuList)           // 分页获取基础menu列表
		MenuRouter.POST("addBaseMenu", v1.AddBaseMenu)           // 新增菜单
		MenuRouter.POST("getBaseMenuTree", v1.GetBaseMenuTree)   // 获取用户动态路由
		MenuRouter.POST("addMenuAuthority", v1.AddMenuAuthority) //	增加menu和角色关联关系
		MenuRouter.POST("getMenuAuthority", v1.GetMenuAuthority) // 获取指定角色menu
		MenuRouter.POST("deleteBaseMenu", v1.DeleteBaseMenu)     // 删除菜单
		MenuRouter.POST("updateBaseMenu", v1.UpdateBaseMenu)     // 更新菜单
		MenuRouter.POST("getBaseMenuById", v1.GetBaseMenuById)   // 根据id获取菜单
	}
	return MenuRouter
}
`