package bot

import (
	tele "gopkg.in/telebot.v4"
)

// Bot - это основная структура, которая инкапсулирует состояние и зависимости бота.
// она содержит экземпляр клиента telebot и хранилище данных пользователей
type Bot struct {
	Bot             *tele.Bot                  // экземпляр бота telebot для взаимодействия с Telegram API
	UserWords       map[int64][]WordPair       // Хранилище слов пользователей. Ключ - ID пользователя, значение - срез его слов
	LearningSession map[int64]*LearningSession // Хранилище активных сессии обучения. Ключ - ID пользователя, значение - указатель на его сессию
	AddWordSession  map[int64]*AddWordSession  // Хранилище активных сессии для добавления слов
	Translator      Translator
}

// NewBot - является конструктором для структуры Bot
// он инициализирует и возвращает новый экземпляр бота с предоставленными зависимостями
func NewBot(bot *tele.Bot, userWords map[int64][]WordPair, learningSession map[int64]*LearningSession, addWordSession map[int64]*AddWordSession, translator Translator) *Bot {
	return &Bot{
		Bot:             bot,
		UserWords:       userWords,
		LearningSession: learningSession,
		AddWordSession:  addWordSession,
		Translator:      translator,
	}
}
