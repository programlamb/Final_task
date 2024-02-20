# Sprint2_Final_task

## Запуск
1. Скачайте репозиторий.
2. Откройте в среде разработки и запустите __main.go__.
3. Откройте главную страницу http://localhost:8080/
4. Начните работать.

## Использование
Чтобы получить значение выражения, введите его в поле и нажмите __ENTER__.
Чтобы задать время выполнения операций(в секундах) перейдите на страницу __Settings__, задайте параметры и сохраните.

## Примеры
- 2+2*2
- (3-1)*5
- ((4+1)*2)*3
- 100/10
- 1/0  (Выдаст ошибку. На ноль делить нельзя)
- (2+2) * (3+3)
- 10 - 100
- -5 * 5
- 2024 / 1012
- -2/-2

## Как это работает?
Пользователь передаёт выражение на сервер. Там оно обрабатывается оркестратором(сервером) и вычисляется 1-м(нельзя изменить) агентом(работает на самом сервере, __НЕ__ постепенно, а вычисляет всё сразу по истечении общего времени; это __НЕ__ демон) с параметрами, задаными в момент отправки(если их не меняли, каждая операция выполняется 1 секунду).
Пользователь может узнать состояние выражений, обновляя страницу с течением времени.
При перезапуске программы выражения и настройки сбрасываются.
### Схема
пользователь <----> сервер(http://localhost:8080/)
