package db

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	CreatedAt time.Time
	DeleteAt  gorm.DeletedAt

	AuthorId uint
	PlayUrl  string
	CoverUrl string

	FavoriteCount uint // 点赞人数
	CommentCount  uint // 评论人数
}

type Favorite struct {
	CreatedAt time.Time
	DeleteAt  gorm.DeletedAt

	UserId  uint // 点赞用户 ID
	VideoId uint // 被点赞视频 ID
}

type Comment struct {
	CreatedAt time.Time
	DeleteAt  gorm.DeletedAt

	UserId  uint   // 评论用户ID
	VideoId uint   // 被评论视频ID
	Content string // 评论内容
}
