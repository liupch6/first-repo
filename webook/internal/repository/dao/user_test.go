package dao

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name string
		// 为什么不用 ctrl ？这里不是 gomock， 这是 sqlmock
		mock func(t *testing.T) *sql.DB

		ctx context.Context
		u   User

		expectedErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				// 这边预期的是正则表达式
				// 这个写法的意思就是只要是 INSERT 到 users 的语句就行
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnResult(sqlmock.NewResult(3, 1))
				require.NoError(t, err)
				return mockDB
			},
			ctx: context.Background(),
			u: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
			},
			expectedErr: nil,
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				// 这边预期的是正则表达式
				// 这个写法的意思就是只要是 INSERT 到 users 的语句就行
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnError(&mysql.MySQLError{Number: 1062})
				require.NoError(t, err)
				return mockDB
			},
			ctx: context.Background(),
			u: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
			},
			expectedErr: ErrUserDuplicate,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				// 这边预期的是正则表达式
				// 这个写法的意思就是只要是 INSERT 到 users 的语句就行
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnError(errors.New("mock db error"))
				require.NoError(t, err)
				return mockDB
			},
			ctx: context.Background(),
			u: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
			},
			expectedErr: errors.New("mock db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := gorm.Open(gormMysql.New(gormMysql.Config{
				Conn:                      tc.mock(t),
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				SkipDefaultTransaction: true,
				DisableAutomaticPing:   true,
			})
			d := NewUserDAO(db)
			err = d.Insert(tc.ctx, tc.u)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
