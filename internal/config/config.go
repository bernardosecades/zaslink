package config

import (
	"fmt"
	"strconv"
)

type Config struct {
	Port            string
	SecretKey       string
	DefaultPassword string
	MongoDBURI      string
	MongoDBName     string
}

type Builder struct {
	port            string
	secretKey       string
	defaultPassword string
	mongoDBURI      string
	mongoDBName     string
}

// Port	Go doesn’t support optional parameters in function signatures so we use builder pattern
func (cb *Builder) Port(port string) *Builder {
	if port == "" {
		cb.port = "8080"
		return cb
	}
	cb.port = port
	return cb
}

func (cb *Builder) SecretKey(secretKey string) *Builder {
	cb.secretKey = secretKey
	return cb
}

func (cb *Builder) DefaultPassword(defaultPassword string) *Builder {
	cb.defaultPassword = defaultPassword
	return cb
}

func (cb *Builder) MongoDBURI(mongoDBURI string) *Builder {
	cb.mongoDBURI = mongoDBURI
	return cb
}

func (cb *Builder) MongoDBName(name string) *Builder {
	cb.mongoDBName = name
	return cb
}

func (cb *Builder) Build() (Config, error) {
	portConverted, err := strconv.Atoi(cb.port)
	if err != nil {
		return Config{}, fmt.Errorf("port %s should be a number", cb.port)
	}
	if portConverted < 0 || portConverted > 65535 {
		return Config{}, fmt.Errorf("port %s should be a number between 0 and 65535", cb.port)
	}

	if cb.secretKey == "" {
		return Config{}, fmt.Errorf("secret key should not be empty")
	}

	if cb.defaultPassword == "" {
		return Config{}, fmt.Errorf("default password should not be empty")
	}

	if cb.mongoDBURI == "" {
		return Config{}, fmt.Errorf("mongodb URI should not be empty")
	}
	if cb.mongoDBName == "" {
		return Config{}, fmt.Errorf("mongodb name should not be empty")
	}

	return Config{
		Port:            cb.port,
		SecretKey:       cb.secretKey,
		DefaultPassword: cb.defaultPassword,
		MongoDBURI:      cb.mongoDBURI,
		MongoDBName:     cb.mongoDBName,
	}, nil
}
