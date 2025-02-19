// Package cmd @Author lanpang
// @Date 2025/1/13 下午3:32:00
// @Desc
package cmd

import (
	"alarm/libs"
	"context"
	"database/sql"
	"log"
	"strings"
	"time"
)

func (alarm *Alarm) gather(queryTime string, config libs.DB) {

	// 创建 mysqlClinet
	mysqlClinet, err := libs.NewMysqlClient(config, "wshoto")
	if err != nil {
		log.Println("Failed to create mysql client. err:", err)
		return
	}
	defer func() {
		if mysqlClinet != nil {
			err := mysqlClinet.Close()
			if err != nil {
				return
			}
		}
	}()

	query := `
	SELECT
		content
	FROM
		dwd_qw_chat_msg_by_markdown_h
	WHERE
		msgtime > ?`

	rows, err := mysqlClinet.Query(query, queryTime)
	if err != nil {
		log.Println("数据查询失败. err:", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Failed info: %s \n", err)
		}
	}(rows)

	var dorisNum, conversationNum, defaultNum = 0, 0, 0
	// 处理查询结果
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			log.Printf("Failed info: %s \n", err)
		}
		handle := alarm.matchHandle(content)
		if handle == "doris" {
			dorisNum++
		} else if handle == "conversation" {
			conversationNum++
		} else {
			defaultNum++
		}
	}
	log.Printf("dorisNum: %d, conversationNum: %d, defaultNum: %d \n", dorisNum, conversationNum, defaultNum)
}

func (alarm *Alarm) matchHandle(input string) string {
	// 按换行符分割字符串
	lines := strings.Split(input, "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		switch firstLine {
		case "# Doris 巡检":
			doris := alarm.DorisHandler.Handle(input)
			alarm.Doris = append(alarm.Doris, doris)
			return "doris"
		case "# 会话数巡检":
			//fmt.Println(input)
			con := alarm.ConHandler.Handle(input)
			alarm.Conversation = append(alarm.Conversation, con)
			return "conversation"
		default:
			//defHandler(input)
			return "default"
		}
	}
	return "default"
}

func (alarm *Alarm) work() {
	// 处理 doris 巡检数据
	for _, d := range alarm.Doris {
		insertDorisLog(d, alarm.PGClient)
		checkDorisData(d, alarm.PGClient)
	}
	// 处理会话数巡检数据
	for _, c := range alarm.Conversation {
		insertConversationLog(c, alarm.PGClient)
		checkConversationData(c, alarm.PGClient)
	}
}

func Execute() {
	// 加载 TOML 配置文件
	_, err := loadConfig("config.toml")
	if err != nil {
		log.Fatal(err)
	}
	// 获取当天时间  queryTime := "2025-01-20"
	queryTime := time.Now().Format("2006-01-02")

	alarm := &Alarm{
		Doris:        []*Doris{},
		Conversation: []*Conversation{},
		Default:      []string{},
		DorisHandler: NewDorisHandler(),
		ConHandler:   NewConversationHandler(),
	}
	alarm.gather(queryTime, cfg.Doris)

	// 创建 pgclient
	pgClient, err := libs.NewPGClient(cfg.PG, "pomp")
	alarm.PGClient = pgClient
	if err != nil {
		log.Println("Failed to create pg client. err:", err)
		return
	}

	defer func() {
		if pgClient != nil {
			err := pgClient.Close(context.Background())
			if err != nil {
				log.Println("Failed to close pg client. err:", err)
			}
		}
	}()

	alarm.work()
}
