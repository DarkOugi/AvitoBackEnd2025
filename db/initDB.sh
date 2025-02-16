for file in ../../db/migrations/*; do
    awk '/-- \+goose Up/,/-- \+goose StatementEnd/' "$file" |
    sed -n '/-- +goose StatementBegin/,/-- +goose StatementEnd/p' |
    sed '1d; $d'
done > var1
