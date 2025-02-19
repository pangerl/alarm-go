// Package cmd @Author lanpang
// @Date 2025/1/20 下午3:50:00
// @Desc
package cmd

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Message struct {
	CorpName            string
	CurrentMessageNum   int64
	YesterdayMessageNum int64
}

type Conversation struct {
	ProjectName    string
	InspectionDate string
	CorpMessage    []*Message
}

type ConversationHandler struct {
	reProjectName    *regexp.Regexp
	reInspectionDate *regexp.Regexp
	reCorpMessage    *regexp.Regexp
}

func NewConversationHandler() *ConversationHandler {
	return &ConversationHandler{
		reProjectName:    regexp.MustCompile(`\*\*项目名称：.*>(.*?)</font>`),
		reInspectionDate: regexp.MustCompile(`\*\*巡检时间：.*>(.*?)</font>`),
		reCorpMessage:    regexp.MustCompile(`> 企业名称：.*>(.*?)</font>\s*> 当前拉取会话数：.*>(\d+)</font>\s*> 昨天拉取会话数：.*>(\d+)</font>`),
	}
}

func (c *ConversationHandler) Handle(content string) *Conversation {

	conversation := &Conversation{}
	// 提取字段
	projectName := extractMatch(c.reProjectName, content)
	conversation.ProjectName = projectName
	inspectionDate := extractMatch(c.reInspectionDate, content)
	conversation.InspectionDate = inspectionDate

	matches := c.reCorpMessage.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) == 4 {
			currentMessageNum, _ := strconv.ParseInt(match[2], 10, 64)
			yesterdayMessageNum, _ := strconv.ParseInt(match[3], 10, 64)
			message := &Message{
				CorpName:            match[1],
				CurrentMessageNum:   currentMessageNum,
				YesterdayMessageNum: yesterdayMessageNum,
			}
			conversation.CorpMessage = append(conversation.CorpMessage, message)
		}
	}
	return conversation
}

func insertConversationLog(c *Conversation, pgClient *pgx.Conn) {
	// 查询项目ID
	projectId := selectProjectId(pgClient, c.ProjectName)
	// 插入数据到数据库
	insertQuery := `
		INSERT INTO public.qw_project_conversation_log (project_name, project_id, corp_name, current_message_num, yesterday_message_num, inspection_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (corp_name, inspection_date) DO UPDATE
		SET
		    project_id = EXCLUDED.project_id,
		    corp_name = EXCLUDED.corp_name,
		    current_message_num = EXCLUDED.current_message_num,
		    yesterday_message_num = EXCLUDED.yesterday_message_num
		RETURNING id;`
	for _, message := range c.CorpMessage {
		_, err := pgClient.Exec(context.Background(), insertQuery,
			c.ProjectName, projectId, message.CorpName, message.CurrentMessageNum, message.YesterdayMessageNum, c.InspectionDate)
		if err != nil {
			log.Println(c.ProjectName, "插入数据失败: ", err)
		} else {
			//log.Println(c.ProjectName, "数据插入成功!")
		}
	}
}

func checkConversationData(c *Conversation, pgClient *pgx.Conn) {

	if checkCorpMessage(c) {
		// 查询项目运维
		operationer := selectProjectOperationer(pgClient, selectProjectId(pgClient, c.ProjectName))
		toList := getToList(operationer)
		var builder strings.Builder
		// 构建邮件内容
		builder.WriteString(conversationMailHead)
		for _, m := range c.CorpMessage {
			corpData := fmt.Sprintf("<td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td></tr>", c.ProjectName, m.CorpName, c.InspectionDate, m.YesterdayMessageNum, m.CurrentMessageNum)
			builder.WriteString(corpData)
		}
		builder.WriteString("</table><br>请查收！</br><br>Send By notify@wshoto.com </br>（自动发送请勿回复）")
		// 发送告警
		mailAlert("Conversation", builder.String(), toList)
	}
}

func checkCorpMessage(c *Conversation) bool {

	for _, m := range c.CorpMessage {
		if m.CurrentMessageNum == 0 && m.YesterdayMessageNum == 0 {
			log.Println("会话数异常!", m)
			return true
		}
	}
	return false
}
