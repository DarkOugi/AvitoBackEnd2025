# AvitoBackend-trainee-assignment-winter-2025  

Тестовое на стажировку Golang

# Что еще нужно сделать в мечтах
1 - навесить индексов и ключей  
2 - проверку логина (?)  
3 - подумать как не допустить ситуацию создания миллионов jwt токкенов


# Как использовать Миграции GOOSE

``goose -dir db/migrations insert_merch_table sql `` - create file  
``goose -dir db/migrations up `` - run migrations  
``goose -dir /db/migrations postgres down-to <VERSION>`` - back to stable version


