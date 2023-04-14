Репозиторий является частью проекта https://github.com/hablof/logistic-package-api

# logistic-package-api-bot

Сервис является телеграм ботом, gRPC-клиентом для [сервера](https://github.com/hablof/logistic-package-api).

## Telegram-bot
Телеграм бот поддерживает команды _help_ и CRUD команды, команду вывода списка.
Пагинация списка использует `ReplyMarkup` с кнопками вперёд и назад.

## gRPC-client
gRPC код описан в описанном в субмодуле `pkg/logistic-package-api` основного [репозитория](https://github.com/hablof/logistic-package-api).

gRPC-клиент делает запросы к сервису из [репозитория](https://github.com/hablof/logistic-package-api).

## Cache
Сервис использует БД Redis для кеширования запросов. 

Стратегия кеширования:
* метод `Create` — записывает данные в кеш;
* метод `Describe` — пытается прочитать запись из кеша, при неудаче - записывает данные в кеш после ответа gRPC-сервера;
* метод `Remove` — удаляет запись из кеша;
* метод `Update` — удаляет запись из кеша;
* метод `List` — не использует кеш;

## Kafka
Сервис отправляет сообщения в Кафка-топики. 
В топик "omp-tgbot-commands" отправляются сообщения о всех обновлениях из Телеграма.
В топик "omp-tgbot-cache-events" отправляются сообщения о событиях кеша.

## Docker

Описаны докерфайлы для образа бота.

В docker-compose прокидываются конфиги для gRPC-сервера, ретранслятора, бота и фасада.

## Makefile

В мейкфайле описаны команды для локального запуска приложения, сборки докер-образа, и запуска контейнера.