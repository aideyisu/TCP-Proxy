package tcpmax

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "net/http"
    pb "src/pb"
    "time"

    "github.com/gin-gonic/gin"
    grpc "google.golang.org/grpc"

    // just import
    _ "github.com/mattn/go-sqlite3"
)

// 回连映射
type ProxyMap struct {
    ID         int64  `json:"id"`
    ListenPort uint32 `json:"listenport"`
    IP         string `json:"ip"`
    Port       uint32 `json:"port"`
    Creator    string `json:"creator"`
    CreateTime uint32 `json:"timestamp"`
}

type proxyMapResponse struct {
    Page     int        `json:"page"`
    Pagesize int        `json:"pagesize"`
    Total    int        `json:"total"`
    Data     []ProxyMap `json:"data"`
}

// 回连映射管理
type ProxyMaps struct {
    Db           *sql.DB
    ProxyMapInfo map[string](*ProxyMap)
}

// Init ProxyMap 创建回连映射
func InitProxyMapTable(dbfile string) *ProxyMaps {
    proxymaps = &ProxyMaps{}

    db, err := sql.Open("sqlite3", dbfile)
    if err != nil {
        log.Fatal("Openning dbfile", err)
    }
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS proxymap(
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            listenport INTEGER NOT NULL UNIQUE,
            ip CHAR(50) NOT NULL,
            port INTEGER  NOT NULL,
            creator CHAR(50) NOT NULL,
            create_time INTEGER NOT NULL,
            UNIQUE(ip,port));
        `)
    if err != nil {
        log.Fatal("Creating table:", err)
    }

    proxymaps.Db = db
    proxymaps.ProxyMapInfo = map[string](*ProxyMap){}
    return proxymaps
}

// Add 添加回连配置映射
func (pm *ProxyMaps) Add(listenport int, ip string, port int, creator string) (*ProxyMap, error) {
    query := fmt.Sprintf(`insert into proxymap(listenport, ip, port, creator, create_time) values ('%d', '%s', '%d', '%s', %d);`,
        listenport, ip, port, creator, int32(time.Now().Unix()))
    _, err := pm.Db.Exec(query)
    if err != nil {
        //log.Print(query, err)
        log.Print(err)
        return nil, err
    }

    return pm.Get(listenport, ip, port), nil
}

// Del 回连映射
func (pm *ProxyMaps) Del(listenport int, ip string, port int) error {
    query := fmt.Sprintf("delete from proxymap where listenport='%d' and ip='%s' and port='%d';", listenport, ip, port)
    _, err := pm.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
        return err
    }
    return nil
}

// Modify 回连映射
func (pm *ProxyMaps) Modify(id int64, listenport int, ip string, port int) (*ProxyMap, error) {
    query := fmt.Sprintf("update proxymap set listenport='%d', ip='%s', port='%d' where id='%d';", listenport, ip, port, id)
    _, err := pm.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
        return nil, err
    }
    return pm.Get(listenport, ip, port), nil
}

// Get 获取回连映射配置
func (pm *ProxyMaps) Get(listenport int, ip string, port int) *ProxyMap {
    query := fmt.Sprintf("select * from proxymap where listenport='%d' and ip='%s' and port='%d';", listenport, ip, port)
    rows, err := pm.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return nil
    }
    defer rows.Close()

    proxymap := ProxyMap{}
    if rows.Next() {
        err = rows.Scan(&proxymap.ID, &proxymap.ListenPort, &proxymap.IP, &proxymap.Port, &proxymap.Creator, &proxymap.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        return &proxymap
    }

    return nil
}

// GetID 获取回连映射配置
func (pm *ProxyMaps) GetID(id int64) *ProxyMap {
    query := fmt.Sprintf("select * from proxymap where id='%d';", id)
    rows, err := pm.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return nil
    }
    defer rows.Close()

    proxymap := ProxyMap{}
    if rows.Next() {
        err = rows.Scan(&proxymap.ID, &proxymap.ListenPort, &proxymap.IP, &proxymap.Port, &proxymap.Creator, &proxymap.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        return &proxymap
    }

    return nil
}

func fmtQueryCondition(listenport int, ip string, port int, creator string, sdate int, edate int) string {
    lpQueryCondition := ""
    andstr := ""
    if listenport != 0 {
        lpQueryCondition = fmt.Sprintf(" listenport='%d'", listenport)
        andstr = " AND"
    }
    ipQueryCondition := ""
    if ip != "" {
        ipQueryCondition = fmt.Sprintf("%s ip='%s'", andstr, ip)
        andstr = " AND"
    }
    portQueryCondition := ""
    if port != 0 {
        portQueryCondition = fmt.Sprintf("%s port='%d'", andstr, port)
        andstr = " AND"
    }
    creatorQueryCondition := ""
    if creator != "" {
        creatorQueryCondition = fmt.Sprintf("%s creator=='%s'", andstr, creator)
        andstr = " AND"
    }
    dateQueryCondition := ""
    if sdate != edate {
        dateQueryCondition = fmt.Sprintf("%s create_time BETWEEN '%d' AND '%d'", andstr, sdate, edate)
    }
    conditions := ""
    if listenport != 0 || ip != "" || port != 0 || creator != "" || sdate != edate {
        conditions = fmt.Sprintf(" where%s%s%s%s%s",
            lpQueryCondition,
            ipQueryCondition,
            portQueryCondition,
            creatorQueryCondition,
            dateQueryCondition)
    }
    return conditions
}

// 查询条目
func (pm *ProxyMaps) Find(listenport int, ip string, port int, creator string, sdate int, edate int, page int, pagesize int) []ProxyMap {
    conditions := fmtQueryCondition(listenport, ip, port, creator, sdate, edate)
    query := fmt.Sprintf("select * from proxymap%s limit %d offset %d;",
        conditions, pagesize, (page-1)*pagesize)
    rows, err := pm.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return nil
    }

    defer rows.Close()

    proxymaps := make([]ProxyMap, 0, pagesize)
    proxymap := ProxyMap{}
    for rows.Next() {
        err = rows.Scan(&proxymap.ID, &proxymap.ListenPort, &proxymap.IP, &proxymap.Port, &proxymap.Creator, &proxymap.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        proxymaps = append(proxymaps, proxymap)
    }
    return proxymaps
}

func (pm *ProxyMaps) Count(listenport int, ip string, port int, creator string, sdate int, edate int) int {
    conditions := fmtQueryCondition(listenport, ip, port, creator, sdate, edate)
    query := fmt.Sprintf("select count(*) from proxymap%s;", conditions)
    rows, err := pm.Db.Query(query)
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
*******************            grcp client      *******************
*******************************************************************/
func AddProxyMapGrpc(address string, creator string, listenPort uint32, ip string, Port uint32) (*pb.ServiceCntrolReplyOnly, error) {
    // address := "localhost:50051"
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(GrpcTimeout))
    if err != nil {
        log.Printf("did not connect: %v", err)
        return nil, err
    }
    defer conn.Close()

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
    defer cancel()

    cn := pb.NewServiceControlClient(conn)
    // 测试创建回连映射配置
    rep, err := cn.ProxyMapAdd(ctx, &pb.ProxyMapAddRequest{Creator: creator, ListenPort: listenPort, IP: ip, Port: Port})
    return rep, err
}

// DelProxyMapGrpc 修改回连映射
func DelProxyMapGrpc(address string, ID int64, ListenPort uint32, IP string, Port uint32) (*pb.ServiceCntrolReplyOnly, error) {
    // address := "localhost:50051"
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(GrpcTimeout))
    if err != nil {
        log.Printf("did not connect: %v", err)
        return nil, err
    }
    defer conn.Close()

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
    defer cancel()

    cn := pb.NewServiceControlClient(conn)
    // 测试回连映射删除
    rep, err := cn.ProxyMapDel(ctx, &pb.ProxyMapDelRequest{ListenPort: uint32(ListenPort), IP: IP, Port: uint32(Port)})
    return rep, err
}

// ModifyProxyMapGrpc 修改回连配置
func ModifyProxyMapGrpc(address string, ID int64, ListenPort uint32, IP string, Port uint32) (*pb.ServiceCntrolReplyOnly, error) {
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(GrpcTimeout))
    if err != nil {
        log.Printf("did not connect: %v", err)
        return nil, err
    }
    defer conn.Close()

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
    defer cancel()

    cn := pb.NewServiceControlClient(conn)
    rep, err := cn.ProxyMapModify(ctx, &pb.ProxyMapModifyRequest{ID: ID, ListenPort: ListenPort, IP: IP, Port: Port})
    return rep, err
}

/******************************************************************
*******************            回连映射          *******************
*******************************************************************/
func AddPorxyMapHandler(c *gin.Context) {
    type addProxyMap struct {
        ListenPort int    `form:"listenport" json:"listenport" binding:"required,numeric,gt=0,lt=65536"`
        IP         string `form:"ip" json:"ip" binding:"required,ipv4"`
        Port       int    `form:"port" json:"port" binding:"required,numeric,gt=0,lt=65536"`
        Creator    string `form:"creator" json:"creator" binding:"required"`
    }
    var param addProxyMap
    if err := c.ShouldBind(&param); err != nil {
        _, reason := checkErr(err)
        c.JSON(http.StatusBadRequest, gin.H{"status_code": 400, "reason": reason})
        return // exit on first error
    }
    rep, err := AddProxyMapGrpc(SCGrpcAddr, param.Creator, uint32(param.ListenPort), param.IP, uint32(param.Port))
    if err != nil {
        StatusCode, Reason := checkErr(err)
        c.JSON(int(StatusCode), gin.H{
            "status_code": StatusCode,
            "reason":      Reason,
        })
    } else {
        c.JSON(int(rep.StatusCode), gin.H{
            "status_code": rep.StatusCode,
            "reason":      rep.Reason,
        })
    }
}

func ModifyProxyMapHandler(c *gin.Context) {
    type modifyProxyMap struct {
        ID         int64  `form:"id" json:"id" binding:"required,numeric,gte=0"`
        ListenPort int    `form:"listenport" json:"listenport" binding:"required,numeric,gt=0,lt=65536"`
        IP         string `form:"ip" json:"ip" binding:"required,ipv4"`
        Port       int    `form:"port" json:"port" binding:"required,numeric,gt=0,lt=65536"`
        Creator    string `form:"creator" json:"creator" binding:"required"`
    }
    var param modifyProxyMap
    if err := c.ShouldBind(&param); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status_code": 400, "reason": fmt.Sprint(err.Error())})
        return // exit on first error
    }
    rep, err := ModifyProxyMapGrpc(SCGrpcAddr, param.ID, uint32(param.ListenPort), param.IP, uint32(param.Port))
    if err != nil {
        StatusCode, Reason := checkErr(err)
        c.JSON(int(StatusCode), gin.H{
            "status_code": StatusCode,
            "reason":      Reason,
        })
    } else {
        c.JSON(int(rep.StatusCode), gin.H{
            "status_code": rep.StatusCode,
            "reason":      rep.Reason,
        })
    }
}

func ListProxyMapHandler(c *gin.Context) {
    type Param struct {
        Page       int    `form:"page" json:"page" binding:"required,numeric,gt=0"`
        PageSize   int    `form:"pagesize" json:"pagesize" binding:"required,numeric,gt=0,lte=100"`
        Creator    string `form:"creator" json:"creator" binding:"omitempty"`
        ListenPort int    `form:"listenport" json:"listenport" binding:"omitempty,numeric,gt=0,lt=65536"`
        IP         string `form:"ip" json:"ip" binding:"omitempty,ipv4"`
        Port       int    `form:"port" json:"port" binding:"omitempty,numeric,gt=0,lt=65536"`
        Sdate      int    `form:"sdate" json:"sdate" binding:"omitempty,numeric,ltfield=Edate"`
        Edate      int    `form:"edate" json:"edate" binding:"omitempty,numeric,gtfield=Sdate"`
    }
    var param Param
    if err := c.ShouldBind(&param); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status_code": 400, "reason": fmt.Sprint(err.Error())})
        return // exit on first error
    }
    response := proxyMapResponse{}
    response.Page = param.Page
    response.Pagesize = param.PageSize
    response.Total = proxymaps.Count(param.ListenPort, param.IP, param.Port,
        param.Creator, param.Sdate, param.Edate)
    response.Data = proxymaps.Find(param.ListenPort, param.IP, param.Port,
        param.Creator, param.Sdate, param.Edate,
        param.Page, param.PageSize)
    c.IndentedJSON(http.StatusOK, response)
}

func DelProxyMapHandler(c *gin.Context) {
    type delProxyMap struct {
        ListenPort int    `form:"listenport" json:"listenport" binding:"required,numeric,gt=0,lt=65536"`
        IP         string `form:"ip" json:"ip" binding:"required,ipv4"`
        Port       int    `form:"port" json:"port" binding:"required,numeric,gt=0,lt=65536"`
    }
    var param delProxyMap
    if err := c.ShouldBind(&param); err != nil {
        _, reason := checkErr(err)
        c.JSON(http.StatusBadRequest, gin.H{"status_code": 400, "reason": reason})
        return // exit on first error
    }
    rep, err := DelProxyMapGrpc(SCGrpcAddr, 1, uint32(param.ListenPort), param.IP, uint32(param.Port))
    if err != nil {
        StatusCode, Reason := checkErr(err)
        c.JSON(int(StatusCode), gin.H{
            "status_code": StatusCode,
            "reason":      Reason,
        })
    } else {
        c.JSON(int(rep.StatusCode), gin.H{
            "status_code": rep.StatusCode,
            "reason":      rep.Reason,
        })
    }
}
