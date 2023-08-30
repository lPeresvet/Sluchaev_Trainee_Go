## Инструкция по запуску
1) Собрать образ приложения. Для этого надо выполнить команду 
`docker build --no-cache --pull -t trainee .`
2) Запустить приложение и базу данных в докере. Нужно перейти в директорию __compose__ и выполнить команду `docker compose up`
3) После того как в терминале появится панель с информацией о Fiber, приложение готово к работе!

## Тестирование