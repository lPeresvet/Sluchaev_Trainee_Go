package main

import (
	"context"
	"fmt"
	"log"
	"os"
	_ "trainee/docs"
	"trainee/internal/repository/postgres"
	"trainee/internal/service"
	"trainee/internal/transport/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

const (
	CREATE_SEGMENTS_SCHEMA = `
	CREATE TABLE IF NOT EXISTS segments (
	id BIGSERIAL PRIMARY KEY,  
	slug VARCHAR(100), 
	status SMALLINT DEFAULT 0);
	`
	CREATE_USERS_SCHEMA = `
	CREATE TABLE IF NOT EXISTS accounts (
	id BIGINT PRIMARY KEY);
	`
	CREATE_USER_SEGMENT_SCHEMA = `
	CREATE TABLE IF NOT EXISTS user_segment (
	user_id BIGINT, 
	segment_id BIGINT, 
	FOREIGN KEY (user_id) REFERENCES accounts(id), 
	FOREIGN KEY (segment_id) REFERENCES segments(id));
	`
	CREATE_LOG_SCHEMA = `
	CREATE TABLE IF NOT EXISTS segment_log (
	id BIGSERIAL PRIMARY KEY, 
	user_id BIGINT, 
	segment_id BIGINT, 
	operation SMALLINT, 
	operation_time TIMESTAMP, 
	FOREIGN KEY (user_id) REFERENCES accounts(id), 
	FOREIGN KEY (segment_id) REFERENCES segments(id));
	`
	DB_URL = "postgres://developer:developer@database:5432/avito_db"
)

// @title Test task
// @version 2.0
// @description This is a avito task
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	if err := SetupViper(); err != nil {
		log.Fatal(err.Error())
	}

	dbUrl := viper.GetString("postgres.database.url")

	conn, err := initDB(dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault)

	userRepository := postgres.NewUserRepository(conn)
	segmentRepository := postgres.NewSegmentRepository(conn)
	logrepository := postgres.NewLogRepository(conn)
	userService := service.NewUserService(userRepository)
	segmentService := service.NewSegmentService(segmentRepository)
	logService := service.NewLogService(logrepository, userRepository)
	userHandler := handler.NewUserHandler(userService)
	segmentHandler := handler.NewSegmentHandler(segmentService)
	logHandler := handler.NewLogHandler(logService)

	userHandler.InitRoutes(app)
	segmentHandler.InitRoutes(app)
	logHandler.InitRoutes(app)

	port := viper.GetString("http.port")

	if err = app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}

}

func initDB(dbUrl string) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}
	conn, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

	_, err = conn.Query(context.Background(), CREATE_SEGMENTS_SCHEMA)
	if err != nil {
		fmt.Println("Segments not created")
		return nil, err
	}

	_, err = conn.Query(context.Background(), CREATE_USERS_SCHEMA)
	if err != nil {
		fmt.Println("Users not created")
		return nil, err
	}

	_, err = conn.Query(context.Background(), CREATE_USER_SEGMENT_SCHEMA)
	if err != nil {
		fmt.Println("User_Segment not created")
		return nil, err
	}

	_, err = conn.Query(context.Background(), CREATE_LOG_SCHEMA)
	if err != nil {
		fmt.Println("Log_Schema not created")
		return nil, err
	}
	return conn, nil
}

func SetupViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
