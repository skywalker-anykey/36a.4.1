package main

import (
	"hw36a.4.1/internal/conf"
	"hw36a.4.1/internal/postgres"
	"hw36a.4.1/internal/rss"
	"log"
	"time"
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
		log.Fatal("ошибка загрузки конфига RSS:", err)
	}

	// Загружаем конфигурацию BD из конфига
	bdConf, err := conf.NewBD(pathBDConfig)
	if err != nil {
		log.Fatal("ошибка загрузки конфига BD:", err)
	}

	// Инициализация БД (таблица пересоздается-обнуляется)
	data, err := postgres.New(bdConf)
	if err != nil {
		log.Fatal("ошибка инициализации БД:,", err)
	}

	// Период опроса серверов RSS
	t := time.Minute * time.Duration(rssConf.RequestPeriod)

	// Запуск пайпов для каждого url (раздельные горутин)
	for _, url := range rssConf.UrlsRSS {
		pipe := cache(readRSS(url, t))
		go reporter(pipe, data)
	}

	// Временное решение для тестов
	// TODO: убрать после теста
	time.Sleep(time.Second * 65)

}

// Читает новости и отправляет в пайп
func readRSS(url string, timer time.Duration) chan rss.Post {
	out := make(chan rss.Post)

	go func() {
		defer close(out)
		for {
			// Получить список новостей если ошибка закрываем пайп
			arPosts, err := rss.GetRSS(url)
			if err != nil {
				log.Println(err)
				return
			}
			// По одной отправляем новости в пайп
			for _, post := range arPosts {
				out <- post
			}

			// После полной отправки ждем период ожидания
			time.Sleep(timer)
		}
	}()
	return out
}

// Пропускает через себя новость только 1 раз
func cache(input <-chan rss.Post) chan rss.Post {
	output := make(chan rss.Post)
	go func() {
		defer close(output)

		// Создаем список уже обработанных новостей
		cacheMap := make(map[string]bool)

		for {
			select {
			case value, ok := <-input:
				// Если канал закрыт, то данных больше не будет и сигнал к завершению работы рутины
				if !ok {
					// Закрываем следующий канал, чтобы оповестить следующую рутину о завершении
					return
				}
				// Если есть ID в cacheMap, то пропускаем, иначе передаем далее
				if cacheMap[value.ID] {
					continue
				} else {
					cacheMap[value.ID] = true
					output <- value
				}
			}
		}
	}()
	return output
}

// Финишная горутина пайпа
func reporter(input <-chan rss.Post, dataBase *postgres.Store) {

	go func() {
		for value := range input {
			// добавляем все новости в БД
			err := dataBase.AddPost(value)
			if err != nil {
				log.Println("ошибка добавления новости в БД:", err)
			}
			log.Printf("Добавлена новость id: '%s' title: '%s'\n", value.ID, value.Title)
		}
	}()
}

// Примерный алгоритм работы

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
