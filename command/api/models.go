package api

var reqUser = `package request
// 客户端的请求参数结构体
import uuid "github.com/satori/go.uuid"

// User register structure
type RegisterStruct struct {
	Username    string // json:"userName"
	Password    string // json:"passWord"
	NickName    string // json:"nickName" gorm:"default:'QMPlusUser'"
	HeaderImg   string // json:"headerImg" gorm:"default:'http://www.henrongyi.top/avatar/lufu.jpg'"
	AuthorityId string // json:"authorityId" gorm:"default:888"
}

// User login structure
type RegisterAndLoginStruct struct {
	Username  string // json:"username"
	Password  string // json:"password"
	Captcha   string // json:"captcha"
	CaptchaId string // json:"captchaId"
}

// Modify password structure
type ChangePasswordStruct struct {
	Username    string // json:"username"
	Password    string // json:"password"
	NewPassword string // json:"newPassword"
}

// Modify  user's auth structure
type SetUserAuth struct {
	UUID        uuid.UUID // json:"uuid"
	AuthorityId string    // json:"authorityId"
}
`

var respUser = `package response

// 服务端的响应参数结构体
import (
	"{{.Appname}}/model"
)

type SysUserResponse struct {
	User model.SysUser // json:"user"
}

type LoginResponse struct {
	User      model.SysUser // json:"user"
	Token     string        // json:"token"
	ExpiresAt int64         // json:"expiresAt"
}




`

var modUser = `package models

// 生产数据库表的实体
import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type SysUser struct {
	gorm.Model
	UUID        uuid.UUID    // json:"uuid" gorm:"comment:'用户UUID'"
	Username    string       // json:"userName" gorm:"comment:'用户登录名'"
	Password    string       // json:"-"  gorm:"comment:'用户登录密码'"
	NickName    string       // json:"nickName" gorm:"default:'系统用户';comment:'用户昵称'" 
	HeaderImg   string       // json:"headerImg" gorm:"default:'http://qmplusimg.henrongyi.top/head.png';comment:'用户头像'"
	Authority   SysAuthority // json:"authority" gorm:"ForeignKey:AuthorityId;AssociationForeignKey:AuthorityId;comment:'用户角色'"
	AuthorityId string       // json:"authorityId" gorm:"default:888;comment:'用户角色ID'"
}
`