package bot

import (
	"fmt"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
)

// HandleCommand - является маршрутизатором для команд, начинающихся с '/'.
// он анализирует текст сообщения, извлекает команду и вызывает соответствующий обработчик
func (b *Bot) HandleCommand(c tele.Context) error {
	text := c.Message().Text

	// разделяем слова для анализа
	parts := strings.Fields(text)
	// проверяем, что сообщение не пустое и начинается с '/', иначе это не команда
	if len(parts) == 0 || !strings.HasPrefix(parts[0], "/") {
		return c.Send("Необходимо ввести корректную команду с '/'")
	}

	// извлекаем чистую команду, убирая '/'
	command := strings.TrimPrefix(parts[0], "/")

	// используем switch для вызова нужного метода в зависимости от команды.
	switch command {
	case "start":
		return b.HandleStartCommand(c)
	case "add":
		return b.HandleAddCommand(c)
	case "go":
		return b.HandleGoCommand(c)
	case "help":
		return b.HandleHelpCommand(c)
	case "repeat":
		return b.HandleRepeatCommand(c)
	case "alldelete":
		return b.HandleAlldeleteCommand(c)
	case "stop":
		return b.HandleStopAddCommand(c)
	default:
		return b.HandleHelpCommand(c)
	}
}

// HandleAnswer - обрабатывает ответ пользователя во время сессии обучения
func (b *Bot) HandleAnswer(c tele.Context) error {
	userID := c.Sender().ID
	text := c.Message().Text
	session := b.LearningSession[userID]
	
	// получаем текущее слово, на которое отвечает пользователь
	currentWord := session.Words[session.CurrentIndex]

	// очищаем ответ пользователя и правильный ответ от лишних пробелов
	userAnswer := strings.TrimSpace(text)
	correctAnswer := strings.TrimSpace(currentWord.Translated)

	// сравниваем ответы без учета регистра
	if strings.EqualFold(userAnswer, correctAnswer) {
		c.Send("✅ Правильный перевод!")
		session.Score++ // увеличиваем счет при правильном ответе
	} else {
		c.Send("❌ Неправильно. Правильный перевод: " + correctAnswer)
	}

	// переходим к следующему слову
	session.CurrentIndex++

	// проверяем, закончились ли слова в сессии
	if session.CurrentIndex >= len(session.Words) {
		// формируем итоговое сообщение
		result := fmt.Sprintf("Обучение завершено! Результат: %s/%s\n\nЕсли хотите повторить эти же слова, используйте /repeat\nИли удалите все слова: /alldelete, чтобы начать с чистого листа!",
		strconv.Itoa(session.Score),
		strconv.Itoa(len(session.Words)))
		
		c.Send(result)
		// удаляем сессию обучения, так как она завершена
		delete(b.LearningSession, userID)
	} else {
		// если слова еще есть, отправляем следующее
		nextWord := session.Words[session.CurrentIndex].Original
		c.Send(nextWord)
	}

	return nil
}

// HandleTextMessage - это главный обработчик для всех входящих текстовых сообщений.
// он определяет, является ли сообщение командой или ответом в рамках сессии обучения
func (b *Bot) HandleTextMessage(c tele.Context) error {
	userID := c.Sender().ID
	text := c.Message().Text

	// проверяем, существует ли для этого пользователя активная сессия обучения
	session, exists := b.LearningSession[userID]
	// проверяем, существует ли для этого пользователя активная сессия добавления слов
	sessionAdd, existsAdd := b.AddWordSession[userID]
	
	// если сессия существует, она активная, и сообщение не является командой,
	// то расцениваем его как ответ на вопрос
	if exists && session.IsActive && !strings.HasPrefix(text, "/") {
		return b.HandleAnswer(c)
	} else if existsAdd && sessionAdd.IsActive && !strings.HasPrefix(text, "/") {
		return b.HandleAddWord(c)
	} else {
		// в противном случае, обрабатываем сообщение как новую команду
		return b.HandleCommand(c)
	}
}