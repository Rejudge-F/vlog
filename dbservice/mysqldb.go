package dbservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"time"
)

type MysqlDB struct {
	Pool    *sql.DB // Mysql连接池
	DBValid bool    // Mysql是否可用
}

// 检测数据库连接
func (mysqlDB *MysqlDB) Ping() (bool, error) {
	err := mysqlDB.Pool.Ping()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 定时检测数据库是否正常
func (mysqlDB *MysqlDB) CheckHealth() {
	t := viper.GetInt("Mysql.MysqlHealthCheckTimer")
	if t == 0 {
		t = 1
	}

	ticker := time.NewTicker(time.Duration(t) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
			mysqlDB.DBValid, _ = mysqlDB.Ping()
		}
	}
}

func (mysqlDB *MysqlDB) Get() (*sql.Conn, error){
	if(mysqlDB.Pool == nil) {
		mysqlDB.InitPool()
		go mysqlDB.CheckHealth()
	}
	return mysqlDB.Pool.Conn(context.Background())
}

func (mysqlDB *MysqlDB) Close() error {
	if nil == mysqlDB.Pool {
		return nil
	}
	return mysqlDB.Pool.Close()
}

func (mysqlDB *MysqlDB) InitPool() error {
	if mysqlDB.Pool != nil {
		return errors.New("initPool repeat")
	}

	addr := viper.GetString("Mysql.Addr")
	user := viper.GetString("Mysql.User")
	pwd := viper.GetString("Mysql.Password")
	rawUrl := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8&parseTime=true&loc=Local", user, pwd, addr)
	maxConns := viper.GetInt("Mysql.MaxConns")
	maxIdle := viper.GetInt("Mysql.MaxIdle")
	connLifeTime := viper.GetInt("Mysql.ConnLifeTime")

	pool, err := sql.Open("mysql", rawUrl)
	if err != nil {
		return err
	}

	pool.SetMaxOpenConns(maxConns)
	pool.SetMaxIdleConns(maxIdle)
	pool.SetConnMaxLifetime(time.Duration(connLifeTime) * time.Second)

	mysqlDB.Pool = pool
	return nil
}
