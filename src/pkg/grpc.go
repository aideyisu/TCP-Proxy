package tcpmax

import (
    "context"
    "errors"
    "fmt"
    "log"
    "net"
    "os/exec"
    "strconv"
    "strings"
    "time"

    pb "src/pb"

    grpc "google.golang.org/grpc"
)

const AddSuccess = "配置添加成功."
const DeleteSuccess = "配置删除成功."
const ModifySuccess = "配置修改成功."
const Status404 = "没有找到配置."

type server struct {
    pb.UnimplementedServiceControlServer
}

func startProxyServer(configFile string, listenip string, listenport int, guardip string, guardport int) error {

    gpp, err := PGPMS.Add(listenip, listenport)
    if err != nil {
        return err
    }
    if gpp == nil {
        return errors.New("Proxy grpc port map not found, can't start proxy server.")
    }
    grpcport := ProxyGrpcPortStart + gpp.GrpcPort.OffSet
    cmd := exec.Command("./tcpmax", "proxyserver",
        "--config", ConfigFile,
        "--listenport", strconv.Itoa(listenport),
        "--listenip", listenip,
        "--grpcport", strconv.Itoa(grpcport),
        "--guardip", guardip,
        "--guardport", strconv.Itoa(guardport))
    err = cmd.Start()
    if err == nil {
        gpp.Cmd = cmd
    }
    return err
}

func stopProxyServer(listenip string, listenport int) (*pb.ProxyReplyOnly, error) {
    gpp := PGPMS.GetProcess(listenip, listenport)
    if gpp == nil {
        return &pb.ProxyReplyOnly{Code: "404", Message: Status404}, nil
    }
    grpcport := ProxyGrpcPortStart + gpp.GrpcPort.OffSet
    addr := fmt.Sprintf("127.0.0.1:%d", grpcport)
    conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(GrpcTimeout))
    if err != nil {
        log.Printf("did not connect: %v", err)
        return nil, err
    }
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
    defer cancel()
    cn := pb.NewProxyClient(conn)
    rep, err := cn.Stop(ctx, &pb.ProxyStopRequest{})
    if rep != nil && rep.Code == "200" {
        if gpp.Cmd != nil {
            go gpp.Cmd.Wait()
        }
        PGPMS.Del(listenip, listenport)
        PGPMS.GrpcPortIDMap.Clear(gpp.GrpcPort.OffSet)
    }
    return rep, err
}

// DealwithListenAddrAdd 新增监听的实际处理
func DealwithListenAddrAdd(IP string, Createor string) (int32, string) {
    if _, err := InterfaceIns.Add(IP); err != nil {
        log.Printf("Add network interface (%s,%s)failed:  %s", IP, Createor, err.Error())
        return checkErr(err)
    }
    if _, err := listenaddrs.Add(IP, Createor); err != nil {
        log.Printf("Add listen address (%s,%s) failed: %s", IP, Createor, err.Error())
        return checkErr(err)
    }

    ProxyMapNum := proxymaps.Count(0, "", 0, "", 0, 0)
    for i := 0; i < ProxyMapNum; i++ {
        proxymap := proxymaps.Find(0, "", 0, "", 0, 0, i+1, 1)
        listenPort := int(proxymap[0].ListenPort)
        port := int(proxymap[0].Port)

        if err := startProxyServer(ConfigFile, IP, listenPort, proxymap[0].IP, port); err != nil {
            log.Printf("Start proxy server (%s,%d,%s,%d) failed: %s", IP, listenPort, proxymap[0].IP, port, err.Error())
            return checkErr(err)
        }
    }
    return 200, AddSuccess
}

// DealwithListenAddrDel 删除监听IP的实际处理
func DealwithListenAddrDel(IP string) (int32, string) {
    if la := listenaddrs.Get(IP); la == nil {
        log.Printf("Get listen address %s failed: %s", IP, Status404)
        return 404, Status404
    }

    ProxyMapNum := proxymaps.Count(0, "", 0, "", 0, 0)
    for i := 0; i < ProxyMapNum; i++ {
        proxymap := proxymaps.Find(0, "", 0, "", 0, 0, i+1, 1)
        _, err := stopProxyServer(IP, int(proxymap[0].ListenPort))
        if err != nil {
            log.Printf("Stop proxy server (%s,%d) failed: %s", IP, proxymap[0].ListenPort, err.Error())
            return checkErr(err)
        }
    }
    if err := listenaddrs.Del(IP); err != nil {
        log.Printf("Delete listen address %s failed: %s.", IP, err.Error())
    }
    InterfaceIns.Del(IP)
    return 200, DeleteSuccess
}

// DealwithListenAddrModify 修改监听IP的配置处理
func DealwithListenAddrModify(id int64, ip string) (int32, string) {
    org := listenaddrs.GetID(id)
    if org == nil {
        log.Printf("Get listen address (%d,%s) failed: %s", id, ip, Status404)
        return 404, Status404
    }

    if _, err := InterfaceIns.Add(ip); err != nil {
        log.Printf("Add network interface (%d,%s) failed: %s", id, ip, err.Error())
        return checkErr(err)
    }

    proxyMapNum := proxymaps.Count(0, "", 0, "", 0, 0)
    for i := 0; i < proxyMapNum; i++ {
        proxymap := proxymaps.Find(0, "", 0, "", 0, 0, i+1, 1)
        _, err := stopProxyServer(org.IP, int(proxymap[0].ListenPort))
        if err != nil {
            log.Printf("Stop proxy server (%d,%s) failed: %s", id, ip, err.Error())
            return checkErr(err)
        }
        if err = startProxyServer(ConfigFile, ip, int(proxymap[0].ListenPort), proxymap[0].IP, int(proxymap[0].Port)); err != nil {
            log.Printf("Start proxy server (%s,%d,%s,%d) failed: %s", ip, proxymap[0].ListenPort, proxymap[0].IP, proxymap[0].Port, err.Error())
            return checkErr(err)
        }
    }
    InterfaceIns.Del(org.IP)
    if _, err := listenaddrs.Modify(id, ip); err != nil {
        log.Printf("Modify listen address (%d,%s) failed: %s", id, ip, err.Error())
        return checkErr(err)
    }
    return 200, ModifySuccess
}

// 新增监听
func (s *server) ListenAddrAdd(ctx context.Context, in *pb.ListenAddrAddRequest) (*pb.ServiceCntrolReplyOnly, error) {
    log.Printf("Add listen address: (%v,%v)", in.GetIP(), in.GetCreateor())
    if IPConflictStatus := ListenIPCheck(in.GetIP()); IPConflictStatus != "" {
        return &pb.ServiceCntrolReplyOnly{StatusCode: 403, Reason: IPConflictStatus}, nil
    }

    statusCode, Reason := DealwithListenAddrAdd(in.GetIP(), in.GetCreateor())

    return &pb.ServiceCntrolReplyOnly{StatusCode: statusCode, Reason: Reason}, nil
}

// 删除监听
func (s *server) ListenAddrDel(ctx context.Context, in *pb.ListenAddrDelRequest) (*pb.ServiceCntrolReplyOnly, error) {
    log.Printf("Delete listen address: %v", in.GetIP())
    statusCode, Reason := DealwithListenAddrDel(in.GetIP())

    return &pb.ServiceCntrolReplyOnly{StatusCode: statusCode, Reason: Reason}, nil
}

// 修改监听
func (s *server) ListenAddrModify(ctx context.Context, in *pb.ListenAddrModifyRequest) (*pb.ServiceCntrolReplyOnly, error) {
    log.Printf("Modify listen address: (%v,%v)", in.GetID(), in.GetIP())
    if IPConflictStatus := ListenIPCheck(in.GetIP()); IPConflictStatus != "" {
        return &pb.ServiceCntrolReplyOnly{StatusCode: 403, Reason: IPConflictStatus}, nil
    }
    StatusCode, Reason := DealwithListenAddrModify(in.GetID(), in.GetIP())

    return &pb.ServiceCntrolReplyOnly{StatusCode: StatusCode, Reason: Reason}, nil
}

// DealwithProxyMapAdd 新增回连映射的基础处理
func DealwithProxyMapAdd(listenPort uint32, guardIP string, guardPort uint32, creator string) (int32, string) {
    if _, err := proxymaps.Add(int(listenPort), guardIP, int(guardPort), creator); err != nil {
        log.Printf("Add proxy map (%d,%s,%d,%s) failed: %s", listenPort, guardIP, guardPort, creator, err.Error())
        return checkErr(err)
    }

    ListenAddrNum := listenaddrs.Count("", "", 0, 0)
    for i := 0; i < ListenAddrNum; i++ {
        las := listenaddrs.Find("", "", 0, 0, i+1, 1)
        err := startProxyServer(ConfigFile, las[0].IP, int(listenPort), guardIP, int(guardPort))
        if err != nil {
            log.Printf("Start proxy server (%d,%s,%d,%s) failed: %s", listenPort, guardIP, guardPort, creator, err.Error())
            return checkErr(err)
        }
    }
    return 200, AddSuccess
}

// DealwithProxyMapDel 删除回连映射的基础处理
func DealwithProxyMapDel(DelListenPort uint32, GuardIP string, GuardPort uint32) (int32, string) {
    if la := proxymaps.Get(int(DelListenPort), GuardIP, int(GuardPort)); la == nil {
        log.Printf("Get proxy map (%d,%s,%d) failed: %s", DelListenPort, GuardIP, GuardPort, Status404)
        return 404, Status404
    }

    ListenAddrNum := listenaddrs.Count("", "", 0, 0)
    for i := 0; i < ListenAddrNum; i++ {
        las := listenaddrs.Find("", "", 0, 0, i+1, 1)
        _, err := stopProxyServer(las[0].IP, int(DelListenPort))
        if err != nil {
            log.Printf("Stop proxy server (%s,%d) failed: %s", GuardIP, GuardPort, err.Error())
            return checkErr(err)
        }
    }

    if err := proxymaps.Del(int(DelListenPort), GuardIP, int(GuardPort)); err != nil {
        log.Printf("Delete proxy map (%d,%s,%d) failed: %s", DelListenPort, GuardIP, GuardPort, Status404)
        return 404, Status404
    }

    return 200, DeleteSuccess
}

// DealwithProxyMapModify  修改回连映射配置
func DealwithProxyMapModify(id int64, listenPort int, ip string, port int) (int32, string) {
    Org := proxymaps.GetID(id)
    if Org == nil {
        log.Printf("Get proxy map (%d,%d,%s,%d) failed: %s", id, id, ip, port, Status404)
        return 404, Status404
    }
    ListenAddrNum := listenaddrs.Count("", "", 0, 0)
    for i := 0; i < ListenAddrNum; i++ {
        las := listenaddrs.Find("", "", 0, 0, i+1, 1)
        _, err := stopProxyServer(las[0].IP, int(Org.ListenPort))
        if err != nil {
            log.Printf("Stop proxy server (%s,%d) failed: %s", las[0].IP, Org.ListenPort, err.Error())
            return checkErr(err)
        }
        err = startProxyServer(ConfigFile, las[0].IP, listenPort, ip, port)
        if err != nil {
            log.Printf("Start proxy server (%s,%d,%s,%d) failed: %s", las[0].IP, listenPort, ip, port, err.Error())
            return checkErr(err)
        }
    }
    if _, err := proxymaps.Modify(id, listenPort, ip, port); err != nil {
        log.Printf("Modify proxy map (%d,%d,%s,%d) failed: %s", id, id, ip, port, err.Error())
        return checkErr(err)
    }

    return 200, ModifySuccess
}

func CheckPortAvailability(listenport uint32) error {
    ListenPort := strconv.Itoa(int(listenport))
    addr, err := net.ResolveTCPAddr("tcp", "localhost:"+ListenPort)
    if err != nil {
        return err
    }

    l, err := net.ListenTCP("tcp", addr)
    if err != nil {
        return err
    }
    defer l.Close()
    return nil
}

// 回连映射增加
func (s *server) ProxyMapAdd(ctx context.Context, in *pb.ProxyMapAddRequest) (*pb.ServiceCntrolReplyOnly, error) {
    if err := CheckPortAvailability(in.GetListenPort()); err != nil {
        log.Printf("Failed add proxy map: (%v,%v). %v", in.GetListenPort(), in.GetIP(), err)
        statusCode, Reason := checkErr(err)
        return &pb.ServiceCntrolReplyOnly{StatusCode: statusCode, Reason: Reason}, err
    }
    log.Printf("Add proxy map: (%v,%v)", in.GetListenPort(), in.GetIP())
    statusCode, Reason := DealwithProxyMapAdd(in.GetListenPort(), in.GetIP(), in.GetPort(), in.GetCreator())

    return &pb.ServiceCntrolReplyOnly{StatusCode: statusCode, Reason: Reason}, nil
}

// 回连映射删除
func (s *server) ProxyMapDel(ctx context.Context, in *pb.ProxyMapDelRequest) (*pb.ServiceCntrolReplyOnly, error) {
    log.Printf("Delete proxy map: (%v,%v,%v)", in.GetListenPort(), in.GetIP(), in.GetPort())
    StatusCode, Reason := DealwithProxyMapDel(in.GetListenPort(), in.GetIP(), in.GetPort())

    return &pb.ServiceCntrolReplyOnly{StatusCode: StatusCode, Reason: Reason}, nil
}

// 回连映射修改
func (s *server) ProxyMapModify(ctx context.Context, in *pb.ProxyMapModifyRequest) (*pb.ServiceCntrolReplyOnly, error) {
    if err := CheckPortAvailability(in.GetListenPort()); err != nil {
        log.Printf("Failed modify proxy map: (%v,%v). %v", in.GetListenPort(), in.GetIP(), err)
        statusCode, Reason := checkErr(err)
        return &pb.ServiceCntrolReplyOnly{StatusCode: statusCode, Reason: Reason}, err
    }
    log.Printf("Modify proxy map: (%v,%v)", in.GetListenPort(), in.GetIP())
    StatusCode, Reason := DealwithProxyMapModify(in.GetID(), int(in.GetListenPort()), in.GetIP(), int(in.GetPort()))

    return &pb.ServiceCntrolReplyOnly{StatusCode: StatusCode, Reason: Reason}, nil
}

// SrartGrpc 启动grpc服务
func SrartGrpc(port string) {

    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()

    pb.RegisterServiceControlServer(s, &server{})
    s.Serve(lis)
}

func checkErr(err error) (int32, string) {

    if err == nil {
        return 200, "操作成功."
    }

    if strings.Contains(err.Error(), "Access denied") {
        return 403, "无权限访问."
    }

    if strings.Contains(err.Error(), "Error") {
        return 404, Status404
    }

    if strings.Contains(err.Error(), "json: ") {
        return 400, "请求类型错误."
    }

    if strings.Contains(err.Error(), "invalid character") {
        return 400, "无效的字符."
    }

    if strings.Contains(err.Error(), "UNIQUE") {
        return 409, "配置冲突."
    }

    if strings.Contains(err.Error(), "ddress already in use") {
        return 409, "监听端口与本地冲突."
    }

    if strings.Contains(err.Error(), "exit status 1") {
        return 403, "检查网卡配置是否正确."
    }

    log.Printf("当前系统发生未知错误：%s", err.Error())

    return 500, "系统发生未知错误."
}

// ListenIPCheck 检测IP是否与本地IP冲突
func ListenIPCheck(ip string) string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        log.Println(err)
        return err.Error()
    }
    for _, address := range addrs {
        // 检查ip地址判断是否回环地址
        // if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
        if ipnet, ok := address.(*net.IPNet); ok {
            if ipnet.IP.String() == ip {
                return "与本地IP冲突"
            }
        }
    }
    return ""
}
