package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var s3Client *s3.S3
var ctx = context.Background()

var cancel context.CancelFunc

var Bucket = "neekcodesolutions"

func init() {
	s3Client = InitS3Client()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	MONGO_URI := os.Getenv("MONGO_URI")
	MongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	defer func() {
		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	router := http.NewServeMux()

	router.HandleFunc("POST /solution", CreateSolution)
	router.HandleFunc("GET /solutions", GetSolutions)
	router.HandleFunc("POST /like", LikeSolution)
	router.HandleFunc("POST /dislike", DislikeSolution)

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))

}
