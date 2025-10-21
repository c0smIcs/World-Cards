package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/c0smIcs/Word-Cards/internal/bot"
	"github.com/c0smIcs/Word-Cards/internal/config"

	tele "gopkg.in/telebot.v4"
)

func main() {
	// загружаем конфигурацию для переменных окружения
	cfg, err := config.Lead()
	if err != nil {
		// в случае ошибки при загрузке конфигурации, приложение не может продолжить работу
		log.Fatal(err)
	}

	// инициализируем глобальный источник случайных чисел
	// это необходимо для корректной работы функции перемешивания слов
	rand.NewSource(time.Now().UnixNano())

	// настраиваем параметры для клиента telebot
	pref := tele.Settings{
		Token: cfg.TgToken, // указываем токен бота
		Poller: &tele.LongPoller{ // используем Long Polling для получения обновлений от tg
			Timeout: 10 * time.Second, // таймаут ожидания обновлении
		},
	}

	// создаем новый экземпляр бота telebot с заданными настройками
	telebot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("Ошибка, при создании экземпляра бота: %v",err)
		return
	}

	// инициализируем in-memory хранилище для данных пользователей.
	// при перезапуске бота эти данные будут утеряны
	userWords := make(map[int64][]bot.WordPair)
	learningSession := make(map[int64]*bot.LearningSession)
	addWordSession := make(map[int64]*bot.AddWordSession)
	translator := bot.NewLibreTranslator("http://localhost:5000", "")
	
	// создаем наш собственный объект Bot, передавая ему клиент telebot и хранилища
	// это позволяет инкапсулировать логику работы с данными внутри нашей структуры
	appBot := bot.NewBot(telebot, userWords, learningSession, addWordSession, translator)

	// регистрируем главный обработчик для всех текстовых сообщений
	// вся логика по маршрутизации сообщений (команды, ответы) будут внутри HandleTextMessage
	telebot.Handle(tele.OnText, appBot.HandleTextMessage)

	fmt.Println("Бот запущен")

	// запускаем бота. этот вызов блокирует выполнение программы
	// пока бот не будет остановлен
	telebot.Start()
}
