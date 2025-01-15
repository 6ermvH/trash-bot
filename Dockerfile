# Используем минимальный образ Python
FROM python:3.10-slim

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем файл с зависимостями (если есть requirements.txt)
COPY requirements.txt ./

# Устанавливаем зависимости (указанные в requirements.txt)
RUN pip install --no-cache-dir -r requirements.txt

# Копируем файлы вашего проекта в контейнер
COPY . .

# Определяем команду запуска контейнера
CMD ["python3", "bot.py"]

