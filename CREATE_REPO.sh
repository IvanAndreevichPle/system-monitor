#!/bin/bash
# Скрипт для создания репозитория на GitHub

echo "=== Создание репозитория system-monitor на GitHub ==="

# Проверка наличия gh CLI
if ! command -v gh &> /dev/null; then
    echo "GitHub CLI (gh) не установлен."
    echo "Установите: https://cli.github.com/"
    echo ""
    echo "Или создайте репозиторий вручную:"
    echo "1. Перейдите на https://github.com/new"
    echo "2. Repository name: system-monitor"
    echo "3. Description: Final project: System monitoring daemon for OTUS Golang course"
    echo "4. Visibility: Public"
    echo "5. НЕ создавайте README, .gitignore или лицензию"
    echo ""
    echo "Затем выполните:"
    echo "  git remote add origin git@github.com:IvanAndreevichPle/system-monitor.git"
    echo "  git push -u origin main"
    exit 1
fi

# Проверка авторизации
if ! gh auth status &> /dev/null; then
    echo "Необходима авторизация в GitHub CLI:"
    echo "  gh auth login"
    exit 1
fi

# Создание репозитория
echo "Создаю публичный репозиторий system-monitor..."
gh repo create system-monitor \
    --public \
    --description "Final project: System monitoring daemon for OTUS Golang course" \
    --source=. \
    --remote=origin \
    --push

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Репозиторий успешно создан!"
    echo "URL: https://github.com/IvanAndreevichPle/system-monitor"
    echo ""
    echo "Проверьте GitHub Actions:"
    echo "  https://github.com/IvanAndreevichPle/system-monitor/actions"
else
    echo ""
    echo "❌ Ошибка при создании репозитория"
    echo "Попробуйте создать вручную через веб-интерфейс GitHub"
fi
