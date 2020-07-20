package services

import (
    "fmt"
    "io"
    "log"
    "net"
    "runtime/debug"
    tcpmax "src/pkg"
    "strconv"
    "strings"
    "sync"
    "time"
)

// TCP 配置
type TCP struct {
    listenIP   string
    listenPort int
    guardIP    string
    guardPort  int
}

// ServerChannel 支持信件服务端通道
type ServerChannel struct {
    ip               string
    port             int
    Listener         net.Listener
    errAcceptHandler func(err error)
}

var tcpSC ServerChannel

// NewServerChannel 为了监听程序支撑
func NewServerChannel(ip string, port int) ServerChannel {
    return ServerChannel{
        ip:   ip,
        port: port,
        errAcceptHandler: func(err error) {
            log.Printf("accept error , ERR:%s", err)
        },
    }
}

// 建立监听
func Trylisten(listenIP string, listenPort int, guardIP string, guardPort int) {
    tcp := TCP{
        listenIP,
        listenPort,
        guardIP,
        guardPort,
    }

    err := tcp.Start()
    if err != nil {
        log.Printf("Listen %s:%d failed: %s.", listenIP, listenPort, err.Error())
    }
}

// ProxyHandler 代理处理流程
func (tcp *TCP) ProxyHandler(inConn net.Conn) {
    //
    inAddr := inConn.RemoteAddr().String()
    log.Printf("in connection %s.\n", inAddr)
    inLocalAddr := inConn.LocalAddr().String()
    SrcIP := strings.Split(inAddr, ":")[0]
    SrcPort, _ := strconv.Atoi(strings.Split(inAddr, ":")[1])
    ListenIP := strings.Split(inLocalAddr, ":")[0]
    ListenPort, _ := strconv.Atoi(strings.Split(inLocalAddr, ":")[1])
    wl := tcpmax.WLS.Get(SrcIP)
    if wl == nil || wl.Status != "启用" {
        inConn.Close()
        tcpmax.AttackCon.Add("TCP", SrcIP, tcp.listenIP, SrcPort, tcp.listenPort)
        log.Printf("white ip:%s not found.\n", SrcIP)
        return
    }
    outConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", tcp.guardIP, tcp.guardPort))
    if err != nil {
        log.Printf("Failed to establish connection with the guard(%s:%d): %s.", tcp.guardIP, tcp.guardPort, err.Error())
        return
    }
    outAddr := outConn.RemoteAddr().String()
    outLocalAddr := outConn.LocalAddr().String()

    // DstIP := strings.Split(outAddr, ":")[0]
    // DstPort, _ := strconv.Atoi(strings.Split(outAddr, ":")[1])
    // 添加流量更新统计
    rep := tcpmax.Conn.Add("TCP", SrcIP, ListenIP, SrcPort, ListenPort, 0, 0, "连接")
    ID := rep.ID
    key := tcpmax.ProxyCBKey{
        SrcIP:   SrcIP,
        SrcPort: uint32(SrcPort),
        DstIP:   tcp.listenIP,
        DstPort: uint32(tcp.listenPort),
    }
    value := &tcpmax.ProxyCB{
        Key:        key,
        Controlled: inConn,
        Control:    outConn,
        Send:       0,
        Received:   0,
        ID:         ID,
    }
    //调用接口实现，在建立监听的时候增加键值
    tcpmax.Conn.ProxyCBAdd(key, value)

    IoBind(inConn, outConn, ID, func(isSrcErr bool, err error) {
        log.Printf("conn %s - %s - %s -%s (%d)released", inAddr, inLocalAddr, outLocalAddr, outAddr, ID)
        CloseConn(&inConn)
        CloseConn(&outConn)
        tcpmax.Conn.UpdateStatus(ID, "断开")
    }, func(n int, d bool) {})
    log.Printf("conn %s - %s - %s -%s (%d)connected", inAddr, inLocalAddr, outLocalAddr, outAddr, ID)
}

// IoBind IO绑定,略微修改
func IoBind(dst io.ReadWriter, src io.ReadWriter, ID int64, fn func(isSrcErr bool, err error), cfn func(count int, isPositive bool)) {
    var one = &sync.Once{}
    go func() {
        defer func() {
            if e := recover(); e != nil {
                log.Printf("IoBind crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
            }
        }()
        var err error
        var isSrcErr bool

        _, isSrcErr, err = ioCopy(src, dst, ID, "Send", func(c int) {
            cfn(c, false)
        })

        if err != nil {
            one.Do(func() {
                fn(isSrcErr, err)
            })
        }
    }()
    go func() {
        defer func() {
            if e := recover(); e != nil {
                log.Printf("IoBind crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
            }
        }()
        var err error
        var isSrcErr bool

        _, isSrcErr, err = ioCopy(dst, src, ID, "Received", func(c int) {
            cfn(c, false)
        })

        if err != nil {
            one.Do(func() {
                fn(isSrcErr, err)
            })
        }
    }()
}

// DataLog 记录
func DataLog(ID int64, nr int, Type string) {
    Org := tcpmax.Conn.Get(ID)
    if Org != nil && Type == "Received" {
        ReceivedSize := int64(Org.Received + uint64(nr))
        _ = tcpmax.Conn.UpdateReceived(ID, ReceivedSize)
    }
    if Org != nil && Type == "Send" {
        SendSize := int64(Org.Send + uint64(nr))
        _ = tcpmax.Conn.UpdateSend(ID, SendSize)
    }
}

// IOCopy 数据拷贝无修改
func ioCopy(dst io.Writer, src io.Reader, ID int64, OrderType string, fn ...func(count int)) (written int64, isSrcErr bool, err error) {
    buf := make([]byte, 32*1024)
    for {
        nr, er := src.Read(buf)
        go DataLog(ID, nr, OrderType)
        if nr > 0 {
            nw, ew := dst.Write(buf[0:nr])
            if nw > 0 {
                written += int64(nw)
                if len(fn) == 1 {
                    fn[0](nw)
                }
            }
            if ew != nil {
                err = ew
                break
            }
            if nr != nw {
                err = io.ErrShortWrite
                break
            }
        }
        if er != nil {
            err = er
            isSrcErr = true
            break
        }
    }
    return written, isSrcErr, err
}

// ListenTCP 监听TCP链接，等待
func (sc *ServerChannel) ListenTCP(fn func(conn net.Conn)) (err error) {
    listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", sc.ip, sc.port))

    if err == nil {
        sc.Listener = listener

        go func() {
            defer func() {
                if e := recover(); e != nil {
                    log.Printf("ListenTCP crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
                }
            }()
            for {
                var conn net.Conn
                conn, err = sc.Listener.Accept()
                if err == nil {
                    go func() {
                        defer func() {
                            if e := recover(); e != nil {
                                log.Printf("connection handler crashed , err : %s , \ntrace:%s", e, string(debug.Stack()))
                            }
                        }()
                        fn(conn)
                    }()
                } else {
                    checkError(err)
                    break
                }
            }
        }()
    }

    return
}

func (sc *ServerChannel) CloseListener() {
    sc.Listener.Close()
}

// Start 启动代理服务
func (tcp *TCP) Start() (err error) {
    tcpSC := NewServerChannel(tcp.listenIP, tcp.listenPort)
    err = tcpSC.ListenTCP(tcp.ProxyHandler)
    if err != nil {
        return
    }
    log.Printf("TCP proxy on %s:%d", tcp.listenIP, tcp.listenPort)
    return
}

func StopListener() {
    tcpSC.CloseListener()
}

// CloseConn 关闭连接
func CloseConn(conn *net.Conn) {
    if conn != nil && *conn != nil {
        (*conn).SetDeadline(time.Now().Add(time.Millisecond))
        (*conn).Close()
    }
}

func checkError(err error) {
    if err != nil {
        log.Printf("Fatal error %s", err.Error())
        //time.Sleep(time.Second * 5)
        //os.Exit(1)
    }
}
