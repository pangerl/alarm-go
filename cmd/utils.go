// Package cmd @Author lanpang
// @Date 2025/1/16 上午10:55:00
// @Desc
package cmd

import (
	"alarm/libs"
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jackc/pgx/v5"

	"github.com/mozillazg/go-pinyin"
	"log"
	"regexp"
	"strings"
)

// 提取字段的公共方法
func extractMatch(re *regexp.Regexp, content string) string {
	match := re.FindStringSubmatch(content)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

// 字符串转整型的辅助函数
func parseInt(value string) int {
	var result int
	_, err := fmt.Sscanf(value, "%d", &result)
	if err != nil {
		return 0
	}
	return result
}

// isChinese 判断字符串是否是中文
func isChinese(s string) bool {
	// 匹配中文的正则表达式
	chineseRegex := regexp.MustCompile(`^\p{Han}+$`)
	return chineseRegex.MatchString(s)
}

// toPinyin 将中文字符串转换为拼音
func toPinyin(s string) string {
	// pinyin 包提供的选项
	a := pinyin.NewArgs()
	py := pinyin.Pinyin(s, a)

	// 将拼音数组拼接为字符串
	var result []string
	for _, p := range py {
		result = append(result, strings.Join(p, ""))
	}
	return strings.Join(result, "")
}

func getToList(p string) []string {
	var toList []string
	var pinyinStr string
	if isChinese(p) {
		pinyinStr = toPinyin(p)
	} else {
		pinyinStr = p
	}
	toList = append(toList, pinyinStr+"@wshoto.com")
	toList = append(toList, cfg.Mail.AddTo...)
	return toList
}

func selectProjectId(pgClient *pgx.Conn, projectName string) int {
	var projectId int
	query := "SELECT project_id FROM public.qw_project WHERE project_name LIKE $1 AND project_type = 'prod';"
	err := pgClient.QueryRow(context.Background(), query, "%"+projectName+"%").Scan(&projectId)
	if err != nil {
		log.Println(projectName, "查询项目ID失败:", err)
		return 999999
	}
	return projectId
}

func selectProjectOperationer(pgClient *pgx.Conn, projectId int) string {
	var operationer string
	query := "SELECT operationer FROM public.qw_project_member WHERE project_id = $1;"
	err := pgClient.QueryRow(context.Background(), query, projectId).Scan(&operationer)
	if err != nil {
		log.Println("查询项目运维失败:", err)
		return "蓝胖"
	}
	return operationer
}

func loadConfig(cfgFile string) (*CfgType, error) {
	// 读取配置文件内容
	cfg = &CfgType{}
	_, err := toml.DecodeFile(cfgFile, cfg)
	if err != nil {
		return nil, fmt.Errorf("could not parse config file: %v", err)
	}
	return cfg, nil
}

func mailAlert(job, body string, to []string) {

	subject := fmt.Sprintf("%s Alert: 巡检数据异常!", job)

	err := libs.SendAlertEmail(cfg.Mail, to, subject, body)
	if err != nil {
		log.Printf("Error sending email: %v\n", err)
	}
}
