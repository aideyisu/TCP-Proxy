package tcpmax

import (
    "database/sql"
    "math/rand"
    "fmt"
    "log"
    "time"

    jwt "github.com/appleboy/gin-jwt/v2"
    "github.com/gin-gonic/gin"

    // just import
    _ "github.com/mattn/go-sqlite3"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandStringBytesRmndr 生成长度为n的随机字符串
func RandStringBytesRmndr(n int) string {
    rand.Seed(time.Now().UnixNano())
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
    }
    return string(b)
}

type login struct {
    Username string `form:"username" json:"username" binding:"required"`
    Password string `form:"password" json:"password" binding:"required"`
}

// User 用户信息
type User struct {
    ID                 string
    UserName           string
    Key                string
    CreateTime         int64
    Desc               string
}

// Users 用户管理
type Users struct {
    Db       *sql.DB
    UserInfo map[string](*User)
}

// InitUsers 创建用户管理
func InitUsers(dbfile string) *Users {
    users = &Users{}

    db, err := sql.Open("sqlite3", dbfile)
    if err != nil {
        log.Fatal("Openning dbfile", err)
    }
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS user (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username VARCHAR(16) NOT NULL UNIQUE,
        key VARCHAR(32) NOT NULL,
        create_time INTEGER NOT NULL,
        desc VARCHAR(100));
        `)
    if err != nil {
        log.Fatal("Creating table:", err)
    }

    users.Db = db
    users.UserInfo = map[string](*User){}
    return users
}

// Auth 用户认证
func (us *Users) Auth(username string, key string) bool {
    query := fmt.Sprintf("select username, key from user where username='%s' and key='%s';", username, key)
    rows, err := us.Db.Query(query)
    if err != nil {
        log.Printf(err.Error())
        return false
    }
    defer rows.Close()

    if rows.Next() {
        return true
    }
    return false
}

// Add 添加用户
func (us *Users) Add(username string, desc string) *User {
    key := RandStringBytesRmndr(20)
    query := fmt.Sprintf(`insert into user(username, key, create_time, desc) values ('%s', '%s', %d, '%s');`,
        username, key, int32(time.Now().Unix()), desc)
    _, err := us.Db.Exec(query)
    if err != nil {
        fmt.Print(query, err)
        return nil
    }

    return us.Get(username)
}

// Del 删除用户
func (us *Users) Del(username string) {
    query := fmt.Sprintf("delete from user where username='%s';", username)
    _, err := us.Db.Exec(query)
    if err != nil {
        fmt.Print(query, err)
    }
}

// SetDesc 更新用户描述信息
func (us *Users) SetDesc(username string, desc string) *User {
    query := fmt.Sprintf("update user set desc='%s' where username='%s';", desc, username)
    _, err := us.Db.Exec(query)
    if err != nil {
        fmt.Print(query, err)
        return nil
    }
    return us.Get(username)
}

// UpdateKey 更新用户密钥
func (us *Users) UpdateKey(username string) *User {
    key := RandStringBytesRmndr(20)
    query := fmt.Sprintf("update user set key='%s' where username='%s';", key, username)
    _, err := us.Db.Exec(query)
    if err != nil {
        fmt.Print(query, err)
        return nil
    }
    return us.Get(username)
}

// Get 根据username获取用户数据
func (us *Users) Get(username string) *User {
    query := fmt.Sprintf("select * from user where username='%s';", username)
    rows, err := us.Db.Query(query)
    if err != nil {
        log.Printf(query, err.Error())
        return nil
    }

    defer rows.Close()

    user := User{}
    if rows.Next() {
        err = rows.Scan(&user.ID, &user.UserName, &user.Key, &user.CreateTime, &user.Desc)
        if err != nil {
            log.Printf(err.Error())
            return nil
        }
        return &user
    }
    return nil
}

var identityKey = "id"

func AuthorizatorHandler(data interface{}, c *gin.Context) bool {
    // 这里可以做权限控制
    if v, ok := data.(*User); ok && v.UserName == "admin" {
        return true
    }

    return false
}

func PayloadFuncHandler(data interface{}) jwt.MapClaims {
    if v, ok := data.(*User); ok {
        return jwt.MapClaims{
            identityKey: v.UserName,
        }
    }
    return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
    claims := jwt.ExtractClaims(c)
    return &User{
        UserName: claims[identityKey].(string),
    }
}

func AuthenticatorHandler(c *gin.Context) (interface{}, error) {
    var loginVals login
    if err := c.ShouldBind(&loginVals); err != nil {
        return "", jwt.ErrMissingLoginValues
    }
    username := loginVals.Username
    password := loginVals.Password

    if users.Auth(username, password) == true {
        return users.Get(username), nil
    }
    return nil, jwt.ErrFailedAuthentication
}
