// This package for registrate environment variable and run arguments.
package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
)

type Parameters struct {
	FlagRunAddr     string `json:"server_address"`
	ShortURL        string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DBpath          string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
	Path            string
}

// This func registrate environment variable and run arguments.
func ParseFlags() Parameters {

	var (
		flagRunAddr     string
		shortURL        string
		fileStoragePath string
		dbpath          string
		enableHTTPS     string
		httpsVar        bool
		path            string
	)

	// регистрируем переменную Path
	// как аргумент -c со значением
	flag.StringVar(&path, "c", "", "Config path")

	if envPath := os.Getenv("CONFIG"); envPath != "" {
		path = envPath
	}

	if len(path) != 0 {
		return *getConfig(path)
	}

	// регистрируем переменную FlagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")

	// регистрируем переменную ShortURL
	// как аргумент -b со значением http://localhost:8080/ по умолчанию
	flag.StringVar(&shortURL, "b", "http://localhost:8080", "response URL")

	// регистрируем переменную FileStoragePath
	// как аргумент -f со значением /tmp/short-url-db.json по умолчанию
	flag.StringVar(&fileStoragePath, "f", "/tmp/short-url-db.json", "path to file storage")

	// регистрируем переменную DBpath
	// как аргумент -d со значением
	flag.StringVar(&dbpath, "d", "", "path to database storage")

	// регистрируем переменную EnableHTTPS
	// как аргумент -s со значением
	flag.StringVar(&enableHTTPS, "s", "", "enable https")

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
	// если есть для пути сохранения файла
	if envDBStoragePath := os.Getenv("DATABASE_DSN"); envDBStoragePath != "" {
		dbpath = envDBStoragePath
	}

	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		enableHTTPS = envEnableHTTPS
	}

	https, err := strconv.ParseBool(enableHTTPS)
	if err != nil {
		httpsVar = false
	}

	httpsVar = https

	return Parameters{
		FlagRunAddr:     flagRunAddr,
		ShortURL:        shortURL,
		FileStoragePath: fileStoragePath,
		DBpath:          dbpath,
		EnableHTTPS:     httpsVar,
	}
}

func getConfig(filename string) *Parameters {
	var cfg *Parameters

	file, err := os.ReadFile(filename)
	if err != nil {
		log.Print("read file faild ", err)
		return nil
	}

	err = json.Unmarshal(file, cfg)
	if err != nil {
		log.Print("unmarshal config fail ", err)
		return nil
	}
	return cfg
}
