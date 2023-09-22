//go:build !k8s

// aaa go:build dev
// aaa go:build test
// aaa go:build e2e

package config

var Config = config{
	DB: DBConfig{
		DSN: "localhost:13316",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
