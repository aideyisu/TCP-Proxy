package main

import (
    //"fmt"
    "flag"
    "fmt"
    "log"
    "time"

    //"net/http"
    "io"
    "os"
    tcpmax "src/pkg"
    services "src/services"

    "github.com/fsnotify/fsnotify"
    "github.com/spf13/viper"
)

type serverOption struct {
    Config string
}

type apiServerOption struct {
    Config string
}

type proxyServerOption struct {
    Config     string
    ListenIP   string
    ListenPort int
    GrpcPort   int
    GuardIP    string
    GuardPort  int
}

type userAddOption struct {
    UserName string
    Desc     string
    Config   string
}

type userDelOption struct {
    UserName string
    Config   string
}

type userInfoOption struct {
    UserName string
    Config   string
}

type userUpdateOption struct {
    UserName string
    Key      bool
    Desc     string
    Config   string
}

func apiServer(confFile string) {
    readConfig(confFile)
    logFile, err := os.OpenFile(viper.Get("log.tcpmax").(string), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("Failed to open log file:", err)
    }

    log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
    log.SetOutput(io.MultiWriter(logFile))

    tcpmax.InitDB(viper.Get("dbfile").(string))

    r, authMiddleware := services.CreateRouter(viper.Get("log.gin").(string),
        tcpmax.PayloadFuncHandler,
        tcpmax.AuthenticatorHandler,
        tcpmax.IdentityHandler,
        tcpmax.AuthorizatorHandler)
    auth := r.Group("/")
    auth.Use(authMiddleware.MiddlewareFunc())
    {
        // 监听地址
        auth.GET("/listenaddr", tcpmax.ListListenaddrHandler)
        auth.PUT("/listenaddr", tcpmax.ModifyListenaddrHandler)
        auth.POST("/listenaddr", tcpmax.AddListenaddrHandler)
        auth.DELETE("/listenaddr", tcpmax.DelListenaddrHandler)

        // 回连映射
        auth.GET("/proxymap", tcpmax.ListProxyMapHandler)
        auth.PUT("/proxymap", tcpmax.ModifyProxyMapHandler)
        auth.POST("/proxymap", tcpmax.AddPorxyMapHandler)
        auth.DELETE("/proxymap", tcpmax.DelProxyMapHandler)

        // IP白名单
        auth.GET("/whitelist", tcpmax.ListWhiteListHandler)
        auth.PUT("/whitelist", tcpmax.ModifyWhiteListHandler)
        auth.POST("/whitelist", tcpmax.AddWhiteListHandler)
        auth.DELETE("/whitelist", tcpmax.DelWhiteListHandler)
        auth.PATCH("/whitelist", tcpmax.ChangeWhiteStatusHandler)

        // 连接管理
        auth.GET("/connection", tcpmax.ListOnlineConnectionsHandler)
        auth.POST("/connection", tcpmax.DisConnectionHandler)

        // 日志
        auth.GET("/logs/connection/authorized", tcpmax.ListConnectionsLogsHandler)
        auth.GET("/logs/connection/unauthorized", tcpmax.ListAtkConnectionsHandler)

        // 版本信息
        auth.GET("/version", tcpmax.GetVersion)
    }
    tcpmax.SCGrpcAddr = viper.Get("service.scgrpc").(string)
    tcpmax.GrpcTimeout = time.Millisecond * time.Duration(viper.Get("service.grpctimeout").(int))

    tcpmax.ProxyGrpcPortStart = viper.Get("service.proxygrpcportstart").(int)
    tcpmax.InitGrpcPortMaps(viper.Get("dbfile").(string))

    services.RunServer(viper.Get("service.apiserver").(string), r)
}

func proxyServer(confFile string, listenIP string, listenPort int, grpcport int, guardIP string, guardPort int) {
    readConfig(confFile)
    logFile, err := os.OpenFile(viper.Get("log.proxy").(string), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("Failed to open log file:", err)
    }

    log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
    log.SetOutput(io.MultiWriter(logFile))
    tcpmax.InitDB(viper.Get("dbfile").(string))
    tcpmax.InitGrpcPortMaps(viper.Get("dbfile").(string))
    services.ProxyServer(listenIP, listenPort, grpcport, guardIP, guardPort)
}

func serviceControlServer(confFile string) {
    readConfig(confFile)
    tcpmax.ConfigFile = confFile
    logFile, err := os.OpenFile(viper.Get("log.srvctl").(string), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("Failed to open log file:", err)
    }

    log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
    log.SetOutput(io.MultiWriter(logFile))

    // 初始化网卡管理
    tcpmax.InitInterface(viper.Get("service.ifname").(string), viper.Get("dbfile").(string))
    tcpmax.InitDB(viper.Get("dbfile").(string))
    tcpmax.ProxyGrpcPortStart = viper.Get("service.proxygrpcportstart").(int)
    tcpmax.InitGrpcPortMaps(viper.Get("dbfile").(string))
    tcpmax.RecoverProxyProcess(confFile)
    tcpmax.SrartGrpc(viper.Get("service.scgrpc").(string))
}

func onConfigChange(e fsnotify.Event) {
    log.Printf("Config changed: %s", e.Name)
}

func readConfig(file string) error {
    viper.SetConfigName(file)
    viper.AddConfigPath(".")
    viper.SetConfigType("yaml")
    err := viper.ReadInConfig()
    if err != nil {
        log.Printf("Config file error: %s", err)
        return err
    }

    viper.WatchConfig()
    viper.OnConfigChange(onConfigChange)
    return nil
}

func showUser(u *tcpmax.User) {
    if u == nil {
        return
    }
    fmt.Println("username\tkey\t\t\tcreatetime\t\tdescription")
    timeNow := time.Unix(u.CreateTime, 0)
    timeString := timeNow.Format("2006-01-02 15:04:05")
    fmt.Printf("%s\t\t%s\t%s\t%s\n",
        u.UserName, u.Key, timeString, u.Desc)
}

func showVersion() {
    fmt.Println(tcpmax.FullVersion())
}

func main() {
    serverOpt := serverOption{}
    serverCmd := flag.NewFlagSet("servicectlserver", flag.ExitOnError)
    serverCmd.StringVar(&serverOpt.Config, "config", "config.yaml", "config file")

    apiServerOpt := apiServerOption{}
    apiServerCmd := flag.NewFlagSet("apiserver", flag.ExitOnError)
    apiServerCmd.StringVar(&apiServerOpt.Config, "config", "config.yaml", "config file")

    proxyServerOpt := proxyServerOption{}
    proxyServerCmd := flag.NewFlagSet("proxyserver", flag.ExitOnError)
    proxyServerCmd.StringVar(&proxyServerOpt.Config, "config", "config.yaml", "config file")
    proxyServerCmd.StringVar(&proxyServerOpt.ListenIP, "listenip", "127.0.0.1", "listen ip")
    proxyServerCmd.IntVar(&proxyServerOpt.ListenPort, "listenport", 5000, "listen port")
    proxyServerCmd.IntVar(&proxyServerOpt.GrpcPort, "grpcport", 60000, "grpc port")
    proxyServerCmd.StringVar(&proxyServerOpt.GuardIP, "guardip", "127.0.0.1", "guard ip")
    proxyServerCmd.IntVar(&proxyServerOpt.GuardPort, "guardport", 5500, "listen ip")

    userAddOpt := userAddOption{}
    useraddCmd := flag.NewFlagSet("useradd", flag.ExitOnError)
    useraddCmd.StringVar(&userAddOpt.UserName, "username", "", "user name")
    useraddCmd.StringVar(&userAddOpt.Desc, "desc", "", "description")
    useraddCmd.StringVar(&userAddOpt.Config, "config", "config.yaml", "config file")

    userDelOpt := userDelOption{}
    userdelCmd := flag.NewFlagSet("userdel", flag.ExitOnError)
    userdelCmd.StringVar(&userDelOpt.UserName, "username", "", "user name")
    userdelCmd.StringVar(&userDelOpt.Config, "config", "config.yaml", "config file")

    userInfoOpt := userInfoOption{}
    userinfoCmd := flag.NewFlagSet("userinfo", flag.ExitOnError)
    userinfoCmd.StringVar(&userInfoOpt.UserName, "username", "", "user name")
    userinfoCmd.StringVar(&userInfoOpt.Config, "config", "config.yaml", "config file")

    userUpdateOpt := userUpdateOption{}
    userupdateCmd := flag.NewFlagSet("userupdate", flag.ExitOnError)
    userupdateCmd.StringVar(&userUpdateOpt.UserName, "username", "", "user name")
    userupdateCmd.StringVar(&userUpdateOpt.Desc, "desc", "", "description")
    userupdateCmd.BoolVar(&userUpdateOpt.Key, "key", false, "key")
    userupdateCmd.StringVar(&userUpdateOpt.Config, "config", "config.yaml", "config file")

    if len(os.Args) < 2 {
        fmt.Println("expected 'servicectlserver', 'apiserver', 'proxyserver', 'version', 'uesradd', 'userdel', 'userupdate' or 'userinfo' subcommands")
        os.Exit(1)
    }

    switch os.Args[1] {

    case "servicectlserver":
        serverCmd.Parse(os.Args[2:])
        serviceControlServer(serverOpt.Config)

    case "proxyserver":
        proxyServerCmd.Parse(os.Args[2:])
        proxyServer(proxyServerOpt.Config, proxyServerOpt.ListenIP, proxyServerOpt.ListenPort, proxyServerOpt.GrpcPort, proxyServerOpt.GuardIP, proxyServerOpt.GuardPort)

    case "apiserver":
        apiServerCmd.Parse(os.Args[2:])
        apiServer(apiServerOpt.Config)

    case "useradd":
        useraddCmd.Parse(os.Args[2:])
        readConfig(userAddOpt.Config)
        users := tcpmax.InitUsers(viper.Get("dbfile").(string))
        showUser(users.Add(userAddOpt.UserName, userAddOpt.Desc))

    case "userdel":
        userdelCmd.Parse(os.Args[2:])
        readConfig(userDelOpt.Config)
        users := tcpmax.InitUsers(viper.Get("dbfile").(string))
        users.Del(userDelOpt.UserName)

    case "userupdate":
        userupdateCmd.Parse(os.Args[2:])
        readConfig(userUpdateOpt.Config)
        users := tcpmax.InitUsers(viper.Get("dbfile").(string))
        u := &tcpmax.User{}
        fmt.Println(userUpdateOpt)
        if userUpdateOpt.Desc != "" {
            u = users.SetDesc(userUpdateOpt.UserName, userUpdateOpt.Desc)
        }
        if userUpdateOpt.Key == true {
            u = users.UpdateKey(userUpdateOpt.UserName)
        }
        showUser(u)

    case "userinfo":
        userinfoCmd.Parse(os.Args[2:])
        readConfig(userInfoOpt.Config)
        users := tcpmax.InitUsers(viper.Get("dbfile").(string))
        fmt.Println(userInfoOpt)
        showUser(users.Get(userInfoOpt.UserName))

    case "version":
        showVersion()

    default:
        fmt.Println("expected 'servicectlserver', 'apiserver', 'proxyserver', 'version', 'uesradd', 'userdel', 'userupdate' or 'userinfo' subcommands")
        os.Exit(1)
    }
}
