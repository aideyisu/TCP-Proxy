package tcpmax

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "net"
    "net/http"
    pb "src/pb"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    grpc "google.golang.org/grpc"
    "gopkg.in/go-playground/validator.v9"

    // just import
    _ "github.com/mattn/go-sqlite3"
)

type ProxyCBKey struct {
    SrcIP   string
    SrcPort uint32
    DstIP   string
    DstPort uint32
}

// ProxyCB 代理控制模块
type ProxyCB struct {
    Key        ProxyCBKey
    Controlled net.Conn //受控端socket
    Control    net.Conn //控制端socket
    Send       uint64   // 发送的字节数
    Received   uint64   // 接收端字节数
    ID         int64
}

// 连接信息
type Connection struct {
    ID         int64  `json:"id"`
    Protocol   string `json:"protocol"`
    SrcIP      string `json:"src_ip"`
    DstIP      string `json:"dst_ip"`
    SrcPort    uint32 `json:"src_port"`
    DstPort    uint32 `json:"dst_port"`
    Received   uint64 `json:"received"`
    Send       uint64 `json:"send"`
    Status     string `json:"status"`
    CreateTime uint32 `json:"timestamp"`
}

type ConnectionResponse struct {
    Page     int          `json:"page"`
    Pagesize int          `json:"pagesize"`
    Total    int          `json:"total"`
    Data     []Connection `json:"data"`
}

// 白名单管理
type Connections struct {
    Db *sql.DB
    //ConnectionInfo map[string](*Connection)
    proxyCBMap map[ProxyCBKey](*ProxyCB)
}

// ProxyCBAdd 新增key
func (conn *Connections) ProxyCBAdd(key ProxyCBKey, value *ProxyCB) {
    log.Println("PCBADD", key, value)
    conn.proxyCBMap[key] = value
}

// ProxyCBFind 查找key
func (conn *Connections) ProxyCBFind(key ProxyCBKey) *ProxyCB {
    Value := conn.proxyCBMap[key]
    log.Println("PCBFIND", key, Value)
    return Value
}

// ProxyCBDelete 删除key
func (conn *Connections) ProxyCBDelete(key ProxyCBKey) {
    log.Println("PCBDELETE", key, conn.proxyCBMap[key])
    delete(conn.proxyCBMap, key)
}

func (conn *Connections) ProxyCBDeleteAll(fn func(cb *ProxyCB)) {
    for key := range conn.proxyCBMap {
        fn(conn.proxyCBMap[key])
        delete(conn.proxyCBMap, key)
    }
}

func (conn *Connections) ProxyCBAll() {
    for key := range conn.proxyCBMap {
        log.Println("PBCALL:", key, conn.proxyCBMap[key])
    }
}

// InitConnectionTable 创建回连映射
func InitConnectionTable(dbfile string) *Connections {
    connections = &Connections{}

    db, err := sql.Open("sqlite3", dbfile)
    if err != nil {
        log.Fatal("Openning dbfile", err)
    }
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS connection(
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            protocol CHAR(50) NOT NULL,
            src_ip CHAR(50) NOT NULL,
            dst_ip CHAR(50) NOT NULL,
            src_port INTEGER NOT NULL,
            dst_port INTEGER NOT NULL,
            received INTEGER NOT NULL,
            send INTEGER NOT NULL,
            status CHAR(50) NOT NULL,
            create_time INTEGER NOT NULL);
        `)
    if err != nil {
        log.Fatal("Creating table:", err)
    }

    connections.Db = db
    //connections.ConnectionInfo = map[string](*Connection){}
    connections.proxyCBMap = make(map[ProxyCBKey](*ProxyCB))
    return connections
}

// Add 添加回连配置映射
func (conn *Connections) Add(protocol string, src_ip string, dst_ip string, src_port int, dst_port int, received int, send int, status string) *Connection {
    query := fmt.Sprintf(`insert into connection(protocol, src_ip, dst_ip, src_port, dst_port, received, send, status, create_time) values ('%s', '%s', '%s', '%d', '%d', '%d', '%d', '%s', %d);`,
        protocol, src_ip, dst_ip, src_port, dst_port, received, send, status, int64(time.Now().Unix()))
    result, err := conn.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
        return nil
    }
    id, err := result.LastInsertId()

    return conn.Get(id)
}

// UpdateStatus 更新状态启用停用
func (conn *Connections) UpdateStatus(id int64, status string) *Connection {
    query := fmt.Sprintf("update connection set status='%s' where id='%d';", status, id)
    _, err := conn.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
        return nil
    }
    return conn.Get(id)
}

// UpdateReceived 更新流量统计-接收
func (conn *Connections) UpdateReceived(id int64, received int64) *Connection {
    query := fmt.Sprintf("update connection set received=%d where id='%d';", received, id)
    _, err := conn.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
        return nil
    }
    return conn.Get(id)
}

// UpdateSend 更新流量统计-发送
func (conn *Connections) UpdateSend(id int64, send int64) *Connection {
    query := fmt.Sprintf("update connection set send=%d where id='%d';", send, id)
    _, err := conn.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
        return nil
    }
    return conn.Get(id)
}

// Del 回连映射
func (conn *Connections) Del(id int64) {
    query := fmt.Sprintf("delete from connection where id='%d';", id)
    _, err := conn.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
    }
}

// Get 获取回连映射配置
func (conn *Connections) Get(id int64) *Connection {
    query := fmt.Sprintf("select * from connection where id='%d';", id)
    rows, err := conn.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return nil
    }
    defer rows.Close()

    connection := Connection{}
    if rows.Next() {
        err = rows.Scan(&connection.ID,
            &connection.Protocol,
            &connection.SrcIP,
            &connection.DstIP,
            &connection.SrcPort,
            &connection.DstPort,
            &connection.Received,
            &connection.Send,
            &connection.Status,
            &connection.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        return &connection
    }

    return nil
}

func connFmtQueryCondition(protocol string, status string, srcIp string, dstIp string, srcPort int, dstPort int, sdate int, edate int) string {
    protocolQueryCondition := ""
    andstr := ""
    if protocol != "" {
        protocolQueryCondition = fmt.Sprintf(" protocol=='%s'", protocol)
        andstr = " AND"
    }
    statusQueryCondition := ""
    if status != "" {
        statusQueryCondition = fmt.Sprintf("%s status=='%s'", andstr, status)
        andstr = " AND"
    }
    srcIpQueryCondition := ""
    if srcIp != "" {
        srcIpQueryCondition = fmt.Sprintf("%s src_ip=='%s'", andstr, srcIp)
        andstr = " AND"
    }
    dstIpQueryCondition := ""
    if dstIp != "" {
        dstIpQueryCondition = fmt.Sprintf("%s dst_ip=='%s'", andstr, dstIp)
        andstr = " AND"
    }
    srcPortQueryCondition := ""
    if srcPort != 0 {
        srcPortQueryCondition = fmt.Sprintf("%s src_port=='%d'", andstr, srcPort)
        andstr = " AND"
    }
    dstPortQueryCondition := ""
    if dstPort != 0 {
        dstPortQueryCondition = fmt.Sprintf("%s dst_port=='%d'", andstr, dstPort)
        andstr = " AND"
    }
    dateQueryCondition := ""
    if sdate != edate {
        dateQueryCondition = fmt.Sprintf("%s create_time BETWEEN '%d' AND '%d'", andstr, sdate, edate)
    }
    conditions := ""
    if protocol != "" || status != "" || srcIp != "" || dstIp != "" || srcPort != 0 || dstPort != 0 || sdate != edate {
        conditions = fmt.Sprintf(" where%s%s%s%s%s%s%s",
            protocolQueryCondition,
            statusQueryCondition,
            srcIpQueryCondition,
            dstIpQueryCondition,
            srcPortQueryCondition,
            dstPortQueryCondition,
            dateQueryCondition)
    }
    return conditions
}

// Find 查询符合条件的连接日志条目
func (conn *Connections) Find(protocol string, status string, srcIp string, dstIp string, srcPort int, dstPort int, sdate int, edate int, page int, pagesize int) []Connection {
    conditions := connFmtQueryCondition(protocol, status, srcIp, dstIp, srcPort, dstPort, sdate, edate)
    query := fmt.Sprintf("select * from connection%s limit %d offset %d;",
        conditions, pagesize, (page-1)*pagesize)
    log.Printf(query)
    rows, err := conn.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return nil
    }

    defer rows.Close()

    connections := make([]Connection, 0, pagesize)
    connection := Connection{}
    for rows.Next() {
        err = rows.Scan(&connection.ID,
            &connection.Protocol,
            &connection.SrcIP,
            &connection.DstIP,
            &connection.SrcPort,
            &connection.DstPort,
            &connection.Received,
            &connection.Send,
            &connection.Status,
            &connection.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        connections = append(connections, connection)
    }
    return connections
}

// Count 获取符合条件的连接数量
func (conn *Connections) Count(protocol string, status string, srcIp string, dstIp string, srcPort int, dstPort int, sdate int, edate int) int {
    conditions := connFmtQueryCondition(protocol, status, srcIp, dstIp, srcPort, dstPort, sdate, edate)
    query := fmt.Sprintf("select count(*) from connection%s;", conditions)
    rows, err := conn.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return 0
    }

    defer rows.Close()

    var count int
    for rows.Next() {
        err = rows.Scan(&count)
        if err != nil {
            log.Printf(err.Error())
            return 0
        }
    }
    return count
}

/******************************************************************
*******************       在线连接管理 grpc       *******************
*******************************************************************/
func DisConnectionGrpc(protocol string, srcIP string, dstIP string, srcPort int, dstPort int) (*pb.ProxyReplyOnly, error) {
    pgpms := PGPMS.Get(dstIP, dstPort)
    if pgpms == nil {
        return &pb.ProxyReplyOnly{Code: "404", Message: "The proxy service is not found."}, nil
    }
    addr := fmt.Sprintf("127.0.0.1:%d", ProxyGrpcPortStart+pgpms.OffSet)
    conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(GrpcTimeout))
    if err != nil {
        log.Printf("did not connect: %v", err)
        return nil, err
    }
    defer conn.Close()

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
    defer cancel()

    c := pb.NewProxyClient(conn)
    rep, err := c.Disconnect(ctx,
        &pb.ProxyDisconnectRequest{Protocol: protocol, InConnIP: srcIP, OutConnIP: dstIP, InConnPort: uint32(srcPort), OutConnPort: uint32(dstPort)})
    return rep, err
}

/******************************************************************
*******************          在线连接管理         *******************
*******************************************************************/
func ListOnlineConnectionsHandler(c *gin.Context) {
    type Param struct {
        Page     int    `form:"page" json:"page" binding:"required,numeric,gt=0"`
        PageSize int    `form:"pagesize" json:"pagesize" binding:"required,numeric,gt=0,lte=100"`
        Protocol string `form:"protocol" json:"protocol" binding:"omitempty,oneof=TCP"`
        SrcIP    string `form:"src_ip" json:"src_ip" binding:"omitempty,ipv4"`
        DstIP    string `form:"dst_ip" json:"dst_ip" binding:"omitempty,ipv4"`
        SrcPort  int    `form:"src_port" json:"src_port" binding:"omitempty,numeric,gt=0,lt=65536"`
        DstPort  int    `form:"dst_port" json:"dst_port" binding:"omitempty,numeric,gt=0,lt=65536"`
        Sdate    int    `form:"sdate" json:"sdate" binding:"omitempty,numeric,ltfield=Edate"`
        Edate    int    `form:"edate" json:"edate" binding:"omitempty,numeric,gtfield=Sdate"`
    }
    var param Param
    if err := c.ShouldBindQuery(&param); err != nil {
        for _, fieldErr := range err.(validator.ValidationErrors) {
            c.JSON(http.StatusBadRequest, gin.H{"status_code": 400, "reason": fmt.Sprint(fieldErr)})
            return // exit on first error
        }
        return
    }
    response := ConnectionResponse{}
    response.Page = param.Page
    response.Pagesize = param.PageSize
    response.Total = connections.Count(param.Protocol, "连接",
        param.SrcIP, param.DstIP,
        param.SrcPort, param.DstPort,
        param.Sdate, param.Edate)
    response.Data = connections.Find(param.Protocol, "连接",
        param.SrcIP, param.DstIP,
        param.SrcPort, param.DstPort,
        param.Sdate, param.Edate,
        param.Page, param.PageSize)

    c.IndentedJSON(http.StatusOK, response)
}

// 断开连接
func DisConnectionHandler(c *gin.Context) {
    type Param struct {
        Protocol string `form:"protocol" json:"protocol" binding:"required,oneof=TCP"`
        SrcIP    string `form:"src_ip" json:"src_ip" binding:"required,ipv4"`
        DstIP    string `form:"dst_ip" json:"dst_ip" binding:"required,ipv4"`
        SrcPort  int    `form:"src_port" json:"src_port" binding:"required,numeric,gt=0,lt=65536"`
        DstPort  int    `form:"dst_port" json:"dst_port" binding:"required,numeric,gt=0,lt=65536"`
    }
    var param Param
    if err := c.ShouldBind(&param); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status_code": 400, "reason": fmt.Sprint(err.Error())})
        return // exit on first error
    }

    rep, err := DisConnectionGrpc(param.Protocol, param.SrcIP, param.DstIP, param.SrcPort, param.DstPort)
    if err != nil {
        StatusCode, Reason := checkErr(err)
        c.JSON(int(StatusCode), gin.H{
            "status_code": StatusCode,
            "reason":      Reason,
        })
    } else {
        code, _ := strconv.Atoi(rep.Code)
        c.JSON(code, gin.H{
            "status_code": code,
            "reason":      rep.Message,
        })
    }
}

/******************************************************************
*******************          连接日志            *******************
*******************************************************************/
func ListConnectionsLogsHandler(c *gin.Context) {
    type Param struct {
        Page     int    `form:"page" json:"page" binding:"required,numeric,gt=0"`
        PageSize int    `form:"pagesize" json:"pagesize" binding:"required,numeric,gt=0,lte=100"`
        Protocol string `form:"protocol" json:"protocol" binding:"omitempty,oneof=TCP"`
        SrcIP    string `form:"src_ip" json:"src_ip" binding:"omitempty,ipv4"`
        DstIP    string `form:"dst_ip" json:"dst_ip" binding:"omitempty,ipv4"`
        SrcPort  int    `form:"src_port" json:"src_port" binding:"omitempty,numeric,gt=0,lt=65536"`
        DstPort  int    `form:"dst_port" json:"dst_port" binding:"omitempty,numeric,gt=0,lt=65536"`
        Sdate    int    `form:"sdate" json:"sdate" binding:"omitempty,numeric,ltfield=Edate"`
        Edate    int    `form:"edate" json:"edate" binding:"omitempty,numeric,gtfield=Sdate"`
    }
    var param Param
    if err := c.ShouldBind(&param); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status_code": 400, "reason": fmt.Sprint(err.Error())})
        return // exit on first error
    }

    response := ConnectionResponse{}
    response.Page = param.Page
    response.Pagesize = param.PageSize
    response.Total = connections.Count(param.Protocol, "断开",
        param.SrcIP, param.DstIP,
        param.SrcPort, param.DstPort,
        param.Sdate, param.Edate)
    response.Data = connections.Find(param.Protocol, "断开",
        param.SrcIP, param.DstIP,
        param.SrcPort, param.DstPort,
        param.Sdate, param.Edate,
        param.Page, param.PageSize)

    c.IndentedJSON(http.StatusOK, response)
}
