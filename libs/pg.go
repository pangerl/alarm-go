// Package libs @Author lanpang
// @Date 2024/8/8 下午1:43:00
// @Desc
package libs

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/jackc/pgx/v5"
)

func NewPGClient(conf DB, dbName string) (*pgx.Conn, error) {
	connString := connStr(conf, dbName)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Printf("Failed to connect to database %s: %s\n", dbName, err)
		return nil, err
	}
	log.Printf("%s 数据库连接成功！", dbName)
	return conn, nil
}

func connStr(conf DB, db string) string {
	scheme := map[bool]string{true: "require", false: "disable"}[conf.Sslmode]
	// 对密码进行 URL 编码
	encodedPassword := url.QueryEscape(conf.Password)
	str := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		conf.Username, encodedPassword, conf.Ip, conf.Port, db, scheme)
	return str
}
