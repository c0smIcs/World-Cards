package bot

import (
	"strconv"

	tele "gopkg.in/telebot.v4"
)

// startLearningSession - инициализирует и запускает новую сессию обучения для пользователя.
// она может быть вызвана как для нового набора слов, так и для повторения старого
//
// параметры:
//
//	c - контекст telebot, содержащий информацию о сообщении и отправителе.
//	isRepeat - булевый флаг; если true, сессия запускается в режиме "повторения"
func (b *Bot) startLearningSession(c tele.Context, isRepeat bool) error {
	// получаем ID пользователя из контекста
	userID := c.Sender().ID

	// проверяем, есть ли у пользователя слова для изучения
	words, exists := b.UserWords[userID]
	if !exists || len(words) == 0 {
		return c.Send("У вас нет слов для обучения. Добавьте слова с помощью /add")
	}

	// перемешиваем слова пользователя, чтобы они каждый раз были в случайном порядке
	shuffled := shuffleWords(words)

	// создаем новую сессию обучения
	session := LearningSession{
		UserID:       userID,
		Words:        shuffled,
		CurrentIndex: 0,    // начинаем с первого слова (индекс 0)
		Score:        0,    // изначальный счет равен 0
		IsActive:     true, // помечаем сессию как активную
	}

	// сохраняем указатель на созданную сессию в карту активных сессий бота.
	// ключом является ID пользователя
	b.LearningSession[userID] = &session

	if isRepeat {
		c.Send("🔁 Начинаем повторение слов.")
	} else {
		start := strconv.Itoa(len(shuffled))
		c.Send("Начинаем обучение! Слов: " + start)
	}

	// получаем первое слово из перемешанного списка
	getWord := session.Words[0].Original

	// отправляем пользователю первое слово для перевода
	return c.Send(getWord)
}
