package bot

import (
	tele "gopkg.in/telebot.v4"
)

// HandleHelpCommand - отправляет пользователю справочное сообщение со списком доступных команд
func (b *Bot) HandleHelpCommand(c tele.Context) error {
	helpText := "Доступные команды:\n" +
		"/start - начало работы\n" +
		"/add активация режима\n" +
		"/repeat - повторить предыдущие слова\n" +
		"/go - начать обучение с новыми словами\n" +
		"/alldelete - удалить все слова\n" +
		"/help - помощь\n" +
		"/stop - завершает добавление слов"

	return c.Send(helpText)
}

// HandleUnknownCommand - отправляет сообщение о том, что введенная команда не распознана
func (b *Bot) HandleUnknownCommand(c tele.Context) error {
	return c.Send("Неизвестная команда. Используйте /help для списка команд")
}

// HandleStartCommand - отправляет приветственное сообщение новому пользователю
func (b *Bot) HandleStartCommand(c tele.Context) error {
	return c.Send("Привет! Это бот, который поможет тебе изучить/повторить слова на английском языке. Удачи!")
}

// здесь старый HandleAddCommand
/*
// HandleAddCommand - обрабатывает команду добавления новой пары слов.
// ожидается формат: /add <оригинал> <перевод>
func (b *Bot) HandleAddCommand(c tele.Context) error {
	userID := c.Sender().ID  // получаем ID пользователя
	text := c.Message().Text // получаем текст сообщения

	// убираем команду /add из текста, чтобы остались только аргументы
	text = strings.TrimPrefix(text, "/add")
	// убираем лишние пробелы по краям
	text = strings.TrimSpace(text)

	// проверяем что пользователь ввел что-то кроме команды
	if text == "" {
		return c.Send("Неверный формат. Используйте: /add слово перевод")
	}

	// разбиваем оставшийся текст на части по пробелу
	parts := strings.Split(text, " ")

	// проверяем, что есть минимум 2 части (слово и перевод)
	if len(parts) < 2 {
		return c.Send("Неверный формат. Используйте: /add слово перевод")
	}

	// первая часть - оригинальное слово
	original := parts[0]

	// перевод слова
	translation := parts[1]

	// создаем новое слово
	word := WordPair{
		UserID:     userID,
		Original:   original,
		Translated: translation,
	}

	// сохраняем в память
	b.UserWords[userID] = append(b.UserWords[userID], word)

	response := "Добавлено: " + original + " - " + translation
	return c.Send(response)
}
*/

// HandleGoCommand является оберткой для запуска сессии обучения в обычном режиме
func (b *Bot) HandleGoCommand(c tele.Context) error {
	return b.startLearningSession(c, false)
}

// HandleRepeatCommand является оберткой для запуска сессии обучения в режиме повторения
func (b *Bot) HandleRepeatCommand(c tele.Context) error {
	return b.startLearningSession(c, true)
}

// HandleAlldeleteCommand - удаляет все слова и активные сессии для текущего пользователя
func (b *Bot) HandleAlldeleteCommand(c tele.Context) error {
	userID := c.Sender().ID

	// проверяем, есть ли у пользователя слова, чтобы их удалять
	words, exists := b.UserWords[userID]
	if !exists || len(words) == 0 {
		return c.Send("У вас нет слов для обучения. Добавьте слова с помощью /add")
	}

	// удаляем запись о словах пользователя из хранилища
	delete(b.UserWords, userID)
	// также удаляем активную сессию, если она была
	delete(b.LearningSession, userID)

	return c.Send("Все слова были успешно удалены!")
}

// ------------------------------------------------------------------

// HandleAddCommand активирует режим добавления
func (b *Bot) HandleAddCommand(c tele.Context) error {
	userID := c.Sender().ID

	addWordSession := AddWordSession{
		UserID:   userID,
		IsActive: true,
	}

	b.AddWordSession[userID] = &addWordSession

	c.Send("Вы можете добавлять русские слова. Для завершения введите /stop")

	return nil
}

// HandleStopAddCommand завершает режим добавления
func (b *Bot) HandleStopAddCommand(c tele.Context) error {
	userID := c.Sender().ID

	if session, exists := b.AddWordSession[userID]; exists && session.IsActive {
		delete(b.AddWordSession, userID)
		c.Send("Режим добавления слов завершен")
		return b.startLearningSession(c, false)
	}

	return c.Send("У вас нет активной сессии добавления слов")
}

// isAddSessionActive - метод для проверки активности
func (b *Bot) isAddSessionActive(userID int64) bool {
	session, exists := b.AddWordSession[userID]
	if exists && session.IsActive {
		return true
	}

	return false
}

func (b *Bot) HandleAddWord(c tele.Context) error {
	userID := c.Sender().ID
	text := c.Message().Text

	translation, err := b.Translator.Translate(text, "ru", "en")
	if err != nil {
		return c.Send("Не удалось перевести слово. Попробуйте еще раз")
	}

	word := WordPair{
		UserID:     userID,
		Original:   text,
		Translated: translation,
	}

	b.UserWords[userID] = append(b.UserWords[userID], word)
	return c.Send("Добавлено: " + word.Original + " - " + word.Translated)
}
