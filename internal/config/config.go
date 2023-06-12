package config

import (
	"flag"
	"os"
)

var (
	FlagRunAddr string
	ShortURL    string
)

func ParseFlags() {

	// регистрируем переменную FlagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")

	// регистрируем переменную ShortURL
	// как аргумент -b со значением http://localhost:8080/ по умолчанию
	flag.StringVar(&ShortURL, "b", "http://localhost:8080", "response URL")

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	// возьмем из переменной окружения SERVER_ADDRESS адрес запуска сервера
	// переопределим переменную из переменного окружения, если есть
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}

	// переопределим базовый адрес результирующего сокращенного URL если есть
	if envShortUrl := os.Getenv("BASE_URL"); envShortUrl != "" {
		ShortURL = envShortUrl
	}

}
