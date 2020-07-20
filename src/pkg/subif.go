package tcpmax

import (
    "database/sql"
    "log"
    "fmt"
    "time"
    "os/exec"
    "errors"
    // just import
    _ "github.com/mattn/go-sqlite3"
)

type SubInterface struct {
    ID         int64
    IP         string
    Number     int
    CreateTime int
}

type SubInterfaces struct {
    Db *sql.DB
    IfName string
    SubItfMap map[string](*SubInterface)
    IfBitMap BitMap
}

func (sif *SubInterfaces) addDB(ip string, number int) (*SubInterface, error) {
    query := fmt.Sprintf(`insert into if_manage (ip, number, create_time) values ('%s', '%d', '%d');`,
        ip, number, int32(time.Now().Unix()))
    _, err := sif.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
        return nil, err
    }
    return sif.getDB(ip, number), nil
}

func (sif *SubInterfaces) delDB(ip string, number int) error {
    query := fmt.Sprintf(`delete from if_manage where ip='%s' and number='%d';`, ip, number)
    _, err := sif.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
        return err
    }
    return nil
}

func (sif *SubInterfaces) getDB(ip string, number int) *SubInterface {
    query := fmt.Sprintf("select * from if_manage where ip='%s' and number='%d';", ip, number)
    rows, err := sif.Db.Query(query)
    if err != nil {
        log.Print(query, err)
        return nil
    }
    defer rows.Close()

    subift := SubInterface{}
    if rows.Next() {
        err = rows.Scan(&subift.ID, &subift.IP, &subift.Number, &subift.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        return &subift
    }
    return nil
}

func (sif *SubInterfaces) countDB() int {
    query := fmt.Sprintf("select count(*) from if_manage;")
    rows, err := sif.Db.Query(query)
    if err != nil {
        log.Print(query, err)
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

func (sif *SubInterfaces) findDB(limit int, offset int) []SubInterface {
    query := fmt.Sprintf("select * from if_manage limit %d offset %d;", limit, offset)
    rows, err := sif.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return nil
    }

    defer rows.Close()

    subifts := make([]SubInterface, 0, limit)
    subift := SubInterface{}
    for rows.Next() {
        err = rows.Scan(&subift.ID, &subift.IP, &subift.Number, &subift.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        subifts = append(subifts, subift)
    }
    return subifts
}

func (sif *SubInterfaces) configSubIf(ip string) (*SubInterface, error) {
    id := sif.IfBitMap.GetID()
    if id == -1{
        return nil, errors.New("The interface has run out.") 
    }
    ifname := fmt.Sprintf("%s:%d", sif.IfName, id)
    cmd := exec.Command("ifconfig", ifname, ip, "netmask", "255.255.255.0", "up")
    if err := cmd.Run(); err != nil {
        return nil, err
    }
    subift, err := sif.addDB(ip, id)
    if err != nil {
        cmd = exec.Command("ifconfig", ifname, "down")
        cmd.Run()
        return nil, err
    }
    sif.SubItfMap[ip] = subift
    return subift, nil
}

func (sif *SubInterfaces) delSubIf(ip string) error {
    subift := sif.SubItfMap[ip]
    ifname := fmt.Sprintf("%s:%d", sif.IfName, subift.Number)
    cmd := exec.Command("ifconfig", ifname, "down")
    if err := cmd.Run(); err != nil {
        return err
    }
    sif.delDB(ip, subift.Number)
    return nil
}

func (sif *SubInterfaces) Add(ip string) (*SubInterface, error) {
    subift, ok := sif.SubItfMap[ip]
    if (ok) {
        return subift, nil
    } else {
        subift, err := sif.configSubIf(ip)
        return subift, err
    }
}

func (sif *SubInterfaces) Del(ip string) {
    subift, ok := sif.SubItfMap[ip]
    if (ok) {
        // 删除失败也没关系，下次这个ip可以被重新配置覆盖
        sif.delSubIf(ip)
        sif.IfBitMap.Clear(subift.Number)
        delete(sif.SubItfMap, ip)
    }
}

func (sif *SubInterfaces) Modify(oldip string, newip string) (*SubInterface, error ){
    subift, ok := sif.SubItfMap[oldip]
    if (ok) {
        sif.Del(oldip)
        subift, err := sif.Add(newip)
        return subift, err
    }
    return subift, errors.New("Interface not found.")
}

func (sif *SubInterfaces) Get(ip string) *SubInterface {
    subift, ok := sif.SubItfMap[ip]
    if (ok) {
        return subift
    }
    return nil
}

func (sif *SubInterfaces) Loads() {
    count := sif.countDB()
    for i := 0; i < count; i++ {
        subifts := sif.findDB(1, i)
        sif.SubItfMap[subifts[0].IP] = &subifts[0]
        sif.IfBitMap.Set(uint16(subifts[0].Number))
        // 配置网卡
        ifname := fmt.Sprintf("%s:%d", sif.IfName, subifts[0].Number)
        cmd := exec.Command("ifconfig", ifname, subifts[0].IP, "netmask", "255.255.255.0", "up")
        if err := cmd.Run(); err != nil {
            log.Fatal("Config interface error", err)
        }
    }
}

func (sif *SubInterfaces) Init(ifname string, dbfile string) {
    db, err := sql.Open("sqlite3", dbfile)
    if err != nil {
        log.Fatal("Openning dbfile", err)
    }
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS if_manage(
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            ip CHAR(50) NOT NULL UNIQUE,
            number INTEGER NOT NULL UNIQUE,
            create_time INTEGER NOT NULL
            );
        `)
    if err != nil {
        db.Close()
        log.Fatal("Openning dbfile", err)
    }

    sif.IfName = ifname
    sif.Db = db
    sif.SubItfMap = make(map[string](*SubInterface))
    sif.IfBitMap.Init()
    sif.Loads()
}
