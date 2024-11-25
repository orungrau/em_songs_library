# Song Library for Effective Mobile

## Запуск приложения

1. **Настройка конфигурации**  
   Переименуйте файл `default.env` в `.env` и укажите необходимые настройки в этом файле.

2. **Поднятие базы данных**  
   Для запуска базы данных PostgreSQL используйте команду:
   ```bash
   docker-compose -f docker-compose.dev.yml up -d
   ```

3. **Запуск приложения**  
   Выполните следующую команду:
   ```bash
   go run cmd/song_library.go
   ```

4. **Доступ к Swagger**  
   После запуска приложения Swagger будет доступен по адресу:  
   [http://{host:port}/swagger/index.html](http://{host:port}/swagger/index.html)