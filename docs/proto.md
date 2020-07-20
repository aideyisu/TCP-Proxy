# proto使用说明

- TODO 举例说明

proto自动生成代码：
    1:安装protoc插件
        go get -u github.com/golang/protobuf/protoc-gen-go
        或者
        brew install protobuf

        编译器插件protoc-gen-go将会安装在中$GOBIN，默认为$GOPATH/bin。PATH协议编译器protoc 必须在您的目录中才能找到它。
    2:配置环境变量
        export PATH=$PATH:$GOPATH/bin
    3:自动生成代码
        在本项目中,在src路径下使用指令
        protoc -I services/ services/services.proto --go_out=plugins=grpc:services
        即可对services中的proto文件进行代码生成,自动生成设计文档中的代码格式
