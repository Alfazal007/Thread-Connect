package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Alfazal007/worker/controllers"
	"github.com/Alfazal007/worker/internal/database"
	"github.com/Alfazal007/worker/send_mail"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	portNumber := os.Getenv("PORT")
	if portNumber == "" {
		log.Fatal("Error getting port number from the env variables")
	}
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("Error getting database url from the env variables")
	}
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		log.Fatal("Error getting the redis url from the env variables")
	}
	grpcPort := os.Getenv("GRPC_POST")
	if grpcPort == "" {
		log.Fatal("Error getting the grpc port from the env variables")
	}
	redisPassword := os.Getenv("REDIS_PASSWORD") // TODO:: ADD ACTUAL PASSWORD AND UNCOMMENT THE LINE
	// if redisPassword == "" {
	// log.Fatal("Error getting the redis password from the env variables")
	// }
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error opening database connection", err)
	}
	apiCfg := controllers.ApiCfg{DB: database.New(conn)}
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword,
		DB:       0,
	})
	apiCfg.Rdb = rdb
	// grpc client
	grpcConnection, err := grpc.NewClient("127.0.0.1:"+grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	clientConn := send_mail.NewSendMailServiceClient(grpcConnection)
	apiCfg.GrpcClient = clientConn
	for {
		data := apiCfg.Rdb.BRPop(context.Background(), 0, "worker")
		val, err := data.Result()
		if err != nil {
			fmt.Println("Error popping from Redis:", err)
			return
		}

		var dataInsert controllers.Incoming
		err = json.Unmarshal([]byte(val[1]), &dataInsert)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		if dataInsert.Type == "repost" {
			apiCfg.ChangeRepostCount(dataInsert)
		} else if dataInsert.Type == "like" {
			apiCfg.ChangeLikeCount(dataInsert)
		} else if dataInsert.Type == "notification" {
			to, username := apiCfg.GetFollowers(dataInsert)
			if username == "" || len(to) == 0 {
				continue
			}
			// need to update
			data := send_mail.CreateRequest{
				Fromusername: username,
				Link:         "http://localhost:8000/" + dataInsert.TweetId,
				To:           to,
			}
			resp, err := apiCfg.GrpcClient.SendMail(context.Background(), &data)
			if err != nil {
				println("Error sending the notifications", err.Error())
			}
			println("Status : ", resp.Done)
		} else {
			println("Improper request")
		}
	}
}
