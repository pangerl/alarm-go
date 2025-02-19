// Package cmd @Author lanpang
// @Date 2025/1/14 下午2:37:00
// @Desc
package cmd

import (
	"alarm/libs"
	"github.com/jackc/pgx/v5"
)

var (
	cfg *CfgType
)

type CfgType struct {
	Doris libs.DB   `toml:"doris"`
	PG    libs.DB   `toml:"pg"`
	Mail  libs.Mail `toml:"mail"`
}

type Alarm struct {
	Doris        []*Doris
	Conversation []*Conversation
	Default      []string
	DorisHandler *DorisHandler
	ConHandler   *ConversationHandler
	PGClient     *pgx.Conn
}

const dorisMailHead = // 定义邮件表格样式
`<style>
		table {
			border-collapse: collapse;
		}
		th {
			background-color: #007fff;
			color: white;
		}
		table, th, td {
			border: 1px solid black;
			padding: 5px;
			text-align: left;
		}
		</style>
		<table>
			<tr>
				<th>企业名称</th>
				<th>巡检日期</th>
				<th>Job失败数</th>
				<th>员工统计</th>
				<th>使用分析</th>
				<th>客户群统计</th>
				<th>BE节点总数</th>
				<th>BE存活总数</th>
			</tr>`

const conversationMailHead = // 定义邮件表格样式
`<style>
		table {
			border-collapse: collapse;
		}
		th {
			background-color: #007fff;
			color: white;
		}
		table, th, td {
			border: 1px solid black;
			padding: 5px;
			text-align: left;
		}
		</style>
		<table>
			<tr>
				<th>企业名称</th>
				<th>租户名称</th>
				<th>巡检日期</th>
				<th>昨天会话记录数</th>
				<th>当前会话记录数</th>
			</tr>`
