# Описание задачи

Реализовать сервис, который будет получать по апи ФИО, из открытых апи обогащать ответ наиболее вероятными возрастом, полом и национальностью и сохранять данные в БД. По запросу выдавать инфу о найденных людях. Необходимо реализовать следующее:

1. Выставить rest методы
    - Для получения данных с различными фильтрами и пагинацией
    - Для удаления по идентификатору
    - Для изменения сущности
    - Для добавления новых людей в формате

```
{
  "name": "Dmitriy",
  "surname": "Ushakov",
  "patronymic": "Vasilevich" // необязательно
}
```

2. Корректное сообщение обогатить
- Возрастом - https://api.agify.io/?name=Dmitriy
- Полом - https://api.genderize.io/?name=Dmitriy
- Национальностью - https://api.nationalize.io/?name=Dmitriy

3. Обогащенное сообщение положить в БД postgres (структура БД должна быть создана путем миграций)
4. Покрыть код debug- и info-логами
5. Вынести конфигурационные данные в .env


# Описание структуры сервиса
На верхнем уровне реализован объект App, который инкапсулирует через интерфейсы взаимодействие с базой данных, сервисами для обогащения данных.

Прочие сервисы вынесены в отдельные пакеты и удовлетворяют необходимым интерфейсам из пакета app.


# TODO
- [x] Базу данных SQLite заменить на Postgres. Это потребует изменения метода Migration: необходимо изменить типы данных в таблице.
- [ ] Метод Get сейчас не осуществляет поиска в базе данных, а возвращает все объекты из нее. Нужно подумать над тем, как реализовать поиск по запросу пользователя.
- [ ] Тестирование реализовано верхнеуровнево, нужно покрыть код тестами.
- [x] В настоящий момент логируются только ошибки, реализовать логгирование отладочной информации.
- [x] Вынести конфигурационные данные в .env


# Идеи для расширения сервиса
- [ ] Обогащение возраста работает только в том случае, когда имя указано латиницей. Можно подумать над транслитерацией имени. Это позволит обогащать данные возрастом для имен, написанных кирилицей.
- [ ] Необходимо реализовать кэширование данных, чтобы не собирать вновь данные, которые уже получены от сторонних сервисов. Нужно фиксировать период актуальности, в течение которого не делать запросы к API, а забирать данные из кэша.