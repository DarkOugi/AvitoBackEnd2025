conn_string="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"
goose postgres "${conn_string}" up
goose postgres "${conn_string}" status
exit 0