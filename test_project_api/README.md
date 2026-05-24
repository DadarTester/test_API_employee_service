# API организационной структуры

REST API для управления подразделениями и сотрудниками с поддержкой древовидной иерархии.

## Технологии

- **Go 1.23**
- **Gorilla Mux** – роутер
- **GORM** – ORM
- **PostgreSQL** – база данных
- **Goose** – миграции
- **Docker / Docker Compose**
- **Logrus** – логирование

## Быстрый старт

### Требования

- Docker Desktop (Windows/Mac) + Docker Compose (Linux)
- Git

### Запуск

1. Клонируйте репозиторий  
   ```bash
   git clone https://github.com/yourusername/org-api.git
   cd org-api

2. Запустите контейнеры  
   ```bash
   docker-compose up --build

Миграции применятся автоматически

API будет доступен на http://localhost:8080

3. Остановка
```bash
   docker-compose down

________________________________________________________________________
________________________________________________________________________
API Эндпоинты

ПОДРАЗДЕЛЕНИЕ

Метод	URL	Описание
POST	/departments	Создать подразделение
GET	/departments/{id}	Получить подразделение (с деревом до depth)
PATCH	/departments/{id}	Переместить/переименовать подразделение
DELETE	/departments/{id}?mode=cascade/reassign	Удалить подразделение

СОЗДАНИЕ подразделения
// POST /departments
{
  "name": "IT",
  "parent_id": null   // или число
}

ПОЛУЧЕНИЕ с деревом
// GET /departments/1?depth=2&include_employees=true

ПЕРЕМЕЩЕНИЕ
// PATCH /departments/1
{
  "parent_id": 2
}


УДАЛЕНИЕ
// DELETE /departments/1?mode=cascade – удалить подразделение, всех сотрудников и всех потомков

// DELETE /departments/1?mode=reassign&reassign_to_department_id=2 – удалить только указанное подразделение, а его сотрудников и дочерние подразделения перевести в подразделение 2


ПРИМЕР:
Создание сотрудника
json
// POST /departments/1/employees
{
  "full_name": "Иван Петров",
  "position": "Разработчик",
  "hired_at": "2025-01-15"   // опционально
}

ТЕСТИРОВАНИЕ:
go test ./tests/...