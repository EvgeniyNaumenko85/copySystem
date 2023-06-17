package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	tasks "tasks_API"
	"tasks_API/configs"
	"tasks_API/db"
	"tasks_API/pkg/handler"
	"tasks_API/pkg/logger"
	"tasks_API/pkg/repository"
	"tasks_API/pkg/service"
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
	service := service.NewService(repos)
	handlers := handler.NewHandler(service)

	srv := new(tasks.Server)

	// Поднимаем таблицы в БД
	logger.Info.Println("Поднимаем таблицы в БД")
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

	// Закрытие соединения с базой данных
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
