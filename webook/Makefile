.PHONY: docker
docker:
	# 把上次编译的东西删除
	@rm webook || true
	@docker rmi -f liupch/webook:v0.0.1
	# 运行一下 go mod tidy，防止 go.sum 文件不对，编译失败
	@go mod tidy
	# 指定编译在 ARM 架构的 linux 操作系统上的
	@GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -o webook .
	@docker build -t liupch/webook:v0.0.1 .