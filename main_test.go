package main

import (
	"bytes"
	"context"
	"fmt"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"net/http"
	"os"
	"testing"
)

// e2e tests

//распишем сценарий покупки мерча

// 1 регистрация нового пользователя
// 1.1 покупка мерча
// 2 вход старого пользователя (удачный/не удачный)
// 2.1 покупка мерча удачная
// 2.2 покупка мерча не удачная
// 2.2.1 перевод денег от хорошего коллеги
// 2.2.2 покупка мерча

// сценарий перевода

// 1 регистрация нового пользователя
// 1.1 не удачная перевод по не существующему логину
// 1.2 не удачный перевод самому себе
// 1.3 удачный перевод
// 2 вход старого пользователя (удачный/не удачный)
// 2.1 не удачная перевод по не существующему логину
// 2.2 не удачный перевод самому себе
// 2.3 удачный перевод

func TestMain(m *testing.M) {
	os.Setenv("VOLUME", "avitoTest")
	defer func() {
		os.Setenv("VOLUME", "avito")
	}()

	compose, err := tc.NewDockerCompose("docker-compose.yaml", ".env")
	compose.WithEnv(map[string]string{
		"VOLUME":            "avitoTest",
		"POSTGRES_HOST":     "db",
		"POSTGRES_PORT":     "5432",
		"POSTGRES_USER":     "avito",
		"POSTGRES_PASSWORD": "0000",
		"POSTGRES_DB":       "avitodb",
		"HTTP_PORT":         "8080",
		"jwtSecretKey":      "VeryStrongKey",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = compose.Up(context.Background(), tc.Wait(true))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		compose.Down(context.Background(), tc.RemoveOrphans(true))
	}()
	resp, err := http.Post("http://localhost:8080/api/auth", "application/json", bytes.NewReader([]byte{}))
	//err := testDB.InitUser(context.Background(), "test1", "0000")
	//assert.Nil(t, err, "POST")
	fmt.Println(resp)
	// развернуть контейнер с сервером
	// контейнер с базой
	// они общаются
	// мы кидаем запросы
	os.Exit(m.Run())
}
