package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type Server struct {
	MongoClient *mongo.Client
	Router      *mux.Router
}

func (server *Server) initialize() {
	server.MongoClient = initializeDBConnection("mongodb://localhost:27017/")
	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

func initializeDBConnection(uri string) *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("Error occurred in connecting to mongo")
		log.Fatal(err)
	}
	return mongoClient
}

func (server *Server) run(port string) {
	log.Fatal(http.ListenAndServe("localhost:"+port, server.Router))
}
