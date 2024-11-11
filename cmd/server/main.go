package main

import (
	"fmt"
	"hw36a.4.1/internal/conf"
	"hw36a.4.1/internal/postgres"
	"log"
)

const (
	// pathRSSConfig путь к конфигурации RSS
	pathRSSConfig = "./cmd/server/config.json"
	// pathBDConfig путь к конфигурации BD
	pathBDConfig = "./cmd/server/BD.json"
)

func main() {
	// Загружаем конфигурацию RSS из конфига
	rssConf, err := conf.NewRSS(pathRSSConfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rssConf)

	// Загружаем конфигурацию BD из конфига
	bdConf, err := conf.NewBD(pathBDConfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bdConf)

	// Инициализация БД (таблица пересоздается-обнуляется)
	data, err := postgres.New(bdConf)
	if err != nil {
		log.Fatal(err)
	}
	_ = data

}

// Примерный алгоритм работы
// init BD
// init RSS config
// for rss goroutines - read rss - check cache (if rss not in cache/db - add db)

// init http (api - show news from DB)

/*
// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера.
	var srv server

	// Создаём объекты баз данных.

	// БД в памяти.
	db1 := memDB.New()

	// Реляционная БД Postgres SQL.
	db2, err := postgres.New("postgres://sandbox:sandbox@localhost:5432/news")
	if err != nil {
		log.Fatal(err)
	}

	// Документная БД MongoDB.
	db3, err := mongo.New("mongodb://192.168.1.20:27017/", "news")
	if err != nil {
		log.Fatal(err)
	}
	_, _, _ = db1, db2, db3

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db3

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	_ = http.ListenAndServe(":8080", srv.api.Router())
}
*/
