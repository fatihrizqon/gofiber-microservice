package main

import (
	"fmt"
	"log"

	"github.com/fatihrizqon/gofiber-microservice/config"
	"github.com/fatihrizqon/gofiber-microservice/internal/worker"
	"github.com/hibiken/asynq"
)

func main() {
	viper := config.NewViper()
	logger := config.NewLogger(viper)

	redisOpt := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	}

	emailConfig := worker.EmailProcessorConfig{
		Host:        viper.GetString("smtp.host"),
		Port:        viper.GetInt("smtp.port"),
		Username:    viper.GetString("smtp.username"),
		Password:    viper.GetString("smtp.password"),
		SenderName:  viper.GetString("smtp.sender_name"),
		SenderEmail: viper.GetString("smtp.sender_email"),
	}

	emailProcessor := worker.NewEmailProcessor(emailConfig, logger)

	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"default": 10,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(worker.TypeEmailDelivery, emailProcessor.ProcessTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
