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

// 白名单
type WhiteList struct {
    ID         int64  `json:"id"`
    IP         string `json:"ip"`
    Creator    string `json:"creator"`
    Status     string `json:"status"`
    CreateTime uint32 `json:"timestamp"`
}

type whiteListResponse struct {
    Page     int         `json:"page"`
    Pagesize int         `json:"pagesize"`
    Total    int         `json:"total"`
    Data     []WhiteList `json:"data"`
}

// 白名单管理
type WhiteLists struct {
    Db            *sql.DB
    WhiteListInfo map[string](*WhiteList)
}

// InitWhiteListTable 创建白名单
func InitWhiteListTable(dbfile string) *WhiteLists {
    whitelists = &WhiteLists{}

    db, err := sql.Open("sqlite3", dbfile)
    if err != nil {
        log.Fatal("Openning dbfile", err)
    }
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS whitelist(
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            ip CHAR(50) NOT NULL UNIQUE,
            status CHAR(50) NOT NULL,
            creator CHAR(50) NOT NULL,
            create_time INTEGER NOT NULL);
        `)
    if err != nil {
        log.Fatal("Creating table:", err)
    }

    whitelists.Db = db
    whitelists.WhiteListInfo = map[string](*WhiteList){}
    return whitelists
}

// Add 添加白名单
func (wls *WhiteLists) Add(ip string, status string, creator string) (*WhiteList, error) {
    query := fmt.Sprintf(`insert into whitelist(ip, status, creator, create_time) values ('%s', '%s', '%s', %d);`,
        ip, status, creator, int32(time.Now().Unix()))
    _, err := wls.Db.Exec(query)
    if err != nil {
        //log.Print(query, err)
        log.Print(err)
        return nil, err
    }

    return wls.Get(ip), nil
}

// Del 通过ip删除白名单
func (wls *WhiteLists) Del(ip string) {
    query := fmt.Sprintf("delete from whitelist where ip='%s';", ip)
    _, err := wls.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
    }
}

// Modify 通过ID修改白名单
func (wls *WhiteLists) Modify(id int64, ip string) (*WhiteList, error) {
    query := fmt.Sprintf("update whitelist set ip='%s' where id='%d';", ip, id)
    _, err := wls.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
        return nil, err
    }
    return wls.Get(ip), nil
}

// Get 通过ip查询白名单
func (wls *WhiteLists) Get(ip string) *WhiteList {
    query := fmt.Sprintf("select * from whitelist where ip='%s';", ip)
    rows, err := wls.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return nil
    }
    defer rows.Close()

    whitelist := WhiteList{}
    if rows.Next() {
        err = rows.Scan(&whitelist.ID, &whitelist.IP, &whitelist.Status, &whitelist.Creator, &whitelist.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        return &whitelist
    }

    return nil
}

// getByID 通过ID查询白名单
func (wls *WhiteLists) getByID(id int64) *WhiteList {
    query := fmt.Sprintf("select * from whitelist where id='%d';", id)
    rows, err := wls.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return nil
    }
    defer rows.Close()

    whitelist := WhiteList{}
    if rows.Next() {
        err = rows.Scan(&whitelist.ID, &whitelist.IP, &whitelist.Status, &whitelist.Creator, &whitelist.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        return &whitelist
    }

    return nil
}

// ChangeStatus 启动和禁用上游IP白名单
func (wls *WhiteLists) ChangeStatus(status string, id int64) (*WhiteList, error) {
    query := fmt.Sprintf("update whitelist set status='%s' where id='%d';", status, id)
    _, err := wls.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
        return nil, err
    }
    return wls.getByID(id), nil
}

func wlFmtQueryCondition(ip string, creator string, status string, sdate int, edate int) string {
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
    statusQueryCondition := ""
    if status != "" {
        statusQueryCondition = fmt.Sprintf("%s status='%s'", andstr, status)
        andstr = " AND"
    }
    dateQueryCondition := ""
    if sdate != edate {
        dateQueryCondition = fmt.Sprintf("%s create_time BETWEEN '%d' AND '%d'", andstr, sdate, edate)
    }
    conditions := ""
    if ip != "" || creator != "" || status != "" || sdate != edate {
        conditions = fmt.Sprintf(" where%s%s%s%s",
            ipQueryCondition,
            creatorQueryCondition,
            statusQueryCondition,
            dateQueryCondition)
    }
    return conditions
}

// 查询符合条件白名单条目
func (wls *WhiteLists) Find(ip string, creator string, status string, sdate int, edate int, page int, pagesize int) []WhiteList {
    conditions := wlFmtQueryCondition(ip, creator, status, sdate, edate)
    query := fmt.Sprintf("select * from whitelist%s limit %d offset %d;",
        conditions, pagesize, (page-1)*pagesize)
    log.Printf(query)
    rows, err := wls.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return nil
    }

    defer rows.Close()

    whitelists := make([]WhiteList, 0, pagesize)
    whitelist := WhiteList{}
    for rows.Next() {
        err = rows.Scan(&whitelist.ID, &whitelist.IP, &whitelist.Status, &whitelist.Creator, &whitelist.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        whitelists = append(whitelists, whitelist)
    }
    return whitelists
}

// 获取符合条件白名单数量
func (wls *WhiteLists) Count(ip string, creator string, status string, sdata int, edata int) int {
    conditions := wlFmtQueryCondition(ip, creator, status, sdata, edata)
    query := fmt.Sprintf("select count(*) from whitelist%s;", conditions)
    rows, err := wls.Db.Query(query)
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
*******************            IP白名单          *******************
*******************************************************************/
func AddWhiteListHandler(c *gin.Context) {
    type addWhiteIP struct {
        IP      string `form:"ip" json:"ip" binding:"required,ipv4"`
        Creator string `form:"creator" json:"creator" binding:"required"`
        Status  string `form:"status" json:"status" binding:"required,oneof=启用 停用"`
    }
    var param addWhiteIP
    if err := c.ShouldBind(&param); err != nil {
        c.String(http.StatusBadRequest, err.Error())
        return
    }
    cfg, err := whitelists.Add(param.IP, param.Status, param.Creator)
    if cfg != nil {
        c.JSON(200, gin.H{
            "status_code": 200,
            "reason":      "The configuration was added successfully.",
        })
    } else {
        StatusCode, Reason := checkErr(err)
        c.JSON(int(StatusCode), gin.H{
            "status_code": StatusCode,
            "reason":      Reason,
        })
    }
}

// 修改白名单配置
func ModifyWhiteListHandler(c *gin.Context) {
    type Param struct {
        ID int64  `form:"id" json:"id" binding:"required,numeric,gte=0"`
        IP string `form:"ip" json:"ip" binding:"required,ipv4"`
    }
    var param Param
    if err := c.ShouldBind(&param); err != nil {
        c.String(http.StatusBadRequest, err.Error())
        return
    }
    cfg, err := whitelists.Modify(param.ID, param.IP)
    if cfg != nil {
        c.JSON(200, gin.H{
            "status_code": 200,
            "reason":      "The specified configuration was modified successfully.",
        })
    } else {
        StatusCode, Reason := checkErr(err)
        c.JSON(int(StatusCode), gin.H{
            "status_code": StatusCode,
            "reason":      Reason,
        })
    }
}

// 获取符合条件的白名单列表
func ListWhiteListHandler(c *gin.Context) {
    type Param struct {
        Page       int    `form:"page" json:"page" binding:"required,numeric,gt=0"`
        PageSize   int    `form:"pagesize" json:"pagesize" binding:"required,numeric,gt=0,lte=100"`
        Creator    string `form:"creator" json:"creator" binding:"omitempty"`
        IP         string `form:"ip" json:"ip" binding:"omitempty,ipv4"`
        Status     string `form:"status" json:"status" binding:"omitempty,oneof=启用 停用"`
        Sdate      int    `form:"sdate" json:"sdate" binding:"omitempty,numeric,ltfield=Edate"`
        Edate      int    `form:"edate" json:"edate" binding:"omitempty,numeric,gtfield=Sdate"`
    }
    var param Param

    if err := c.ShouldBindQuery(&param); err != nil {
        c.String(http.StatusBadRequest, err.Error())
        return
    }

    response := whiteListResponse{}
    response.Page = param.Page
    response.Pagesize = param.PageSize
    response.Total = whitelists.Count(param.IP, param.Creator, param.Status, param.Sdate, param.Edate)
    response.Data = whitelists.Find(param.IP, param.Creator, param.Status, param.Sdate, param.Edate, param.Page, param.PageSize)
    c.IndentedJSON(http.StatusOK, response)
}

// 删除监听配置
func DelWhiteListHandler(c *gin.Context) {
    type Param struct {
        IP string `form:"ip" json:"ip" binding:"required,ipv4"`
    }
    var param Param
    if err := c.ShouldBind(&param); err != nil {
        c.String(http.StatusBadRequest, err.Error())
        return
    }
    whitelists.Del(param.IP)
    c.JSON(200, gin.H{
        "status_code": 200,
        "reason":      "success",
    })
}

// ChangeWhiteStatusHandler 启动和禁用上有IP白名单的Web接口
func ChangeWhiteStatusHandler(c *gin.Context) {
    type TurnWhiteStatus struct {
        ID     int64  `form:"id" json:"id" binding:"required,numeric,gte=0"`
        Status string `form:"status" json:"status" binding:"required,oneof=启用 停用"`
    }
    var param TurnWhiteStatus

    if err := c.ShouldBind(&param); err != nil {
        c.String(http.StatusBadRequest, err.Error())
        return
    }
    cfg, err := whitelists.ChangeStatus(param.Status, param.ID)
    if cfg != nil {
        c.JSON(200, gin.H{
            "status_code": 200,
            "reason":      "The Status of WhiteList was changed successfully.",
        })
    } else {
        StatusCode, Reason := checkErr(err)
        c.JSON(int(StatusCode), gin.H{
            "status_code": StatusCode,
            "reason":      Reason,
        })
    }
}
