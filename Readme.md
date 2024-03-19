# Task Manager REST API

Этот репозиторий содержит REST API сервис Task Manager, который позволяет пользователям управлять заданиями и отслеживать их выполнение. API включает в себя методы для создания пользователей, создания квестов (задач), записи выполнения задачи для пользователей, получения баланса пользователей и списка выполненных задач.

## Конечные Точки

- `/user/`: POST метод для создания нового пользователя. Тело запроса должно содержать name (обязательно) и balance (опционально) в формате JSON.
- `/quest`: POST метод для создания нового квеста.  Тело запроса должно содержать name и cost в формате JSON.
- `/user/{userID}/quests/{questID}`: POST метод для записи выполнения задания для пользователя.
- `/user/{userID}/history`: GET метод для получения баланса пользователя и списка выполненных задач.

## Запуск Сервиса

### Использование Docker

1. Готовый Docker образ: Вы можете загрузить готовый Docker образ с Docker Hub.
```bash
    docker pull yasminworks/taskmanager
```

2. Создание Docker образа и запуск контейнера:
```bash
   docker volume create task-manager
   docker build . -t taskmanager:latest
   docker run -d -it -p 8082:8082 -v task-manager:/app/storage taskmanager
``` 
   
### Компиляция Исходного Кодаt
0. Перейдите в корневую папку проекта

1. Установка зависимостей:
```bash
  go mod download
 ```

2. Подготовка базы данных:
```bash
   go run ./cmd/migrator --storage-path=./storage/storage.db --migrations-path=./migrations
 ```  
3. Компиляция и запуск:
 ```bash 
    go build -o task-manager ./cmd/task-manager/main.go
    ./task-manager
 ```  
### Использование Утилиты Task

Если у вас установлена утилита Task, можно запустить сервис командой
```bash
    task build
```
