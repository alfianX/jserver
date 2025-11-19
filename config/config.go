package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	CnfGlob ConfigGlobal
	CnfLoc  ConfigLocal
}

type ConfigGlobal struct {
	Mode          string `envconfig:"MODE" default:"debug"`
	Database      string `envconfig:"MYSQL_DSN"`
	DatabaseParam string `envconfig:"MYSQL_DSN_PARAM"`
	TimeoutTrx    int    `envconfig:"TIMEOUT_TRX" default:"50"`
	TestTimeout   int    `envconfig:"TEST_TIMEOUT" default:"0"`
	OdooURL       string `envconfig:"ODOO_URL"`
}

type ConfigLocal struct {
	ListenPort int    `envconfig:"LISTEN" default:"88"`
	Host       string `envconfig:"HOST"`
}

func NewParsedConfig() (Config, error) {
	_ = godotenv.Load(".env")
	cnf := Config{}
	cnfLoc := ConfigLocal{}
	err := envconfig.Process("", &cnfLoc)
	if err != nil {
		return Config{}, err
	}
	cnf.CnfLoc = cnfLoc

	_ = godotenv.Load("../.env")
	cnfGlob := ConfigGlobal{}
	err = envconfig.Process("", &cnfGlob)
	if err != nil {
		return Config{}, err
	}
	cnf.CnfGlob = cnfGlob

	return cnf, nil
}
