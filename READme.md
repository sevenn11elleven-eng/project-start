# архитектура
	
	API сервис на [[Go]]
	
	база данных [[PostgreSQL]]
	
	кеш [[Redis]]
	
	очередь задач [[RabbitMQ]]
	
	всё запускается через [[Docker]]

# сущности 
User
Product
Order

# api
POST /users
POST /products
POST /orders
GET /orders/{id}

# процесс

## создание проекта

### реализация структуры 
```bash
mkdir -p project-start/cmd/api
mkdir -p project-start/internal/{handler,service,repository}
mkdir -p project-start/migrations
mkdir -p project-start/docker

touch project-start/cmd/api/main.go
touch project-start/docker-compose.yml
```

### настройка docker
для этого проекта мы будем использоватьь redis psql rabbitq
в файле docker-compose.yml задаем для каждого сервиса настройки
любой файл docker состоит из двух параметров это версия и сервисы
- `postgres` — база данных. `volumes` сохраняет данные между перезапусками.
- `redis` — кеш.
- `rabbitmq` — очередь задач с веб-интерфейсом (`http://localhost:15672`).
```json
version: "3.9"
services:
	postgres:
	redis:
	rabbitmq:
volumes: 
	postgres_data:
```
в свою очередь каждый сервис имеет параметры image enviroment ports
итог получаем такой .yml файлик
```json
version: "3.9"
services:
	postgres:
	image: postgres:16
	environment:
		POSTGRES_USER: dev
		POSTGRES_PASSWORD: dev
		POSTGRES_DB: app
		ports:
			- "5432:5432"
	redis:
		image: redis:7
		ports:
			- "6379:6379"
	rabbitmq:
		image: rabbitmq:3-management
		ports:
			- "5672:5672"
			- "15672:15672"
```

### подключение к postgres
```bash
docker exec -it project-start-postgres-1 psql -U dev -d app
```

### http сервер
добавляем 1 простой эндпоинт
```go
func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"status":"ok"}
		json.NewEncoder(w).Encode(response)
	})
	fmt.Println("Server strarted at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
```

---
# to do 



Шаг 6. Кеш.

Когда появится endpoint:

GET /users/{id}

добавь кеширование через Redis.

Логика:

1 запрос → БД  
следующие → Redis

Шаг 7. Очередь.

Когда создаётся order:

POST /orders

ты отправляешь событие в RabbitMQ:

order_created

Отдельный worker читает очередь и, например, имитирует отправку email.

Это покажет **асинхронную архитектуру**.

---

Критически важная вещь:  
не пытайся сделать идеальную архитектуру сразу.

Большинство новичков ломаются на этом этапе.

Твоя цель сейчас:

**система должна работать end-to-end**.

client → API → DB → cache → queue → worker.

---

Если хочешь, дальше я могу показать:

1. **идеальную структуру backend-проекта на Go (которую используют в проде)**
    
2. **полный docker-compose для твоего стека**
    
3. **первый минимальный код API (~100 строк)**
    

Это сильно ускорит тебе старт.
