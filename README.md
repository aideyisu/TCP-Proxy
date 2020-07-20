# TCP-Proxy

## 简介

  Go语言实现的流量代理.核心转发部分借鉴了gopacket源码.每个流量tcp连接独占线程,通过exec建立线程

后端服务与内部tcp连接之间通过grpc进行通信

通过build.sh可以规定内部版本号

通过Gin以restful格式进行通信

## 目标群体

  如果缺乏go实战基础,可以来看看这个项目.
  
  虽然(๑•̀ㅂ•́)و✧项目不大.但是用到了很多有趣的技术
  
  咳咳.项目里有很多单词拼错QAQ,但是隔段时间忘记在哪里了....如果看到还请提醒我一下Otz
  
## 技术特点
  
  1 Web框架Gin,内部通信grpc,sqlite数据库,log日志记录
  
  2 核心流量转发仿照gopacket,未使用连接池.
  
  3 build.sh 可以添加版本号
  
  4 exec添加子进程
  
  5 config.yaml 配置管理

## webapi

### 运行
```
cd src
go run services/web.go 
```

### 编译
```
cd src
sh build.sh
```

### 运行

#### 启动服务
```
# 启动代理进程
./main  proxyserver --listenip 127.0.0.1 --listenport 6000 --guardip 127.0.0.1 --guardport 8000

# 启动服务控制进程
./main server

# 启动apiserver
./main apiserver
```

#### 用户管理
```
# 添加用户
./main useradd --username username

# 删除用户
./main userdel --username username

# 更新用户信息
./main userupdate --desc "爱的伊苏到此一游" --key

# 查看用户信息
./main userinfo --username username
```

## 注意事项

先添加用户再操作!
