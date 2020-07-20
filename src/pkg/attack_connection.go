package tcpmax

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"

    // just import
    _ "github.com/mattn/go-sqlite3"
)

// 连接信息
type AttackConnection struct {
    ID         int64  `json:"id"`
    Protocol   string `json:"protocol"`
    SrcIP      string `json:"src_ip"`
    DstIP      string `json:"dst_ip"`
    SrcPort    uint32 `json:"src_port"`
    DstPort    uint32 `json:"dst_port"`
    CreateTime uint32 `json:"timestamp"`
}

type AttackConnectionResponse struct {
    Page     int                `json:"page"`
    Pagesize int                `json:"pagesize"`
    Total    int                `json:"total"`
    Data     []AttackConnection `json:"data"`
}

// 非授权名单管理
type AttackConnections struct {
    Db                   *sql.DB
    AttackConnectionInfo map[string](*AttackConnection)
}

// InitAttackConnectionTable 创建非授权名单管理
func InitAttackConnectionTable(dbfile string) *AttackConnections {
    atkconnections := AttackConnections{}

    db, err := sql.Open("sqlite3", dbfile)
    if err != nil {
        log.Fatal("Openning dbfile", err)
    }
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS attack_connection(
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            protocol CHAR(50) NOT NULL,
            src_ip CHAR(50) NOT NULL,
            dst_ip CHAR(50) NOT NULL,
            src_port INTEGER NOT NULL,
            dst_port INTEGER NOT NULL,
            create_time INTEGER NOT NULL);
        `)
    if err != nil {
        log.Fatal("Creating table:", err)
    }

    atkconnections.Db = db
    atkconnections.AttackConnectionInfo = map[string](*AttackConnection){}
    attackconnections = &atkconnections
    return &atkconnections
}

// Add 添加非授权日志
func (conn *AttackConnections) Add(protocol string, src_ip string, dst_ip string, src_port int, dst_port int) *AttackConnection {
    query := fmt.Sprintf(`insert into attack_connection(protocol, src_ip, dst_ip, src_port, dst_port, create_time) values ('%s', '%s', '%s', '%d', '%d', '%d');`,
        protocol, src_ip, dst_ip, src_port, dst_port, int64(time.Now().Unix()))
    result, err := conn.Db.Exec(query)
    if err != nil {
        //log.Print(query, err)
        log.Print(err)
        return nil
    }
    id, err := result.LastInsertId()

    return conn.Get(id)
}

// Del 删除非授权日志
func (conn *AttackConnections) Del(id int64) {
    query := fmt.Sprintf("delete from attack_connection where id='%d';", id)
    _, err := conn.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
    }
}

// Get 通过ID查询非授权日志
func (conn *AttackConnections) Get(id int64) *AttackConnection {
    query := fmt.Sprintf("select * from attack_connection where id='%d';", id)
    rows, err := conn.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return nil
    }
    defer rows.Close()

    atkconnection := AttackConnection{}
    if rows.Next() {
        err = rows.Scan(&atkconnection.ID,
            &atkconnection.Protocol,
            &atkconnection.SrcIP,
            &atkconnection.DstIP,
            &atkconnection.SrcPort,
            &atkconnection.DstPort,
            &atkconnection.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        return &atkconnection
    }

    return nil
}

func atkconnFmtQueryCondition(protocol string, srcIp string, dstIp string, srcPort int, dstPort int, sdate int, edate int) string {
    protocolQueryCondition := ""
    andstr := ""
    if protocol != "" {
        protocolQueryCondition = fmt.Sprintf(" protocol=='%s'", protocol)
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
    if protocol != "" || srcIp != "" || dstIp != "" || srcPort != 0 || dstPort != 0 || sdate != edate {
        conditions = fmt.Sprintf(" where%s%s%s%s%s%s",
            protocolQueryCondition,
            srcIpQueryCondition,
            dstIpQueryCondition,
            srcPortQueryCondition,
            dstPortQueryCondition,
            dateQueryCondition)
    }
    return conditions
}

// 查询符合条件的攻击记录条目
func (conn *AttackConnections) Find(protocol string, srcIp string, dstIp string, srcPort int, dstPort int, sdate int, edate int, page int, pagesize int) []AttackConnection {
    conditions := atkconnFmtQueryCondition(protocol, srcIp, dstIp, srcPort, dstPort, sdate, edate)
    query := fmt.Sprintf("select * from attack_connection%s limit %d offset %d;",
        conditions, pagesize, (page-1)*pagesize)
    rows, err := conn.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return nil
    }

    defer rows.Close()

    atkconnections := make([]AttackConnection, 0, pagesize)
    atkconnection := AttackConnection{}
    for rows.Next() {
        err = rows.Scan(&atkconnection.ID,
            &atkconnection.Protocol,
            &atkconnection.SrcIP,
            &atkconnection.DstIP,
            &atkconnection.SrcPort,
            &atkconnection.DstPort,
            &atkconnection.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        atkconnections = append(atkconnections, atkconnection)
    }
    return atkconnections
}

// 获取符合条件的攻击记录数量
func (conn *AttackConnections) Count(protocol string, srcIp string, dstIp string, srcPort int, dstPort int, sdate int, edate int) int {
    conditions := atkconnFmtQueryCondition(protocol, srcIp, dstIp, srcPort, dstPort, sdate, edate)
    query := fmt.Sprintf("select count(*) from attack_connection%s;", conditions)
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
*******************          连接日志            *******************
*******************************************************************/
func ListAtkConnectionsHandler(c *gin.Context) {
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
    response := AttackConnectionResponse{}
    response.Page = param.Page
    response.Pagesize = param.PageSize
    response.Total = attackconnections.Count(param.Protocol,
        param.SrcIP, param.DstIP,
        param.SrcPort, param.DstPort,
        param.Sdate, param.Edate)
    response.Data = attackconnections.Find(param.Protocol,
        param.SrcIP, param.DstIP,
        param.SrcPort, param.DstPort,
        param.Sdate, param.Edate,
        param.Page, param.PageSize)
    c.IndentedJSON(http.StatusOK, response)
}
