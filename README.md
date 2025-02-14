# AvitoBackend-trainee-assignment-winter-2025  

Тестовое на стажировку Golang

# Что еще нужно сделать в мечтах 
1 - подумать как не допустить ситуацию создания миллионов jwt токкенов

# важное о Аунтификации (регистрации/входе)  
Выдвинутых требований к логину со сторону ТЗ не было  
Но и хранить любой логин мы не должны, как минимум это странно  
Логин приведен в строгий формат second_name.first_name  
(если у нас есть уже сотрудник с данным именем то добавляются цифры,  
сейчас это не играет роли, но в целом в развитии идеи)  
и домен @avito.ru  
Пароль должен быть не менее 4 символов

# Как использовать Миграции GOOSE

``goose -dir db/migrations insert_merch_table sql `` - create file  
``goose -dir db/migrations up `` - run migrations  
``goose -dir /db/migrations postgres down-to <VERSION>`` - back to stable version  

# Схема базы данных

![Cхема базы данных](AvitoDBsheme.png)

# Текущие покрытие тестами


``avito           coverage: 0.0% of statements``  
``?       avito/internal/entity   [no test files]``  
``ok      avito/internal/db       0.428s  coverage: 45.7% of statements``  
``ok      avito/internal/js       0.626s  coverage: 100.0% of statements``  
``ok      avito/internal/server   1.333s  coverage: 0.0% of statements [no tests to run]``  
``ok      avito/internal/service  1.579s  coverage: 80.0% of statements``  
``ok      avito/pkg/auth  1.112s  coverage: 100.0% of statements``  
``ok      avito/pkg/jwt   0.868s  coverage: 88.9% of statements``

