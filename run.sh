echo "🚀 Запуск проекта Word-Cards..."

# Проверяем, установлен ли Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker не установлен. Пожалуйста, установите Docker сначала."
    exit 1
fi

# Останавливаем старые контейнеры
echo "🧹 Останавливаем старые контейнеры..."
docker-compose down

# Запускаем LibreTranslate в фоне
echo "📦 Запускаем LibreTranslate..."
docker-compose up -d libretranslate

# Ждем запуска сервиса перевода
echo "⏳ Ждем запуска сервиса перевода..."
sleep 10

# Проверяем работу перевода
echo "🔍 Проверяем работу перевода..."
for i in {1..10}; do
    response=$(curl -s -X POST "http://localhost:5000/translate" \
      -H "Content-Type: application/json" \
      -d '{"q":"привет", "source":"ru", "target":"en", "format":"text"}' || echo "error")
    
    if [[ $response == *"hello"* ]]; then
        echo "✅ Переводчик готов!"
        break
    else
        echo "⏳ Переводчик еще не готов. Попытка $i/10..."
        sleep 5
    fi
    
    if [ $i -eq 10 ]; then
        echo "❌ Переводчик не запустился за отведенное время"
        docker-compose logs libretranslate
        exit 1
    fi
done

echo "🤖 Запускаем бота..."
go run ./cmd/bot