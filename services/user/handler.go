package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"
	"xiaodouyin/kitex_gen/user"
	"xiaodouyin/services/user/db"
)

type UserServiceImpl struct{}

// 验证用户名密码
func (s *UserServiceImpl) AuthUser(ctx context.Context, req *user.AuthRequest) (res *user.AuthResp, err error) {
	res = new(user.AuthResp)
	// 对密码进行 sha256 散列
	password := sha256.Sum256([]byte(req.Password))
	passwordString := fmt.Sprintf("%x", password)

	// 查询用户信息
	userInfo, err := db.QueryUserByUsername(ctx, req.Username)

	// 数据库查询失败
	if err != nil {
		res.StatusCode = -1
		res.StatusMessage = "authenticate failed"
		return
	}

	// 比较密码
	if strings.Compare(passwordString, userInfo.Password) != 0 {
		res.StatusCode = -1
		res.StatusMessage = "authenticate failed"
		return
	}

	res.StatusCode = 0
	res.StatusMessage = "success"
	res.UserId = uint64(userInfo.ID)

	return
}

func (*UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (resp *user.CreateUserResponse, err error) {
	resp = new(user.CreateUserResponse)

	// sha256
	password := sha256.Sum256([]byte(req.Password))
	passwordString := fmt.Sprintf("%x", password)

	user := &db.User{
		Username: req.Username,
		Password: passwordString,
	}

	if err = db.AddUser(ctx, user); err != nil {
		resp.StatusCode = -1
		resp.StatusMessage = "add user failed"
		return resp, nil
	}

	resp.StatusCode = 0
	resp.StatusMessage = "success"
	resp.UserId = uint64(user.ID)

	return
}

func (*UserServiceImpl) UserInfo(ctx context.Context, req *user.UserInfoRequest) (resq *user.UserInfoResponse, err error) {
	resq = new(user.UserInfoResponse)

	userId := req.UserId

	user, err := db.Query(ctx, userId)

	if err != nil {
		resq.StatusCode = -1
		resq.StatusMsg = "error"

		return resq, nil
	}

	resq.StatusCode = 0
	resq.StatusMsg = "success"

	resq.UserId = uint64(user.ID)
	resq.Username = user.Username
	resq.FollowCount = uint64(user.FollowCount)
	resq.FollowerCount = uint64(user.FollowerCount)
	resq.IsFollow = false // todo

	return resq, nil

}
