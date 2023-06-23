package main

import (
	"context"
	"copySys"
	"copySys/configs"
	"copySys/db"
	"copySys/pkg/handler"
	"copySys/pkg/logger"
	"copySys/pkg/repository"
	"copySys/pkg/service"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	configs.PutAdditionalSettings()
	logger.Init()

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	db.StartDbConnection(db.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(copySys.Server)

	// Поднимаем таблицы в БД
	logger.Info.Println("Raising tables in the database")
	if err := db.Up(); err != nil {
		log.Fatalf("Error while migrating tables, err is: %s", err.Error())
		return
	}

	//if err := db.Down(); err != nil {
	//	log.Fatalf("Error while dropping tables, err is: %s", err.Error())
	//	return
	//}

	go func() {
		if err := srv.Run(viper.GetString("PORT"), handlers.InitRoutes()); err != nil {
			//logger.Info.Println(err.Error())
			//log.Fatalf("Error occured while running http server: %s", err.Error())

			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	if err := db.CloseDbConnection(); err != nil {
		fmt.Errorf("error occurred on database connection closing: %s", err.Error())
	}

	fmt.Println("Shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Errorf("error occurred on server shutting down: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
