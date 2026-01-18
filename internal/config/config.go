package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Address string
}

// /env-default:"production
type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" `
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
	Configflag:=flag.String("config","","path to config file")
	flag.Parse()
	configPath = *Configflag

	if configPath == "" {
		log.Fatal("config path is not set")
	}
	}


	if _,err:=os.Stat(configPath);os.IsNotExist(err){//it checks if the file exists
		log.Fatalf("config file does not exist :%s",configPath)

	}

	var cfg Config

	err:=cleanenv.ReadConfig(configPath,&cfg)//it reads the config file and stores it in the cfg variable
	if err!=nil{
		log.Fatalf("cannot read config :%s",err.Error())
	}

	return &cfg
}