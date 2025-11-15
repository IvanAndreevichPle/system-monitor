# Инструкция по созданию репозитория на GitHub

## Шаги:

1. **Создайте новый репозиторий на GitHub:**
   - Перейдите на https://github.com/new
   - Repository name: `system-monitor`
   - Description: "Final project: System monitoring daemon for OTUS Golang course"
   - Visibility: **Public** ✅
   - НЕ создавайте README, .gitignore или лицензию (уже есть локально)

2. **Подключите локальный репозиторий:**
   ```bash
   cd /home/drplekhanov/OTUS/Golang/otus_hw/system-monitor
   git remote add origin git@github.com:IvanAndreevichPle/system-monitor.git
   git push -u origin main
   ```

3. **Проверьте GitHub Actions:**
   - После пуша перейдите в раздел "Actions" на GitHub
   - Должен запуститься CI пайплайн

## Готово! 🎉

После этого можно начинать работу по этапам:
- Этап 1: `git checkout -b stage1-protobuf-grpc`
