package tcpmax

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "time"

    "net/http"
    pb "src/pb"

    "github.com/gin-gonic/gin"
    grpc "google.golang.org/grpc"

    // just import
    _ "github.com/mattn/go-sqlite3"
)

// User 用户信息
type ListenAddr struct {
    ID         int64  `json:"id"`
    IP         string `json:"ip"`
    Creator    string `json:"creator"`
    CreateTime uint32 `json:"timestamp"`
}

type listenAddrResponse struct {
    Page     int          `json:"page"`
    Pagesize int          `json:"pagesize"`
    Total    int          `json:"total"`
    Data     []ListenAddr `json:"data"`
}

// Users 用户管理
type ListenAddrs struct {
    Db             *sql.DB
    ListenAddrInfo map[string](*ListenAddr)
}

// InitUsers 创建用户管理
func InitListenAddrTable(dbfile string) *ListenAddrs {
    listenaddrs = &ListenAddrs{}

    db, err := sql.Open("sqlite3", dbfile)
    if err != nil {
        log.Fatal("Openning dbfile", err)
    }
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS listenip (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        ip VARCHAR(50) NOT NULL UNIQUE,
        creator VARCHAR(64) NOT NULL,
        create_time INTEGER NOT NULL);
        `)
    if err != nil {
        log.Fatal("Creating table:", err)
    }

    listenaddrs.Db = db
    listenaddrs.ListenAddrInfo = map[string](*ListenAddr){}
    return listenaddrs
}

// Add 添加监听地址
func (la *ListenAddrs) Add(ip string, creator string) (*ListenAddr, error) {
    query := fmt.Sprintf(`insert into listenip(ip, creator, create_time) values ('%s', '%s', %d);`,
        ip, creator, int32(time.Now().Unix()))
    _, err := la.Db.Exec(query)
    if err != nil {
        //log.Print(query, err)
        log.Print(err)
        return nil, err
    }

    return la.Get(ip), nil
}

// Del 删除监听地址
func (la *ListenAddrs) Del(ip string) error {
    query := fmt.Sprintf("delete from listenip where ip='%s';", ip)
    _, err := la.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
        return err
    }
    return nil
}

// Modify 修改监听地址
func (la *ListenAddrs) Modify(id int64, ip string) (*ListenAddr, error) {
    query := fmt.Sprintf("update listenip set ip='%s' where id='%d';", ip, id)
    _, err := la.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
        return nil, err
    }
    return la.Get(ip), nil
}

// Get 根据ip地址获取监听地址
func (la *ListenAddrs) Get(ip string) *ListenAddr {
    query := fmt.Sprintf("select * from listenip where ip='%s';", ip)
    rows, err := la.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return nil
    }
    defer rows.Close()

    listenaddr := ListenAddr{}
    if rows.Next() {
        err = rows.Scan(&listenaddr.ID, &listenaddr.IP, &listenaddr.Creator, &listenaddr.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        return &listenaddr
    }

    return nil
}

// GetID 根据id查询监听地址信息
func (la *ListenAddrs) GetID(id int64) *ListenAddr {
    query := fmt.Sprintf("select * from listenip where id='%d';", id)
    rows, err := la.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return nil
    }
    defer rows.Close()

    listenaddr := ListenAddr{}
    if rows.Next() {
        err = rows.Scan(&listenaddr.ID, &listenaddr.IP, &listenaddr.Creator, &listenaddr.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        return &listenaddr
    }

    return nil
}

// 查询条目
func (la *ListenAddrs) Find(ip string, creator string, sdate int, edate int, page int, pagesize int) []ListenAddr {
    ipQueryCondition := ""
    andstr := ""
    if ip != "" {
        ipQueryCondition = fmt.Sprintf(" ip='%s'", ip)
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
    if ip != "" || creator != "" || sdate != edate {
        conditions = fmt.Sprintf(" where%s%s%s", ipQueryCondition, creatorQueryCondition, dateQueryCondition)
    }
    query := fmt.Sprintf("select * from listenip%s limit %d offset %d;",
        conditions,
        pagesize, (page-1)*pagesize)
    rows, err := la.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return nil
    }

    defer rows.Close()

    listenaddrs := make([]ListenAddr, 0, pagesize)
    listenaddr := ListenAddr{}
    for rows.Next() {
        err = rows.Scan(&listenaddr.ID, &listenaddr.IP, &listenaddr.Creator, &listenaddr.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        listenaddrs = append(listenaddrs, listenaddr)
    }

    return listenaddrs
}

func (la *ListenAddrs) Count(ip string, creator string, sdate int, edate int) int {
    ipQueryCondition := ""
    andstr := ""
    if ip != "" {
        ipQueryCondition = fmt.Sprintf(" ip='%s'", ip)
        andstr = " AND"
    }
    creatorQueryCondition := ""
    if creator != "" {
        creatorQueryCondition = fmt.Sprintf("%s creator='%s'", andstr, creator)
        andstr = " AND"
    }
    dateQueryCondition := ""
    if sdate != edate {
        dateQueryCondition = fmt.Sprintf("%s create_time BETWEEN '%d' AND '%d'", andstr, sdate, edate)
    }
    conditions := ""
    if ip != "" || creator != "" || sdate != edate {
        conditions = fmt.Sprintf(" where%s%s%s", ipQueryCondition, creatorQueryCondition, dateQueryCondition)
    }
    query := fmt.Sprintf("select count(*) from listenip%s;", conditions)
    rows, err := la.Db.Query(query)
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
*******************            grpc client      *******************
*******************************************************************/
// 新增监听IP
func AddrListenAddrGrpc(address string, IP string, Createor string) (*pb.ServiceCntrolReplyOnly, error) {
    //conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(GrpcTimeout))
    if err != nil {
        log.Printf("did not connect: %v", err)
        return nil, err
    }
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
    defer cancel()
    cn := pb.NewServiceControlClient(conn)
    rep, err := cn.ListenAddrAdd(ctx, &pb.ListenAddrAddRequest{IP: IP, Createor: Createor})
    return rep, err
}

// 删除监听IP
func DelListenAddrGrpc(address string, ip string) (*pb.ServiceCntrolReplyOnly, error) {
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(GrpcTimeout))
    if err != nil {
        log.Printf("did not connect: %v", err)
        return nil, err
    }
    defer conn.Close()

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
    defer cancel()

    c := pb.NewServiceControlClient(conn)
    rep, err := c.ListenAddrDel(ctx, &pb.ListenAddrDelRequest{IP: ip})
    return rep, err
}

// 修改监听IP
func ModifyListenAddrGrpc(address string, ip string, id int64) (*pb.ServiceCntrolReplyOnly, error) {
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(GrpcTimeout))
    if err != nil {
        log.Printf("did not connect: %v", err)
        return nil, err
    }
    defer conn.Close()

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
    defer cancel()

    c := pb.NewServiceControlClient(conn)
    // 测试监听接口修改
    rep, err := c.ListenAddrModify(ctx, &pb.ListenAddrModifyRequest{IP: ip, ID: id})
    return rep, err
}

/******************************************************************
*******************            监听地址          *******************
*******************************************************************/

// 添加监听地址配置
func AddListenaddrHandler(c *gin.Context) {
    type Param struct {
        Creator string `form:"creator" json:"creator" binding:"required"`
        IP      string `form:"ip" json:"ip" binding:"required,ipv4"`
    }
    var param Param
    if err := c.ShouldBind(&param); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status_code": 400, "reason": fmt.Sprint(err.Error())})
        return // exit on first error
    }
    rep, err := AddrListenAddrGrpc(SCGrpcAddr, param.IP, param.Creator)
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

// 修改配置
func ModifyListenaddrHandler(c *gin.Context) {
    type Param struct {
        ID int64  `form:"id" json:"id" binding:"required,numeric,gte=0"`
        IP string `form:"ip" json:"ip" binding:"required,ipv4"`
    }
    var param Param
    if err := c.ShouldBind(&param); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status_code": 400, "reason": fmt.Sprint(err.Error())})
        return // exit on first error
    }
    rep, err := ModifyListenAddrGrpc(SCGrpcAddr, param.IP, param.ID)
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

// 获取监听配置列表
func ListListenaddrHandler(c *gin.Context) {
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

    response := listenAddrResponse{}
    response.Page = param.Page
    response.Pagesize = param.PageSize
    response.Total = listenaddrs.Count(param.IP, param.Creator, param.Sdate, param.Edate)
    response.Data = listenaddrs.Find(param.IP, param.Creator, param.Sdate, param.Edate, param.Page, param.PageSize)
    c.IndentedJSON(http.StatusOK, response)
}

// 删除监听配置
func DelListenaddrHandler(c *gin.Context) {
    type Param struct {
        IP string `form:"ip" json:"ip" binding:"required,ipv4"`
    }
    var param Param
    if err := c.ShouldBind(&param); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status_code": 400, "reason": fmt.Sprint(err.Error())})
        return // exit on first error
    }
    rep, err := DelListenAddrGrpc(SCGrpcAddr, param.IP)
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
