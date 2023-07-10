package config

import (
	"flag"
	"os"
)

type Parameters struct {
	FlagRunAddr     string
	ShortURL        string
	FileStoragePath string
	DBStoragePath   string
}

func ParseFlags() Parameters {

	var (
		flagRunAddr     string
		shortURL        string
		fileStoragePath string
		dbStoragePath   string
	)

	// регистрируем переменную FlagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")

	// регистрируем переменную ShortURL
	// как аргумент -b со значением http://localhost:8080/ по умолчанию
	flag.StringVar(&shortURL, "b", "http://localhost:8080", "response URL")

	// регистрируем переменную FileStoragePath
	// как аргумент -f со значением /tmp/short-url-db.json по умолчанию
	flag.StringVar(&fileStoragePath, "f", "/tmp/short-url-db.json", "path to file storage")

	// регистрируем переменную DBStoragePath
	// как аргумент -d со значением http://localhost:8080/ по умолчанию
	flag.StringVar(&dbStoragePath, "d", "", "path to connect DB")

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	// возьмем из переменной окружения SERVER_ADDRESS адрес запуска сервера
	// переопределим переменную из переменного окружения, если есть
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	// переопределим базовый адрес результирующего сокращенного URL если есть
	if envShortURL := os.Getenv("BASE_URL"); envShortURL != "" {
		shortURL = envShortURL
	}

	// переопределим переменную из переменного окружения,
	// если есть для пути сохранения файла
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		fileStoragePath = envFileStoragePath
	}

	// переопределим переменную из переменного окружения,
	// если есть для пути коннекта к БД
	if envDBStoragePath := os.Getenv("DATABASE_DSN"); envDBStoragePath != "" {
		dbStoragePath = envDBStoragePath
	}

	return Parameters{
		FlagRunAddr:     flagRunAddr,
		ShortURL:        shortURL,
		FileStoragePath: fileStoragePath,
		DBStoragePath:   dbStoragePath,
	}
}
