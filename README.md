# Media Server

Простой сервер для загрузки, просмотра и удаления файлов. Написан на Go с использованием фреймворка Gin. Поддерживает обычную загрузку файлов и base64-бинарные блобы, фильтрацию по MIME-типам, пагинацию, поиск и регулярные выражения для имён.

## Конечные точки

### `POST /upload`
Загружает файл через multipart/form-data.

**Параметры:**
- `file` — файл (обязательный)
- Любые дополнительные текстовые поля будут включены в ответ

**Ответ:**
```json
{
  "file": "/public/uuid.jpg",
  "field1": "value1"
}
```

---

### `POST /upload-blob`
Загружает файл в формате base64 (возможно, с префиксом Data URL).

**Тело запроса:**
```json
{
  "file": "data:image/jpeg;base64,/9j... | /9j... | base64..."
}
```

**Ответ:**
```json
{
  "file": "/public/uuid.jpg"
}
```

---

### `GET /files`
Возвращает список доступных файлов.

**Параметры:**
- `page` — номер страницы (по умолчанию 1)
- `search` — подстрока для поиска в имени файла (регистронезависимо)

**Ответ:**
```json
{
  "files": ["/public/file1.jpg", "/public/file2.png"],
  "total": 2
}
```

---

### `DELETE /files/:name`
Удаляет файл по имени.

**Ответ при успехе:**
```json
{
  "deleted": true,
  "file": "/public/filename.jpg"
}
```

---

### Статика
Файлы доступны напрямую по пути `/public/*`.

---
## Настройка через переменные окружения

| Переменная         | Влияние и особенности использования |
|--------------------|-------------------------------------|
| `PORT`             | Определяет сетевой порт, на котором запускается HTTP-сервер. Должен быть доступен в контейнере и проброшен в хост через Docker. Не может быть изменён без перезапуска сервиса. |
| `DOMAIN`           | URL-префикс, по которому файлы доступны извне (например, `https://example.com`).|
| `UPLOAD_DIR`       | Физический путь **внутри контейнера**, куда сохраняются загруженные файлы. Должен совпадать с путём в volume (если используется), иначе данные не сохранятся при перезапуске. Рекомендуется использовать абсолютный путь (`/app/public`). |
| `PUBLIC_PATH`      | URL-префикс, по которому файлы доступны извне (например, `/public/file.jpg`). Должен начинаться с `/`. Если изменить — нужно обновить клиентскую логику, ссылающуюся на файлы. Используется для формирования ссылок в ответах. |
| `PAGE_LIMIT`       | Лимитит число файлов, возвращаемых в `/files` за один запрос. Защита от перегрузки памяти и медленной загрузки списка. Увеличение требует больше памяти при большом количестве файлов. |
| `FILE_MAX_SIZE`    | Жёсткий лимит размера файла в байтах. Проверяется **до** чтения MIME и сохранения. Превышение приводит к `400 Bad Request`. Задаётся как число (например, `52428800` = 50 МБ). Следует согласовать с `client_max_body_size`, если используется nginx. |
| `FILE_FILTER`      | Регулярное выражение для фильтрации имён файлов при выводе списка. Например, `\.jpg$` покажет только `.jpg`. Не влияет на загрузку или удаление. Пустое значение блокирует все файлы, `.*` — разрешает все. |
| `ALLOWED_MIMES`    | Список разрешённых MIME-типов через запятую. Проверяется **на основе первых 512 байт файла** (`http.DetectContentType`). Формат: `type/subtype`. Пример: `image/png,video/mp4,application/pdf`. Несуществующие или неподдерживаемые типы игнорируются, файл будет отклонён. Регистрозависим. |


---



## Docker

```bash
docker build -t media-server .
```

**Запуск контейнера:**
```bash
docker run -d \
  -p 8080:8080 \
  -v ./public:/app/public \
  --name media-server \
  media-server
```

---



## Пример docker-compose.yml

```yaml
version: "3.9"

services:
  media-service:
    image: borovlioff/go-media-service:latest
    container_name: media-service
    restart: unless-stopped
    environment:
      UPLOAD_DIR: /app/public
      PUBLIC_PATH: /public
      PAGE_LIMIT: 20
      FILE_MAX_SIZE: 104857600
      FILE_FILTER: .*
      ALLOWED_MIMES: image/jpeg,image/png,video/mp4
      DOMAIN: media-service.local
      PORT: 8080
    volumes:
      - ./public:/app/public
    networks:
      - internal

  nginx:
    image: nginx:stable-alpine
    container_name: media-nginx
    restart: unless-stopped
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./public:/app/public:ro
    networks:
      - internal
      - external

networks:
  internal:
    driver: bridge
  external:
    driver: bridge

```

> **Примечание:** Убедитесь, что директория `./public` существует локально 

## Postman Import
```json
{
  "info": {
    "name": "Media Server API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
    "description": "API для загрузки, получения списка и удаления медиафайлов"
  },
  "item": [
    {
      "name": "Upload File (multipart)",
      "request": {
        "method": "POST",
        "header": [],
        "body": {
          "mode": "formdata",
          "formdata": [
            {
              "key": "file",
              "type": "file",
              "src": []
            }
          ]
        },
        "url": {
          "raw": "{{base_url}}/upload",
          "protocol": "http",
          "host": ["{{base_url}}"],
          "path": ["upload"]
        }
      }
    },
    {
      "name": "Upload Blob (base64)",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"file\": \"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJ...\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/upload-blob",
          "protocol": "http",
          "host": ["{{base_url}}"],
          "path": ["upload-blob"]
        }
      }
    },
    {
      "name": "List Files",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/files?page=1&search=image",
          "protocol": "http",
          "host": ["{{base_url}}"],
          "path": ["files"],
          "query": [
            {
              "key": "page",
              "value": "1"
            },
            {
              "key": "search",
              "value": "image"
            }
          ]
        }
      }
    },
    {
      "name": "Delete File",
      "request": {
        "method": "DELETE",
        "header": [],
        "url": {
          "raw": "{{base_url}}/files/filename.jpg",
          "protocol": "http",
          "host": ["{{base_url}}"],
          "path": ["files", "filename.jpg"]
        }
      }
    }
  ],
  "variable": [
    {
      "id": "base_url",
      "value": "localhost:8080",
      "type": "string",
      "name": "base_url"
    }
  ]
}
```
