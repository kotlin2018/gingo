package api

var cors = `package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id\"")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
`

var jwt = `package middleware

import (
	"errors"
	"{{.Appname}}/global"
	"{{.Appname}}/global/response"
	"{{.Appname}}/model"
	"{{.Appname}}/model/request"
	"{{.Appname}}/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localSstorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		if token == "" {
			response.Result(response.ERROR, gin.H{
				"reload": true,
			}, "未登录或非法访问", c)
			c.Abort()
			return
		}
		if service.IsBlacklist(token) {
			response.Result(response.ERROR, gin.H{
				"reload": true,
			}, "您的帐户异地登陆或令牌失效", c)
			c.Abort()
			return
		}
		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				response.Result(response.ERROR, gin.H{
					"reload": true,
				}, "授权已过期", c)
				c.Abort()
				return
			}
			response.Result(response.ERROR, gin.H{
				"reload": true,
			}, err.Error(), c)
			c.Abort()
			return
		}
		if claims.ExpiresAt - time.Now().Unix()<claims.BufferTime {
			claims.ExpiresAt = time.Now().Unix() + 60*60*24*7
			newToken,_ := j.CreateToken(*claims)
			newClaims,_ := j.ParseToken(newToken)
			c.Header("new-token",newToken)
			c.Header("new-expires-at",strconv.FormatInt(newClaims.ExpiresAt,10))
			if global.GVA_CONFIG.System.UseMultipoint {
				err,RedisJwtToken := service.GetRedisJWT(newClaims.Username)
				if err!=nil {
					global.GVA_LOG.Error(err)
				}else{
					service.JsonInBlacklist(model.JwtBlacklist{Jwt: RedisJwtToken})
					//当之前的取成功时才进行拉黑操作
				}
				// 无论如何都要记录当前的活跃状态
				_ = service.SetRedisJWT(newToken,newClaims.Username)
			}
		}
		c.Set("claims", claims)
		c.Next()
	}
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.GVA_CONFIG.JWT.SigningKey),
	}
}

// 创建一个token
func (j *JWT) CreateToken(claims request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*request.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}

}

// 更新token
//func (j *JWT) RefreshToken(tokenString string) (string, error) {
//	jwt.TimeFunc = func() time.Time {
//		return time.Unix(0, 0)
//	}
//	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return j.SigningKey, nil
//	})
//	if err != nil {
//		return "", err
//	}
//	if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
//		jwt.TimeFunc = time.Now
//		claims.StandardClaims.ExpiresAt = time.Now().Unix() + 60*60*24*7
//		return j.CreateToken(*claims)
//	}
//	return "", TokenInvalid
//}
`

var loadTls = `package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

// 用https把这个中间件在router里面use一下就好

func LoadTls() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:443",
		})
		err := middleware.Process(c.Writer, c.Request)
		if err != nil {
			// 如果出现错误，请不要继续
			fmt.Println(err)
			return
		}
		// 继续往下处理
		c.Next()
	}
}
`

var operation = `package middleware

import (
	"bytes"
	"{{.Appname}}/global"
	"{{.Appname}}/model"
	"{{.Appname}}/service"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		if c.Request.Method != http.MethodGet {
			var err error
			body, err = ioutil.ReadAll(c.Request.Body)
			if err != nil {
				global.GVA_LOG.Error("read body from request error:", err)
			} else {
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}
		}
		userId, err := strconv.Atoi(c.Request.Header.Get("x-user-id"))
		if err != nil {
			userId = 0
		}
		record := model.SysOperationRecord{
			Ip:     c.ClientIP(),
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
			Agent:  c.Request.UserAgent(),
			Body:   string(body),
			UserId: userId,
		}
		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		now := time.Now()

		c.Next()

		latency := time.Now().Sub(now)
		record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = c.Writer.Status()
		record.Latency = latency
		record.Resp = writer.body.String()

		if err := service.CreateSysOperationRecord(record); err != nil {
			global.GVA_LOG.Error("create operation record error:", err)
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
`

var casBin = `package middleware

import (
	"{{.Appname}}/global"
	"{{.Appname}}/global/response"
	"{{.Appname}}/model/request"
	"{{.Appname}}/service"
	"github.com/gin-gonic/gin"
)

// 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		waitUse := claims.(*request.CustomClaims)
		// 获取请求的URI
		obj := c.Request.URL.RequestURI()
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := waitUse.AuthorityId
		e := service.Casbin()
		// 判断策略中是否存在
		if global.GVA_CONFIG.System.Env == "develop" || e.Enforce(sub, obj, act) {
			c.Next()
		} else {
			response.Result(response.ERROR, gin.H{}, "权限不足", c)
			c.Abort()
			return
		}
	}
}
`
