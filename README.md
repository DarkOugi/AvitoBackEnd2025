# AvitoBackend-trainee-assignment-winter-2025  

Тестовое на стажировку Golang

# Как запустить проект: 
`docker compose up --build`  
# Что еще нужно сделать в мечтах 
1 - подумать как не допустить ситуацию создания миллионов jwt токкенов  
2 - e2e тесты

# О проекте:
## важное о Аунтификации (регистрации/входе)  
Выдвинутых требований к логину со сторону ТЗ не было  
Но и хранить любой логин мы не должны, как минимум это странно  
Логин приведен в строгий формат second_name.first_name  
(если у нас есть уже сотрудник с данным именем то добавляются цифры,  
сейчас это не играет роли, но в целом в развитии идеи)  
и домен @avito.ru  
Пароль должен быть не менее 4 символов
## Важное о тестах бд и e2e  
- тест идет через поднятие тестовой базы в случае теста internal/db
- тест должен идти через поднятие полного окружения сервек + бек 
(в данный момент имеются не доработки этого пункта, из-за отсутсвия времени) 
Но если поднять образы, то тесты описаны в main_test.go  


# Схема базы данных

![Cхема базы данных](AvitoDBsheme.png)

# Текущие покрытие тестами
`ok  	avito	0.859s`  
`ok  	avito/internal/db	3.995s	coverage: 75.0% of statements`  
`?   	avito/internal/entity	[no test files]`  
`ok  	avito/internal/js	1.001s	coverage: 100.0% of statements`  
`ok  	avito/internal/server	1.481s	coverage: 1.2% of statements`  
`ok  	avito/internal/service	1.242s	coverage: 80.4% of statements`  
`ok  	avito/pkg/auth	2.307s	coverage: 100.0% of statements`  
`ok  	avito/pkg/jwt	2.047s	coverage: 81.8% of statements`   
# ВАЖНОЕ `total:					(statements)		58.6%`



# Нагрузочное тестирование  
Для нагрузочного тестирования будет использована утилита `wrk`  
Скрипты с параметрами находятся в папке `script`  
Running 30s test @ http://localhost:8080/api/auth
8 threads and 100 connections  
Requests/sec:  19657.42  
Transfer/sec:      4.93MB

Running 30s test @ http://localhost:8080/api/info
8 threads and 100 connections
Requests/sec:   2958.86
Transfer/sec:    137.75MB
 

Running 30s test @ http://localhost:8080/api/sendCoin
8 threads and 100 connections
Requests/sec:   5115.77
Transfer/sec:      0.90MB
 

Running 30s test @ http://localhost:8080/api/buy/t-shirt  
8 threads and 100 connections  
Requests/sec:   3368.41  
Transfer/sec:    608.53KB  

# ВАЖНО: 
- колличество юзеров на поднятом сервере < 100 000
- сейчас крипты не ведут статистику успешных и не упешных ответов(довольно сложно реализовть, 
так как мы в тесте списываем даже если по рублю - то у нас буквально за 3 секунды баланс уйдет в 0 
и в такие моменты ошибка сервера - правильная реакция )
  





