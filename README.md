# device-parser-logs

Сервис для обработки TSV файлов, сохранения данных в MongoDB и создания отчетов в PDF формате

## Функциональность
- Мониторинг директории с TSV файлами
- Асинхронная обработка через очередь сообщений
- Сохранение данных в MongoDB
- Генерация отчетов в формате PDF по уникальному идентификатору устройства
- Запись ошибок парсинга в базу данных и файл
- HTTP-сервер для получения данных по устройствам с пагинацией
## Технологии
- Go 1.25
- MongoDB - база данных
- Docker - контейнеризация
- Chi - HTTP роутер
- testify - тестирование
- log/slog - логирование
- gofpf - генерация pdf файлов
## Установка и запуск
### 1. Клонирование репозитория  
```bash
git clone https://github.com/EBichuk/device-parser-logs.git
cd device-parser-logs
```
### 2. Запуск приложения
```bash
make run
```
Эта команда запустит docker контейнеры:
- MongoDB на порту 27017
- Приложение на порту 8080
### 3. Остановка приложения
```bash
make dowm
```

## Структура
```
device-parser-logs/
├── cmd/
│   └── main.go                     # Точка входа приложения              
├── device-reports/                 # Директория с pdf файлами
├── device-storage/                 # Директория с tsv файлами
├── internal/
│   ├── api/                        # Запуск и остановка приложения
│   │   └── api.go                       
│   ├── config/                     # Конфигурация
│   │   └── config.go
│   ├── controller/                 # Слой controller              
|   |   ├── handlers.go             # Handlers
|   |   └── router.go               # HTTP сервер
│   ├── generator/                  # Создание pdf файлов
│   │   └── generator.go
│   ├── models/                     # Модели данных
│   │   └── model.go
│   ├── parser/                     # Обработка tsv файлов
│   │   ├── parser.go
|   |   └── parser_test.go
│   ├── repository/                 # Слой repository
│   │   └── repository.go
│   ├── service/                    # Слой service
│   |   └── service.go
│   └── watcher/                    # Сканирование директории, очередь задач,
│       └── wacher.go               # параллельная обработка файлов воркерами
├── prg/
│   ├── client/                     # Соединение с MongoDB                                          
│   │   └── mongodb.go                       
│   ├── errorx/                     # Кастомные ошибки
│   │   └── errorx.go               
│   └── font/                       # Шрифты для генерации pdf файлов   
├── producer/                       # Продюсер tsv файлов 
│   └── producer.go
├── .golangci.yml                   # Настройки линтера                         
├── .env.example                    # Пример переменных окружения
├── go.mod                          
├── go.sum                          
├── docker-compose.yml              
├── Dockerfile
└── makefile                
```

## API endpoints
- **GET /devices/{guid}?page=&?limit=**

Получение данных по unit_guid с пагинацией

```bash
curl "http://localhost:8080/devices/01749246-95f6-57db-b7c3-2ae0e8be671f?page=1&?limit=10"
```
Ответ:
```
{
    "data": [
        {
            "id": "6999df418082bdf93e17968d",
            "mqtt": "",
            "invid": "G-044322",
            "guid": "01749246-95f6-57db-b7c3-2ae0e8be671f",
            "msg_id": "cold7_Defrost_status",
            "text": "Разморозка",
            "context": "",
            "class_msg": "waiting",
            "level": 100,
            "area": "LOCAL",
            "addr": "cold7_status.Defrost_status",
            "block": "",
            "type": "",
            "bit": "",
            "invert_bit": "",
            "create_at": "2026-02-21T16:37:21.31Z"
        },
        ...
    ],
    "total": 4,
    "page": 1,
    "limit": 10
}
```
