# AvitoBackend-trainee-assignment-winter-2025  

Тестовое на стажировку Golang

# Что еще нужно сделать в мечтах
1 - навесить индексов и ключей  
2 - проверку логина (?)  
3 - подумать как не допустить ситуацию создания миллионов jwt токкенов

# важное о Аунтификации (регистрации/входе)  
Выдвинутых требований к логину со сторону ТЗ не было  
Но и хранить любой логин мы не должны, как минимум это странно  
Логин приведен в строгий формат second_name.first_name  
(если у нас есть уже сотрудник с данным именем то добавляются цифры,  
сейчас это не играет роли, но в целом в развитии идеи)  
и домен @avito.ru  

# Как использовать Миграции GOOSE

``goose -dir db/migrations insert_merch_table sql `` - create file  
``goose -dir db/migrations up `` - run migrations  
``goose -dir /db/migrations postgres down-to <VERSION>`` - back to stable version


