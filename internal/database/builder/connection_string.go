package builder

import (
	"fmt"
	"stage-sync-cli/config"
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
	fmt.Println(userPart)
	if userPart != "" {
		userPart = fmt.Sprintf("%s@", userPart)
	}
	return fmt.Sprintf("postgres://%s%s:%d/%s?sslmode=require", userPart, credentials.Host, credentials.Port, credentials.Database)
}
