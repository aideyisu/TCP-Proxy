package tcpmax

import (
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
)

var SCGrpcAddr string
var GrpcTimeout time.Duration
var ConfigFile string
var ProxyGrpcPortStart int

// db操作变量
var users *Users
var listenaddrs *ListenAddrs
var proxymaps *ProxyMaps
var whitelists *WhiteLists
var connections *Connections
var attackconnections *AttackConnections

var InterfaceIns *SubInterfaces

var AttackCon *AttackConnections
var WLS *WhiteLists
var Conn *Connections

var PGPMS *ProxyGrpcPortMaps

var (
    // Version should be updated by hand at each release
    Version = "0.0.2"

    //will be overwritten automatically by the build system
    GitCommit string
    GoVersion string
    BuildTime string
)

func InitDB(dbfile string) {
    _ = InitUsers(dbfile)
    _ = InitListenAddrTable(dbfile)
    WLS = InitWhiteListTable(dbfile)
    _ = InitProxyMapTable(dbfile)
    Conn = InitConnectionTable(dbfile)
    AttackCon = InitAttackConnectionTable(dbfile)
}

func InitInterface(ifname string, dbfile string) {
    InterfaceIns = &SubInterfaces{}
    InterfaceIns.Init(ifname, dbfile)
}

func InitGrpcPortMaps(dbfile string) {
    PGPMS = InitProxyGrpcPort(dbfile)
}

// FullVersion formats the version to be printed
func FullVersion() string {
    return fmt.Sprintf("Git Version: %6s \nGit commit: %6s \nGo version: %6s \nBuild time: %6s \n",
        Version, GitCommit, GoVersion, BuildTime)
}

// GetVersion 获取当前程序版本
func GetVersion(c *gin.Context) {
    c.JSON(200, gin.H{
        "status_code": 200,
        "git_version": Version,
        "git_commit":  GitCommit,
        "go_version":  GoVersion,
        "build_time":  BuildTime,
    })
}
