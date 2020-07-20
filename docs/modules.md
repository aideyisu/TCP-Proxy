# go module使用说明

- TODO 举例说明

Go111Modules启动项目流程：
    前言:Go modules 是 Go 语言中正式官宣的项目依赖解决方案，Go modules（前身为vgo）于 Go1.11 正式发布，在 Go1.14 已经准备好，并且可以用在生产上（ready for production）了，Go官方也鼓励所有用户从其他依赖项管理工具迁移到 Go modules
    而 Go1.14，在近期也终于正式发布，Go 官方亲自 “喊” 你来用

    在项目中的基本使用:
         go mod init 生成mod文件
         go build main.go 编译生成可执行文件

    注意事项:
        import的路径问题，gomod与gopath有所不同
        假设当前工程路径workspace/testmod;gopath下引用"api/supply/location"
        gomod下引用"testmod/api/supply/location"
