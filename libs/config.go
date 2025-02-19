// Package libs @Author lanpang
// @Date 2024/8/26 下午2:24:00
// @Desc
package libs

type DB struct {
	Ip       string `toml:"ip"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Sslmode  bool   `toml:"sslmode"`
}

type Mail struct {
	Host     string   `toml:"host"`
	Port     int      `toml:"port"`
	Username string   `toml:"username"`
	Password string   `toml:"password"`
	AddTo    []string `toml:"addto"`
}

func (db DB) HasValue() bool {
	return db.Ip != "" && db.Port != 0 && db.Username != "" && db.Password != ""
}
