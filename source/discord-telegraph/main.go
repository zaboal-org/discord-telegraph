package main

import (
	"flag"
	"log"
)

// Token — токен авторизации к АПИ Дискорда для ботов
var Token string

// Получение аргументов к запуску
func init() {
	flag.StringVar(&Token, "token", "", "Authorization token to Discord API for bots")
	flag.Parse()
}

func main() {
	log.Println("Got the token:", Token)
}

// TODO: Создать и зарегстрировать команду, которая регистрирует аккаунт на Телеграфе. Здесь нужна ключ-значение БД.
