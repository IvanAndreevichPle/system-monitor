# Настройка проекта "Системный мониторинг"

## Понимание требований:

✅ **Отдельный репозиторий в GitHub**
- Создать новый репозиторий: `system-monitor` или `final-project-system-monitor`
- URL: `git@github.com:IvanAndreevichPle/system-monitor.git`

✅ **Работа с PR по этапам**
- Каждый этап = отдельная ветка
- Каждый этап = отдельный PR в `main`
- После мерджа PR, переходим к следующему этапу

✅ **GitHub Actions**
- Настроить CI/CD пайплайн
- Запуск golangci-lint
- Запуск юнит-тестов: `go test -race -count 100`
- Сборка бинаря для Go >= 1.22

## Структура работы:

### Этап 1: Protobuf схема + gRPC сервер
- Ветка: `stage1-protobuf-grpc`
- PR: `feat: add protobuf schema and basic gRPC server`
- После мерджа → переходим к Этапу 2

### Этап 2: Load Average и CPU
- Ветка: `stage2-load-cpu`
- PR: `feat: implement load average and CPU metrics collection`
- После мерджа → переходим к Этапу 3

### Этап 3: Дисковые метрики
- Ветка: `stage3-disk`
- PR: `feat: implement disk metrics collection`

### Этап 4: Сетевые метрики
- Ветка: `stage4-network`
- PR: `feat: implement network metrics collection`

### Этап 5: Усреднение
- Ветка: `stage5-averaging`
- PR: `feat: implement metrics averaging over time window`

### Этап 6: Конфигурация
- Ветка: `stage6-config`
- PR: `feat: add configuration support (CLI + file)`

### Этап 7: Тесты + клиент
- Ветка: `stage7-tests-client`
- PR: `feat: add tests and simple client`

### Этап 8: Dockerfile + CI/CD
- Ветка: `stage8-docker-cicd`
- PR: `feat: add Dockerfile, Makefile and CI/CD pipeline`

## План действий:

1. **Создать новый репозиторий в GitHub**
   - Имя: `system-monitor`
   - Описание: "Final project: System monitoring daemon for OTUS Golang course"
   - Публичный или приватный?

2. **Инициализировать локальный репозиторий**
   ```bash
   cd system-monitor
   git init
   git remote add origin git@github.com:IvanAndreevichPle/system-monitor.git
   ```

3. **Создать начальную структуру**
   - README.md
   - .gitignore
   - go.mod
   - Базовая структура проекта

4. **Настроить GitHub Actions**
   - .github/workflows/ci.yml
   - Линтер, тесты, сборка

5. **Начать работу по этапам**
