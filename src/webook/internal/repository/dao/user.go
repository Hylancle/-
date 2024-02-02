package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}
func (dao *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	return dao.db.WithContext(ctx).Create(&u).Error
}

// 数据存储格式
type User struct {
	Id       int64  `grom:"primaryKey,autoIncrement"`
	Email    string `grom:"unique"`
	Password string
	// 时区,UTC0毫秒数
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}
