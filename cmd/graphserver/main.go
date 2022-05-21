package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"

	"github.com/Littlefisher619/cosdisk/config"
	"github.com/Littlefisher619/cosdisk/graph"
	"github.com/Littlefisher619/cosdisk/graph/auth"
	"github.com/Littlefisher619/cosdisk/graph/generated"

	//"github.com/Littlefisher619/cosdisk/model"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/Littlefisher619/cosdisk/service"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

const defaultPort = 8080

func main() {

	var (
		port            = flag.Int("port", defaultPort, "The server port")
		host            = flag.String("host", "127.0.0.1", "The server host")
		checkOriginHost = flag.String("check-origin-host", "", "The server host")
		allowedOrigins  = flag.String("allowed-origins", "", "The allowed origins")
		configPath      = flag.String("config-file", "./config.toml", "The toml config file")
	)
	flag.Parse()

	if port == nil {
		*port = defaultPort
	}
	c, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Println("parse config file error: ", err)
		os.Exit(1)
	}
	service := service.New(c)
	jwtAuthManager := auth.New()
	router := chi.NewRouter()
	router.Use(jwtAuthManager.Middleware)

	origins := []string{}
	if *allowedOrigins == "" {
		origins = append(origins, "*")
	} else {
		origins = append(origins, *allowedOrigins)
	}

	router.Use(cors.New(cors.Options{
		AllowedHeaders:   []string{"*"},
		AllowedOrigins:   origins,
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	graphqlConfig := generated.Config{Resolvers: &graph.Resolver{service, jwtAuthManager}}
	graphqlConfig.Directives.Auth = jwtAuthManager.AuthDirective
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(graphqlConfig))
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				if *checkOriginHost == "" {
					return true
				}
				return r.Host == *checkOriginHost
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://%s:%d/ for GraphQL playground", *host, *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), router))
	// log.Fatal(http.ListenAndServe(":8080", router))
}
