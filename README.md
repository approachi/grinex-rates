# grinex-rates

gRPC сервис для получения курса USDT с биржи [Grinex](https://grinex.io).

## Требования

- [Docker](https://docs.docker.com/get-docker/) и Docker Compose
- [Go 1.25+](https://golang.org/dl/)
- [Make](https://www.gnu.org/software/make/)

## Запуск

**1. Клонируйте репозиторий:**
```bash
git clone <repo-url>
cd grinex-rates
```

**2. Запустите сервис:**
```bash
make run
```

Команда собирает Docker-образ и поднимает приложение вместе с PostgreSQL, Jaeger, Prometheus и Grafana. Миграции применяются автоматически.

**3. Проверьте, что сервис работает:**
```bash
grpcurl -plaintext \
  -d '{"method":"topN","position_from":1}' \
  localhost:50051 rates.RatesService/GetRates
```

Ожидаемый ответ:
```json
{
  "ask": "...",
  "bid": "...",
  "timestamp": "..."
}
```

**4. Остановка:**
```bash
make stop
```

## API

### GetRates

Получает текущий курс USDT с биржи и сохраняет его в БД.

```
rpc GetRates(GetRatesRequest) returns (GetRatesResponse)
```

**Параметры запроса:**

| Поле | Тип | Описание |
|---|---|---|
| `method` | string | Метод расчёта: `topN` или `avgNM` |
| `position_from` | int32 | Позиция N (единственная для `topN`, начало диапазона для `avgNM`) |
| `position_to` | int32 | Позиция M (конец диапазона, только для `avgNM`) |

**Ответ:**

| Поле | Тип | Описание |
|---|---|---|
| `ask` | int64 | Цена ask в минорных единицах |
| `bid` | int64 | Цена bid в минорных единицах |
| `timestamp` | Timestamp | Время получения курса |

**Примеры:**

```bash
# Курс из первой позиции стакана
grpcurl -plaintext \
  -d '{"method":"topN","position_from":1}' \
  localhost:50051 rates.RatesService/GetRates

# Среднее значение позиций 1–5
grpcurl -plaintext \
  -d '{"method":"avgNM","position_from":1,"position_to":5}' \
  localhost:50051 rates.RatesService/GetRates
```

### HealthCheck

```
rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse)
```

Возвращает `{ "ok": true }` если сервис работает.

```bash
grpcurl -plaintext localhost:50051 rates.RatesService/HealthCheck
```

## Тестирование

Юнит-тесты:
```bash
make test
```

Все тесты включая интеграционные (автоматически поднимают отдельную тестовую БД):
```bash
make ftest
```

## Make-команды

| Команда | Описание |
|---|---|
| `make build` | Собрать бинарник в `bin/server` |
| `make run` | Собрать и запустить все сервисы |
| `make stop` | Остановить все сервисы |
| `make test` | Запустить юнит-тесты |
| `make test-unit` | Запустить только юнит-тесты |
| `make test-integration` | Запустить только интеграционные тесты |
| `make ftest` | Запустить юнит + интеграционные тесты |
| `make docker-build` | Собрать Docker-образ |
| `make lint` | Запустить линтер (через Docker) |
| `make db-conn` | Подключиться к БД через psql |
| `make proto` | Перегенерировать код из .proto файлов |

## Конфигурация

Конфигурация загружается из `config.yaml` в рабочей директории. Значения можно переопределить через переменные окружения (с префиксом `APP_`) или флаги запуска — флаг имеет наивысший приоритет.

### config.yaml

```yaml
grpc:
  port: 50051

exchange:
  timeoutseconds: 10      # таймаут HTTP-запроса к бирже в секундах

postgres:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  database: grinex
  sslmode: disable        # disable | require | verify-full

metrics:
  port: 8080

tracing:
  endpoint: localhost:4318

log:
  level: info             # debug | info | warn | error
```

Все поля обязательны — при отсутствии любого значения приложение завершится с ошибкой.

### Переменные окружения

Каждый параметр конфига переопределяется через env var с префиксом `APP_` и `_` вместо `.`:

```bash
APP_POSTGRES_HOST=db APP_POSTGRES_PASSWORD=secret ./bin/server
```

### Флаги запуска

```bash
./bin/server \
  --postgres-host=db \
  --postgres-port=5432 \
  --postgres-user=postgres \
  --postgres-password=secret \
  --postgres-database=grinex \
  --postgres-sslmode=disable \
  --exchange-timeout=10 \
  --grpc-port=50051 \
  --metrics-port=8080 \
  --tracing-endpoint=localhost:4318 \
  --log-level=info
```

## Observability

| Компонент | Адрес |
|---|---|
| Prometheus метрики | http://localhost:8080/metrics |
| Grafana | http://localhost:3000 (admin / admin) |
| Jaeger UI | http://localhost:16686 |

---

Выбор технологий и рекомендации для production описаны в [docs/COVER_LETTER.md](docs/COVER_LETTER.md).
