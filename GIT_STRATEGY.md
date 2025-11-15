# Стратегия работы с Git для проекта "Системный мониторинг"

## Подход: Этапная разработка с PR

Будем работать по этапам, каждый этап = отдельный коммит/PR.

## Этапы и их коммиты:

### Этап 1: Protobuf схема + gRPC сервер (базовая структура)
- Создать структуру проекта
- Определить Protobuf схему
- Реализовать базовый gRPC сервер
- **Коммит**: `feat: add protobuf schema and basic gRPC server`

### Этап 2: Сбор Load Average и CPU (простые метрики)
- Реализовать сбор Load Average
- Реализовать сбор CPU метрик
- **Коммит**: `feat: implement load average and CPU metrics collection`

### Этап 3: Сбор дисковых метрик
- Реализовать сбор дисковых метрик
- **Коммит**: `feat: implement disk metrics collection`

### Этап 4: Сбор сетевых метрик
- Реализовать сбор сетевых метрик
- **Коммит**: `feat: implement network metrics collection`

### Этап 5: Усреднение за период M секунд
- Реализовать алгоритм усреднения
- **Коммит**: `feat: implement metrics averaging over time window`

### Этап 6: Конфигурация
- Реализовать конфигурацию через CLI и файл
- **Коммит**: `feat: add configuration support (CLI + file)`

### Этап 7: Тесты + клиент
- Написать юнит-тесты
- Написать интеграционные тесты
- Реализовать простой клиент
- **Коммит**: `feat: add tests and simple client`

### Этап 8: Dockerfile + CI/CD
- Создать Dockerfile
- Создать Makefile
- Настроить CI/CD пайплайн
- **Коммит**: `feat: add Dockerfile, Makefile and CI/CD pipeline`

## Работа с ветками:

**Вариант 1: Одна ветка с этапными коммитами**
- Создать ветку `final_project_system_monitor`
- Делать коммиты по этапам
- В конце один PR со всеми этапами

**Вариант 2: Отдельные ветки для каждого этапа**
- `final_project_system_monitor_stage1`
- `final_project_system_monitor_stage2`
- и т.д.
- Отдельный PR для каждого этапа

**Рекомендация**: Вариант 1 (одна ветка) - проще для финального проекта.

## Структура коммитов:

```
feat: add protobuf schema and basic gRPC server

- Define protobuf schema for system metrics
- Implement basic gRPC server with streaming
- Add project structure (cmd/, internal/, api/)
```

