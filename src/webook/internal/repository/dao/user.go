package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

// 邮箱冲突
var (
	ErrorDuplicateEmail = errors.New("邮箱已经注册过")
	ErrorRecordNotFound = gorm.ErrRecordNotFound
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
	err := dao.db.WithContext(ctx).Create(&u).Error
	// 我断言他是一个mysql，断言成功，进入语句
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 邮箱错误
			return ErrorDuplicateEmail
		}
	}
	return err
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error // 指针传入
	return u, err
}

// 数据存储格式
type User struct {
	Id       int64  `grom:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	// 时区,UTC0毫秒数
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}
