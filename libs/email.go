// Package libs @Author lanpang
// @Date 2025/1/16 上午11:20:00
// @Desc
package libs

import (
	"gopkg.in/gomail.v2"
	"log"
)

func SendAlertEmail(cfg Mail, to []string, subject, body string) error {
	// 创建邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.Username) // 发件人
	m.SetHeader("To", to...)          // 收件人
	m.SetHeader("Subject", subject)   // 邮件主题
	m.SetBody("text/html", body)      // 邮件正文，支持 "text/html" 或 "text/plain"

	// 使用 SMTP 发送邮件
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	log.Println("Alert email sent successfully!")
	return nil
}
