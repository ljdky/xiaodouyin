package db

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// User 结构体
type User struct {
	gorm.Model

	Username      string `gorm:"not null;unique;index:idx_username_password,expression:username(20)"` // 用户名
	Password      string `gorm:"not null;index:idx_username_password,expression:password(20)"`        // 密码
	FollowCount   uint   `gorm:"not null;default:0"`                                                  // 关注的人数
	FollowerCount uint   `gorm:"not null;default:0"`                                                  // 粉丝数
}

type Follow struct {
	CreatedAt time.Time `gorm:"not null"`
	DeletedAt gorm.DeletedAt

	Follower uint `gorm:"not null"` // 关注者 ID
	Followee uint `gorm:"not null"` // 被关注者 ID
}

func Query(ctx context.Context, userId uint64) (*User, error) {
	user := new(User)
	if err := DB.Take(&user, userId).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func QueryUserByUsername(ctx context.Context, username string) (*User, error) {

	res := new(User)
	if err := DB.WithContext(ctx).Where("username = ?", username).Take(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func AddUser(ctx context.Context, user *User) error {
	if err := DB.Select("username", "password").Create(user).Error; err != nil {
		return err
	}

	return nil
}
