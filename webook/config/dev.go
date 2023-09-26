//go:build !k8s

// aaa go:build dev
// aaa go:build test
// aaa go:build e2e

package config

var Config = config{
	DB: DBConfig{
		// 本地连接
		DSN: "root:root@tcp(localhost:13316)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
