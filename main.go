package main

import (
	"github.com/cristianchaparroa/secret/api"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.SetLevel(logrus.DebugLevel)
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	s := api.NewSecretServer()
	s.Setup()
	s.Run()
	defer s.Close()
}
