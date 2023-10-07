package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserDuplicate = errors.New("邮箱重复")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string

	// 唯一索引允许有多个空值，但是不能有多个 ""
	Phone sql.NullString `gorm:"unique"`

	// 创建和更新时间：毫秒数
	Ctime int64
	Utime int64
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (ud *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := ud.db.WithContext(ctx).Create(&u).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueIndexErrNo uint16 = 1062
		if mysqlErr.Number == uniqueIndexErrNo {
			// 邮箱冲突或者手机号码冲突
			return ErrUserDuplicate
		}
	}
	return err
}

func (ud *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).First(&u, "email = ?", email).Error
	// err := ud.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (ud *UserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).First(&u, "phone = ?", phone).Error
	return u, err
}

func (ud *UserDao) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).First(&u, "id = ?", id).Error
	return u, err
}
