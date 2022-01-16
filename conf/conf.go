package conf

import "os"

type HandlePath struct {
	CreateWager  string
	GetWagerList string
	BuyWager     string
}

type SQLConfig struct {
	DatabaseAddress string
	Username        string
	Password        string
	WagerTable      string
	PurchaseTable   string
}

type Config struct {
	ServerPort int
	Handlers   HandlePath
	SQL        SQLConfig
}

func GetDefaultConfig() *Config {
	return &Config{
		ServerPort: 8080,
		Handlers: HandlePath{
			CreateWager:  "/wagers",
			GetWagerList: "/wagers",
			BuyWager:     "/buy/{wager_id}",
		},
		SQL: SQLConfig{
			DatabaseAddress: "tcp(127.0.0.1:3306)/demo",
			Username:        "root",
			Password:        os.Getenv("MYSQL_ROOT_PASSWORD"),
			WagerTable:      "wagers",
			PurchaseTable:   "purchase",
		},
	}
}

// TODO This function should load configurations from a text file
func LoadConfig() *Config {
	return &Config{}
}
