package app

type HTTPConfig struct {
	Addr string `config:"addr"`
}

type Config struct {
	HTTP HTTPConfig `config:"http"`
}
