package main

import (
	"hw36a.4.1/internal/api"
	"hw36a.4.1/internal/conf"
	"hw36a.4.1/internal/postgres"
	"hw36a.4.1/internal/rss"
	"log"
	"net/http"
	"time"
)

const (
	// pathRSSConfig путь к конфигурации RSS
	pathRSSConfig = "./cmd/server/config.json"
	// pathBDConfig путь к конфигурации BD
	pathBDConfig = "./cmd/server/BD.json"
)

type server struct {
	api *api.API
}

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

	// Создаём объект сервера.
	var srv server

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(*data)

	// Запускаем сервер
	_ = http.ListenAndServe(":80", srv.api.Router())

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
