## Подтверждает запись на подачу паспорта за вас.

Чтобы записаться на подачу российского загран паспорта в Тбилиси нужно каждые 24 часа подтверждать свою заявку на сайте https://q.midpass.ru/, иначе она понижается в очереди.

Этот код подтверждает заявку самостоятельно. 

## Важно
Код предоставляю as-is, так как мое подтверждение на подачу паспорта пришло раньше, чем я его привел в нормальное состояние. 

В частности не было проверено, как оно норм работает в контейнере. Также нет обработки ошибок, когда капча введена неверно.

Для того, чтобы код работал, надо:
1. Прописать свои учетные данные в .env файле
2. Зарегаться на https://2captcha.com/ и вставить свой api key
3. Выполнить `docker compose up -d`

Как подебажить: 
1. Изменить `chromedp.Flag("headless", true)` на `chromedp.Flag("headless", false)`
2. локально запустить go run ./main.go (запуститься, если локально у вас есть хром)