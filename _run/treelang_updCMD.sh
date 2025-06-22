#!/bin/bash
set -e

# Перевіряємо, чи встановлено jq
if ! command -v jq &> /dev/null
then
    echo "jq не встановлено. Будь ласка, встановіть jq і спробуйте ще раз."
    exit 1
fi

# Директорія, у якій виконується скрипт
run_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# API-ендпоінт GitHub для отримання останнього релізу
API_URL="https://api.github.com/repos/Bookshelf-Writer/treelang/releases/latest"

# Назва файлу-бінару, який треба завантажити
FILE="treelang.linux.bin"

#############################################################################

# Отримуємо JSON-інформацію про останній реліз
RELEASE_INFO=$(curl -s "$API_URL")

# Витягаємо посилання на потрібний файл
DOWNLOAD_URL=$(echo "$RELEASE_INFO" | jq -r ".assets[] | select(.name == \"$FILE\") | .browser_download_url")

# Якщо посилання не знайдено — виводимо повідомлення та завершуємо
if [ -z "$DOWNLOAD_URL" ]; then
    echo "Файл $FILE не знайдено у останньому релізі."
    exit 1
fi

# Завантажуємо файл у поточну директорію та надаємо йому права на виконання
curl -L "$DOWNLOAD_URL" -o "$run_dir/$FILE"
chmod +x "$run_dir/$FILE"

echo "Файл $FILE успішно завантажено"
