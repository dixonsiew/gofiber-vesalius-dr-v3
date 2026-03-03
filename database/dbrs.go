package database

import (
    //"database/sql"
    "strconv"
    "vesaliusdr/config"
    "vesaliusdr/utils"

    "github.com/sijms/go-ora/v2"
    "github.com/jmoiron/sqlx"
)

var dbrsVar *sqlx.DB

func SetDbrs(db *sqlx.DB) {
    dbrsVar = db
}

func GetDbrs() *sqlx.DB {
    if dbrsVar == nil {
        ConnectDBRs()
    }

    return dbrsVar
}

func ConnectDBRs() {
    p := config.Config("db.rs.port")
    port, _ := strconv.Atoi(p)
    connStr := go_ora.BuildUrl(config.Config("db.rs.host"), port, config.Config("db.rs.service"), config.Config("db.rs.username"), config.Config("db.rs.pwd"), nil)
    db, err := sqlx.Open("oracle", connStr)
    // dsn := fmt.Sprintf(`user="%s" password="%s" connectString="%s" heterogeneousPool=false standaloneConnection=false`, config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_URL"))
    // fmt.Println(dsn)
    // DB, err := sql.Open("godror", dsn)

    if err != nil {
        utils.LogError(err)
    } else {
        db.SetMaxOpenConns(10)
        db.SetMaxIdleConns(5)
        SetDbrs(db)
        utils.LogInfo("Connection Opened to Database")
    }
}

func CloseDBRs() {
    dbrsVar.Close()
}