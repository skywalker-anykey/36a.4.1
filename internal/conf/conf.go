// Package conf - читает файл с конфигурацией
package conf

import (
	"encoding/json"
	"os"
)

// RSSConfig - конфиг полученный из файла с конфигурацией RSS
type RSSConfig struct {
	UrlsRSS       []string `json:"rss"`            // UrlsRSS набор ссылок на RSS-ленты
	RequestPeriod int      `json:"request_period"` // RequestPeriod интервал опроса в минутах
}

// NewRSS - конструктор загружает из файла конфиг и создает объект RSSConfig
func NewRSS(filePath string) (*RSSConfig, error) {
	configFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	rss := RSSConfig{}
	err = json.Unmarshal(configFile, &rss)
	if err != nil {
		return nil, err
	}
	return &rss, nil
}

// BDConfig - конфиг полученный из файла с конфигурацией BD
type BDConfig struct {
	Name     string `json:"name"`     // Name - имя БД
	Port     int    `json:"port"`     // Port - порт БД
	Table    string `json:"table:"`   // Table - имя таблицы
	User     string `json:"user"`     // User - ползователь БД
	Password string `json:"password"` // Password - пароль БД
}

// NewBD - конструктор загружает из файла конфиг и создает объект BDConfig
func NewBD(filePath string) (*BDConfig, error) {
	configFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	bd := BDConfig{}
	err = json.Unmarshal(configFile, &bd)
	if err != nil {
		return nil, err
	}
	return &bd, nil
}
