package app

import "time"

type HTTPConfig struct {
	Addr string `config:"addr" yaml:"addr"`
}

type StorageConfig struct {
	Num         int    `config:"num_of_storages" yaml:"num_of_storages"`
	StorageType string `config:"type" yaml:"type"`
}

type DBConfig struct {
	Username string `config:"user" yaml:"user"`
	Name     string `config:"name" yaml:"name"`
	Password string `config:"password" yaml:"password"`
	Host     string `config:"host" yaml:"host"`
	Port     string `config:"port" yaml:"port"`
}

type Config struct {
	HTTP            HTTPConfig    `config:"http" yaml:"http"`
	Storage         StorageConfig `config:"storage" yaml:"storage"`
	DBConfig        DBConfig      `config:"db" yaml:"db"`
	ShutdownTimeout time.Duration `config:"shutdown_timeout" yaml:"shutdown_timeout"`
	ChunksCount     int           `config:"chunks_count" yaml:"chunks_count"`
}
