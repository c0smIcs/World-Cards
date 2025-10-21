.PHONY: help start stop status logs clean

help:
	@echo "Доступные команды:"
	@echo "  make start  - запустить LibreTranslate"
	@echo "  make stop   - остановить LibreTranslate"
	@echo "  make status - показать статус контейнера"
	@echo "  make logs   - показать логи LibreTranslate"
	@echo "  make clean  - остановить и удалить контейнер"

start:
	docker-compose -d libretranslate
	@echo "LibreTranslate запущен на http://loclhost:5000"
	@echo "Теперь можно запустить бота: go run ./cmd/bot"

stop:
	docker-compose stop

status:
	docker-compose ps

logs:
	docker-compose logs -f libretranslate

clean:
	docker-compose down



# .SILENT:

# build:
# 	go build -o ./.bin/bot cmd/bot/main.go

# run: build
# 	./.bin/bot