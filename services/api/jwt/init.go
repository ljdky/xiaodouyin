package jwt

import (
	"context"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"xiaodouyin/kitex_gen/user"
	"xiaodouyin/services/api/rpc"
)

// jwt
var AuthMiddleware *jwt.GinJWTMiddleware

// 登录请求体
type login struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

// 登录请求返回体
type loginResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     uint   `json:"user_id"`
	Token      string `json:"token"`
}

// 初始化 jwt 中间件
func init() {
	// 声明错误
	var err error

	// jwt 验证字段
	identityKey := "id"

	// 初始化 jwt middleware
	AuthMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("bainan"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		// 登陆验证，LoginHandler 会首先调用该方法
		Authenticator: func(c *gin.Context) (interface{}, error) {

			var loginVals login
			if err := c.BindQuery(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			// 验证
			if len(username) == 0 ||
				len(username) > 32 ||
				len(password) == 0 ||
				len(password) < 5 ||
				len(password) > 32 {

				return nil, jwt.ErrFailedAuthentication
			}

			// rpc 验证
			rpcReq := user.AuthRequest{
				Username: username,
				Password: password,
			}

			rpcResp, err := rpc.UserService.AuthUser(context.Background(), &rpcReq)

			if err != nil || rpcResp.StatusCode != 0 {
				return nil, jwt.ErrFailedAuthentication
			}

			// 验证成功需要把 userId 加入上下文，以传递给下游返回
			c.Set("userId", uint(rpcResp.UserId))
			return uint(rpcResp.UserId), nil // userId

		},
		// LoginHandler 接下来会调用该方法，配置 jwt 负载
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(uint); ok {
				return jwt.MapClaims{identityKey: v}
			}
			return jwt.MapClaims{}
		},
		// LoginHandler 最终会调用该方法，配置返回体
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			if code == http.StatusOK {
				c.JSON(code, loginResponse{
					StatusCode: 0,
					StatusMsg:  "success",
					UserId:     c.GetUint("userId"),
					Token:      message,
				})
			} else {
				c.JSON(code, loginResponse{StatusCode: -1, StatusMsg: "error, loginFailed"})
			}
		},
		// 验证 handler
		IdentityHandler: func(ctx *gin.Context) interface{} {
			claims := jwt.ExtractClaims(ctx)
			return uint(claims[identityKey].(float64))
		},
		// 验证器，接收上一个方法的返回值
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(uint); ok {
				// 登陆成功, 设置 userId 上下文
				c.Set("userId", v)
				return true
			}

			return false
		},
		// 验证失败处理器
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, loginResponse{StatusCode: -1, StatusMsg: message})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "query:token",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		panic(err.Error())
	}
}
