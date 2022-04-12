package builder

import (
	"fmt"
	"stage-sync/config"
)

// BuildConnectionString builds a postgres connection string in url form from the config
func BuildConnectionString(credentials config.ConfigDB) string {
	userPart := ""
	if credentials.User != "" {
		userPart = fmt.Sprintf("%s", credentials.User)
	}
	if credentials.Password != "" {
		userPart = fmt.Sprintf("%s:%s", userPart, credentials.Password)
	}
	if userPart != "" {
		userPart = fmt.Sprintf("%s@", userPart)
	}
	cs := fmt.Sprintf("postgres://%s%s:%d/%s", userPart, credentials.Host, credentials.Port, credentials.Database)
	if credentials.SslMode != "" {
		cs = fmt.Sprintf("%s?sslmode=%s", cs, credentials.SslMode)
	}
	return cs
}
