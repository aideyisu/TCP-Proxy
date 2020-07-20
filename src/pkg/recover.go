package tcpmax

import (
	"context"
	"fmt"
	"log"
	"time"
	pb "src/pb"
	"google.golang.org/grpc"
)

func checkHealthz(listenip string, listenport int) (*pb.ProxyReplyOnly, error) {
    pgpm := PGPMS.Get(listenip, listenport)
    if pgpm == nil {
        return &pb.ProxyReplyOnly{Code: "404", Message: "Proxy grpc port map not found, can't stop proxy server."}, nil
    }
    grpcport := ProxyGrpcPortStart + pgpm.OffSet
    addr := fmt.Sprintf("127.0.0.1:%d", grpcport)
    conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(GrpcTimeout))
    if err != nil {
        log.Printf("did not connect: %v", err)
        return nil, err
    }
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
    defer cancel()
    cn := pb.NewProxyClient(conn)
    rep, err := cn.Healthz(ctx, &pb.ProxyHealthzRequest{})
    return rep, err
}

func recoverProxyProcessByIP(configFile string, ip string) {
	proxyMapCount := proxymaps.Count(0, "", 0, "", 0, 0)
    for i := 0; i < proxyMapCount; i++ {
        proxymap := proxymaps.Find(0, "", 0, "", 0, 0, i+1, 1)
		pgpm := PGPMS.Get(ip, int(proxymap[0].ListenPort))
		if pgpm != nil {
			// 服务在线检查
			_, err := checkHealthz(ip, int(proxymap[0].ListenPort))
			if err != nil {
				PGPMS.Del(ip, int(proxymap[0].ListenPort))
				PGPMS.GrpcPortIDMap.Clear(pgpm.OffSet)
				log.Printf("Proxy server(%s,%d) restart.", ip, proxymap[0].ListenPort)
				startProxyServer(configFile, ip, int(proxymap[0].ListenPort), proxymap[0].IP, int(proxymap[0].Port))
			} else {
				log.Printf("Proxy server(%s,%d) still running.", ip, proxymap[0].ListenPort)
			}
		} else {
			// 启动代理服务
			startProxyServer(configFile, ip, int(proxymap[0].ListenPort), proxymap[0].IP, int(proxymap[0].Port))
			log.Printf("Proxy server(%s,%d) start.", ip, proxymap[0].ListenPort)
		}
    }
}

// 进程启动时恢复代理服务进程
func RecoverProxyProcess(configFile string) {
	PGPMS.Loads()
	listenAddrCount := listenaddrs.Count("", "", 0, 0)
    for i := 0; i < listenAddrCount; i++ {
		las := listenaddrs.Find("", "", 0, 0, i+1, 1)
		recoverProxyProcessByIP(configFile, las[0].IP)
	}
}
