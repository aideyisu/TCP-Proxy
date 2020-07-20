package tcpmax

import (
    "context"
    "fmt"
    "io/ioutil"
    "log"
    pb "src/pb"
    "strconv"

    "strings"
    "time"

    "google.golang.org/grpc"
)

// Ioutil 读取文件看配置
func Ioutil(name string) string {
    contents, _ := ioutil.ReadFile(name)
    //因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
    result := strings.Replace(string(contents), "\n", "", 1)
    fmt.Println("Use ioutil.ReadFile to read a file:", result)
    return result
}

func proxyReset(address string, inIP string, inPort uint32, outIP string, outPort uint32) {
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    // 单一断开
    c := pb.NewProxyClient(conn)
    r, err := c.Disconnect(ctx, &pb.ProxyDisconnectRequest{InConnIP: inIP, InConnPort: inPort, OutConnIP: outIP, OutConnPort: outPort})
    if err != nil {
        log.Fatalf("could not greet: %v", err)
    }
    log.Printf("ListenAddr结果: %s", r.GetMessage())
}

func proxyStop(address string) {
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    // 单一断开
    c := pb.NewProxyClient(conn)

    // 测试接口全部断开连接,已经测通，但不要轻易调用..毕竟是整个都退出,威力比较大 ！！！！ 已经将关闭暂时关闭，记得在用到时把注释去掉
    r, err := c.Stop(ctx, &pb.ProxyStopRequest{})
    if err != nil {
        log.Fatalf("could not greet: %v", err)
    }
    log.Printf("ListenAddr结果: %s", r.GetMessage())
}

// proxyClient 之前参数 localhost:50053
func proxyClient(address string) {
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
    defer cancel()

    c := pb.NewProxyClient(conn)
    A := Ioutil("testwritefile.txt")
    fmt.Println(A)
    S := strings.Split(A, " ")
    fmt.Println(len(S))

    listenIP := S[0]
    listenPort, _ := strconv.Atoi(S[1])
    guardIP := S[2]
    guardPort, _ := strconv.Atoi(S[3])
    fmt.Println(listenIP, listenPort, guardIP, guardPort) // 直接读取文件数据断开连接

    // 测试接口全部断开连接,已经测通，但不要轻易调用..毕竟是整个都退出,威力比较大 ！！！！ 已经将关闭暂时关闭，记得在用到时把注释去掉
    r, err := c.Stop(ctx, &pb.ProxyStopRequest{})
    if err != nil {
        log.Fatalf("could not greet: %v", err)
    }
    log.Printf("ListenAddr结果: %s", r.GetMessage())

    // 测试接口单一断开
    r, err = c.Disconnect(ctx, &pb.ProxyDisconnectRequest{InConnIP: listenIP, InConnPort: uint32(listenPort), OutConnIP: guardIP, OutConnPort: uint32(guardPort)})
    if err != nil {
        log.Fatalf("could not greet: %v", err)
    }
    log.Printf("ListenAddr结果: %s", r.GetMessage())
}
