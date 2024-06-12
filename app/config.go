package app

type HTTPConfig struct {
	Addr string `config:"addr"`
}

type StorageConfig struct {
	Num         int    `config:"num_of_storages"`
	StorageType string `config:"type"`
}

type Config struct {
	HTTP    HTTPConfig    `config:"http"`
	Storage StorageConfig `config:"storage"`
}
