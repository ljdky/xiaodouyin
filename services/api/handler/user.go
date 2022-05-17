package handler

import (
	"context"
	"net/http"
	"strconv"

	"xiaodouyin/kitex_gen/user"

	"github.com/gin-gonic/gin"

	"xiaodouyin/services/api/jwt"
	"xiaodouyin/services/api/rpc"
)

// 创建用户请求的返回体
type CreateUserResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     uint   `json:"user_id"`
	Token      string `json:"token"`
}

// 创建用户请求处理
func CreateUserHandler(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")

	// 验证
	if len(username) == 0 ||
		len(username) > 32 ||
		len(password) == 0 ||
		len(password) < 5 ||
		len(password) > 32 {
		resp := &CreateUserResponse{
			StatusCode: -1,
			StatusMsg:  "error",
		}
		ctx.JSON(http.StatusBadRequest, resp)

		return
	}

	// rpc 创建用户
	req := &user.CreateUserRequest{
		Username: username,
		Password: password,
	}

	rpcResp, err := rpc.UserService.CreateUser(context.Background(), req)

	// rpc 错误
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			&CreateUserResponse{
				StatusCode: -1,
				StatusMsg:  "error",
			})

		return
	}

	// rpc 返回错误码
	if rpcResp.StatusCode != 0 {
		ctx.JSON(
			http.StatusAccepted,
			&CreateUserResponse{
				StatusCode: int(rpcResp.StatusCode),
				StatusMsg:  rpcResp.StatusMessage,
			})

		return
	}

	// Create the token
	mw := jwt.AuthMiddleware

	tokenString, _, err := mw.TokenGenerator(uint(rpcResp.UserId))

	// token 签发失败
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			&CreateUserResponse{
				StatusCode: -1,
				StatusMsg:  "error",
			})

		return
	}

	// token 签发成功
	ctx.JSON(
		http.StatusOK,
		&CreateUserResponse{
			StatusCode: int(rpcResp.StatusCode),
			StatusMsg:  rpcResp.StatusMessage,
			UserId:     uint(rpcResp.UserId),
			Token:      tokenString,
		},
	)

}

// 处理用户信息请求

type UserInfoResponse struct {
	StatusCode int      `json:"status_code"`
	StatusMsg  string   `json:"status_msg"`
	User       UserInfo `json:"user"`
}

type UserInfo struct {
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

func UserInfoHandler(ctx *gin.Context) {
	// 取出需要查询的 id
	userIdQuery, ok := strconv.ParseUint(ctx.Query("user_id"), 10, 64)
	if ok != nil {
		ctx.JSON(
			http.StatusBadRequest,
			&UserInfoResponse{
				StatusCode: -1,
				StatusMsg:  "error",
			},
		)

		return
	}

	rpcReq := &user.UserInfoRequest{
		UserId:     userIdQuery,
		FromUserId: uint64(ctx.GetUint("userId")),
	}

	rpcResp, err := rpc.UserService.UserInfo(context.Background(), rpcReq)

	if err != nil || rpcResp.StatusCode != 0 {
		ctx.JSON(
			http.StatusInternalServerError,
			&UserInfoResponse{
				StatusCode: -1,
				StatusMsg:  "error",
			},
		)

		return
	}

	ctx.JSON(
		http.StatusOK,
		&UserInfoResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			User: UserInfo{
				Id:            uint(rpcResp.UserId),
				Name:          rpcResp.Username,
				FollowCount:   uint(rpcResp.FollowCount),
				FollowerCount: uint(rpcResp.FollowerCount),
				IsFollow:      rpcResp.IsFollow,
			},
		},
	)
}
