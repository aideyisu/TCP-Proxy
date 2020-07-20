package services

import (
    "context"
    "fmt"
    "log"
    "net"
    pb "src/pb"
    tcpmax "src/pkg"

    "google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedProxyServer
}

var grpcServer *grpc.Server

// Stop 停止监听进程
func (s *server) Stop(ctx context.Context, in *pb.ProxyStopRequest) (*pb.ProxyReplyOnly, error) {
    go func() {
        //StopListener()
        tcpmax.Conn.ProxyCBDeleteAll(func(cb *tcpmax.ProxyCB) {
            CloseConn(&cb.Control)
            CloseConn(&cb.Controlled)
            tcpmax.Conn.UpdateStatus(cb.ID, "断开")
        })
        grpcServer.GracefulStop()
    }()
    return &pb.ProxyReplyOnly{Code: "200", Message: "请求已接收,进程退出"}, nil
}

// Disconnect 断开子进程指定连接
func (s *server) Disconnect(ctx context.Context, in *pb.ProxyDisconnectRequest) (*pb.ProxyReplyOnly, error) {
    // 根据四元组查询哈希表，若找到则断开连接
    key := tcpmax.ProxyCBKey{
        SrcIP:   in.GetInConnIP(),
        SrcPort: in.GetInConnPort(),
        DstIP:   in.GetOutConnIP(),
        DstPort: in.GetOutConnPort(),
    }
    log.Println("Disconnect:", key)
    content := tcpmax.Conn.ProxyCBFind(key)
    if content != nil {
        CloseConn(&content.Control)
        CloseConn(&content.Controlled)
        tcpmax.Conn.UpdateStatus(content.ID, "断开")
        tcpmax.Conn.ProxyCBDelete(key)
        return &pb.ProxyReplyOnly{Code: "200", Message: "连接断开成功."}, nil
    }
    tcpmax.Conn.ProxyCBAll()

    return &pb.ProxyReplyOnly{Code: "404", Message: "连接没有找到."}, nil
}

func (s *server) Healthz(ctx context.Context, in *pb.ProxyHealthzRequest) (*pb.ProxyReplyOnly, error) {
    return &pb.ProxyReplyOnly{Code: "200", Message: "The proxy service is working."}, nil
}

// ProxyServer 代理流量转发进程
func ProxyServer(listenIP string, listenPort int, grpcport int, guardIP string, guardPort int) {

    go Trylisten(listenIP, listenPort, guardIP, guardPort)

    // TODO grpc监听端口需要一个规则确定
    lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", grpcport))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    grpcServer = grpc.NewServer()
    pb.RegisterProxyServer(grpcServer, &server{})
    grpcServer.Serve(lis)
}
