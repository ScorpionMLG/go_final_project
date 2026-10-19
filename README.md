# Веб-планировщик задач

## Описание проекта
Веб-сервер, реализованный в рамках учебного задания. Служит для демонстрации навыков/backend-разработки на Go.

## Задания со звездочкой
- Выполнены все

## Запуск проекта локально
Пароль для авторизации: secret123

Через Docker (Windows, GitBash):
```
docker build -t todo-app .

MSYS_NO_PATHCONV=1 docker run -d -p 7540:7540 -v "$(pwd)/scheduler.db:/app/scheduler.db" -e TODO_DBFILE=/app/scheduler.db -e TODO_PASSWORD=secret123 --name todo-container todo-app
```

Через Bash:
```bash
TODO_PASSWORD=secret123 go run main.go
```

**Адрес в браузере:** `http://localhost:7540`

### Настройка tests/settings.go
Перед запуском тестов убедитесь, что в файле `tests/settings.go` установлены следующие параметры:
```
var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYXNzd29yZEhhc2giOiJmY2Y3MzBiNmQ5NTIzNmVjZDNjOWZjMmQ5MmQ3YjZiMmJiMDYxNTE0OTYxYWVjMDQxZDZjN2E3MTkyZjU5MmU0In0.8cobUdHiKexxxKFqMrKekwadlg77Xd2k8taXzYBKLU8`
```