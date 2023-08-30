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
		status smallint default 0);
		`
	CREATE_USERS_SCHEMA = "CREATE TABLE IF NOT EXISTS accounts (" +
		"id bigint primary key);"
	CREATE_USER_SEGMENT_SCHEMA = "CREATE TABLE IF NOT EXISTS user_segment (" +
		"user_id bigint, " +
		"segment_id bigint, " +
		"foreign key (user_id) references accounts(id), " +
		"foreign key (segment_id) references segments(id));"
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

	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	conn, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = initDB(conn)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault)

	userRepository := postgres.NewUserRepository(conn)
	segmentRepository := postgres.NewSegmentRepository(conn)
	userService := service.NewUserService(userRepository)
	segmentService := service.NewSegmentService(segmentRepository)
	userHandler := handler.NewUserHandler(userService)
	segmentHandler := handler.NewSegmentHandler(segmentService)

	userHandler.InitRoutes(app)
	segmentHandler.InitRoutes(app)

	port := viper.GetString("http.port")

	if err = app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}

}

func initDB(conn *pgxpool.Pool) error {
	_, err := conn.Query(context.Background(), CREATE_SEGMENTS_SCHEMA)
	if err != nil {
		fmt.Println("Segments not created")
		return err
	}

	_, err = conn.Query(context.Background(), CREATE_USERS_SCHEMA)
	if err != nil {
		fmt.Println("Users not created")
		return err
	}

	_, err = conn.Query(context.Background(), CREATE_USER_SEGMENT_SCHEMA)
	if err != nil {
		fmt.Println("User_Segment not created")
		return err
	}
	return nil
}

func SetupViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
