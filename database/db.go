package database

import (
    //"database/sql"
    "fmt"
    "net/url"
    "vesaliusdr/config"
    "vesaliusdr/utils"

    _ "github.com/microsoft/go-mssqldb"
    "github.com/jmoiron/sqlx"
)

var dbVar *sqlx.DB

func SetDb(db *sqlx.DB) {
    dbVar = db
}

func GetDb() *sqlx.DB {
    if dbVar == nil {
        ConnectDB()
    }

    return dbVar
}

func ConnectDB() {
    port := 1433
    query := url.Values{}
    u := &url.URL{
        Scheme:   "sqlserver",
        User:     url.UserPassword(config.Config("db.username"), config.Config("db.pwd")),
        Host:     fmt.Sprintf("%s:%d", config.Config("db.host"), port),
        RawQuery: query.Encode(),
    }
    db, err := sqlx.Open("sqlserver", u.String())
    if err != nil {
        utils.LogError(err)
    } else {
        db.SetMaxOpenConns(10)
        db.SetMaxIdleConns(5)
        SetDb(db)
        utils.LogInfo("Connection Opened to Database")
    }
}

func CloseDB() {
    dbVar.Close()
}