.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

APP = kube-fake-ssl

# go env -w GO111MODULE=on
# go env -w GOPROXY=http://mvn.res.local/repository/go,direct
# go env -w GOPROXY=https://proxy.golang.com.cn,direct
# go env -w GOPROXY=https://goproxy.cn,direct
# go env -w GOSUMDB=sum.golang.google.cn
# go env -w GOSUMDB=sum.golang.org
# go env -w GOSUMDB=off

# 初始化mod
init:
	go mod init github.com/suisrc/${APP}

# 修正依赖
tidy:
	go mod tidy

run:
	go run main.go

# 打包应用
build:
	go build -ldflags "-w -s" -o $(SERVER_BIN) .

deploy:
	git checkout deploy-auto && git merge master && git push && git checkout master
	
github:
	git checkout plus && git merge master && git push -u github && git checkout master