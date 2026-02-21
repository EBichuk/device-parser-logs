package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Directory
	Db
	HttpServer
	Worker
}

type Directory struct {
	DirectoryTsv string `env:"DIRECTORY_TSV"`
	DirectoryPdf string `env:"DIRECTORY_PDF"`
}

type Db struct {
	Name     string `env:"DB_NAME"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
}

type HttpServer struct {
	Addr string `env:"HTTP_ADDRESS"`
}

type Worker struct {
	Interval string `env:"INTERVAL"`
	Workers  int    `env:"WORKERS_POOL"`
}

func LoadConfig() (Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
