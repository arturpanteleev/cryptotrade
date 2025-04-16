#!/bin/bash

# Путь к папке с проектом
PROJECT_DIR="/var/cryptotrade"

# Имя сессии screen (можно выбрать любое имя для сессии)
SCREEN_SESSION="cryptotrade-session"

# Заходим в папку проекта
cd $PROJECT_DIR

# Останавливаем сервер (отключаем screen сессию)
echo "Останавливаем сервер..."
screen -S $SCREEN_SESSION -X quit

# Обновляем код
echo "Получаем последние изменения из репозитория..."
git pull origin master

# Собираем новый бинарник
echo "Собираем новый бинарник..."
go build -o cryptotrade

# Запускаем сервер в новой screen сессии
echo "Запускаем сервер..."
screen -S $SCREEN_SESSION -d -m ./cryptotrade

# Лог
echo "Процесс завершен!"
