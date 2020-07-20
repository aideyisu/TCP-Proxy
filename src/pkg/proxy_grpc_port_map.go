package tcpmax

import (
    "database/sql"
    "fmt"
    "log"
    "time"
    "errors"
    "os/exec"

    // just import
    _ "github.com/mattn/go-sqlite3"
)

type ProxyGrpcPort struct {
    ID         int64
    IP         string
    Port       int
    OffSet     int
    CreateTime int
}

type ProxyProcess struct {
    GrpcPort *ProxyGrpcPort
    Cmd      *exec.Cmd
}

type ProxyGrpcPortMaps struct {
    Db            *sql.DB
    GrpcPortMap   map[string](*ProxyProcess)
    GrpcPortIDMap BitMap
}

func InitProxyGrpcPort(dbfile string) *ProxyGrpcPortMaps {
    proxygrpcportmaps := &ProxyGrpcPortMaps{}
    db, err := sql.Open("sqlite3", dbfile)
    if err != nil {
        log.Fatal("Openning dbfile", err)
    }
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS proxy_grpc_port_map(
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            ip CHAR(50) NOT NULL,
            port INTEGER  NOT NULL,
            offset INTEGER NOT NULL UNIQUE,
            create_time INTEGER NOT NULL,
            UNIQUE(ip,port));
        `)
    if err != nil {
        db.Close()
        log.Fatal("Openning dbfile", err)
    }
    proxygrpcportmaps.Db = db
    proxygrpcportmaps.GrpcPortMap = make(map[string](*ProxyProcess))
    proxygrpcportmaps.GrpcPortIDMap.Init()
    return proxygrpcportmaps
}

func (pgpm *ProxyGrpcPortMaps) Add(ip string, port int) (*ProxyProcess, error) {
    cardID := pgpm.GrpcPortIDMap.GetID()
    if cardID == -1 {
        return nil, errors.New("Proxy grpc port resource exhausted.")
    }
    query := fmt.Sprintf(`insert into proxy_grpc_port_map(ip, port, offset, create_time) values ('%s', '%d', '%d', '%d');`,
        ip, port, cardID, int32(time.Now().Unix()))
    _, err := pgpm.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
        return nil, err
    }
    gpm := pgpm.Get(ip, port)
    if gpm != nil {
        pgpm.GrpcPortMap[fmt.Sprintf("%s:%d", ip, port)] = &ProxyProcess{gpm, nil}
    }
    return pgpm.GrpcPortMap[fmt.Sprintf("%s:%d", ip, port)], nil
}

func (pgpm *ProxyGrpcPortMaps) Del(ip string, port int) error {
    query := fmt.Sprintf("delete from proxy_grpc_port_map where ip='%s' and port='%d';", ip, port)
    _, err := pgpm.Db.Exec(query)
    if err != nil {
        log.Print(query, err)
        return err
    }
    delete(pgpm.GrpcPortMap, fmt.Sprintf("%s:%d", ip, port))
    return nil
}

func (pgpm *ProxyGrpcPortMaps) Get(ip string, port int) *ProxyGrpcPort {
    query := fmt.Sprintf("select * from proxy_grpc_port_map where ip='%s' and port='%d';", ip, port)
    rows, err := pgpm.Db.Query(query)
    if err != nil {
        log.Print(query, err)
        return nil
    }
    defer rows.Close()

    grpcport := ProxyGrpcPort{}
    if rows.Next() {
        err = rows.Scan(&grpcport.ID, &grpcport.IP, &grpcport.Port, &grpcport.OffSet, &grpcport.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        return &grpcport
    }

    return nil
}

func (pgpm *ProxyGrpcPortMaps) Count() int {
    query := fmt.Sprintf("select count(*) from proxy_grpc_port_map;")
    rows, err := pgpm.Db.Query(query)
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

func (pgpm *ProxyGrpcPortMaps) Find(limit int, offset int) []ProxyGrpcPort {
    query := fmt.Sprintf("select * from proxy_grpc_port_map limit %d offset %d;", limit, offset)
    rows, err := pgpm.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return nil
    }

    defer rows.Close()

    grpcports := make([]ProxyGrpcPort, 0, limit)
    grpcport := ProxyGrpcPort{}
    for rows.Next() {
        err = rows.Scan(&grpcport.ID, &grpcport.IP, &grpcport.Port, &grpcport.OffSet, &grpcport.CreateTime)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        grpcports = append(grpcports, grpcport)
    }
    return grpcports
}

func (pgpm *ProxyGrpcPortMaps) Loads(){
    count := pgpm.Count()
	for i := 0; i < count; i++ {
		pgpms := pgpm.Find(1, i)
        pgpm.GrpcPortIDMap.Set(uint16(pgpms[0].OffSet))
        pgpm.GrpcPortMap[fmt.Sprintf("%s:%d", pgpms[0].IP, pgpms[0].Port)] = &ProxyProcess{&pgpms[0], nil}
	}
}

func (pgpm *ProxyGrpcPortMaps) GetProcess(ip string, port int) *ProxyProcess {
    Value := pgpm.GrpcPortMap[fmt.Sprintf("%s:%d", ip, port)]
    return Value
}
