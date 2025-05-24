package config

type MongoDBConfig struct {
	Host       string
	Port       string
	DBName     string
	Collection string
}

func GetMongoConfig() MongoDBConfig {
	return MongoDBConfig{
		Host:       "localhost",
		Port:       "27017",
		DBName:     "echo",
		Collection: "test",
	}
}

func (c MongoDBConfig) URI() string {
	return "mongodb://" + c.Host + ":" + c.Port + "/" + c.DBName
}
