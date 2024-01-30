package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/KapDmitry/WB_L0/internal/cache"
	"github.com/KapDmitry/WB_L0/internal/config"
	"github.com/KapDmitry/WB_L0/internal/handler"
	"github.com/KapDmitry/WB_L0/internal/logger"
	"github.com/KapDmitry/WB_L0/internal/repo"
	"github.com/KapDmitry/WB_L0/internal/subscriber"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := logger.NewCustomLogger(zapConfig)
	if err != nil {
		panic(err.Error())
	}

	//DB
	postCfg := &config.PostgresConfig{}
	err = postCfg.Load("../../db/")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "localhost", postCfg.PortExt, postCfg.Username, postCfg.Password, postCfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfully connected to the database!")
	postDB := repo.NewPostgresRepo(db)

	//Cache
	memCache := cache.NewInMemoryCash()
	memCache.Recover(context.Background(), postDB)

	//Sub
	newcfg := &config.NATSConfig{}
	err = newcfg.Load("../../config/nats/")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(newcfg.ClusterID)

	sub := subscriber.Subscriber{
		Config: *newcfg,
		CTX:    context.Background(),
		Cache:  memCache,
		DB:     postDB,
		Log:    logger,
	}
	go sub.Listen()

	//Net
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("../../web/static/html")))

	localHandler := &handler.OrderHandler{
		OrderCash: memCache,
		Log:       logger,
	}
	r.HandleFunc("/order/{ID}", localHandler.GetOrder).Methods("GET")

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic("server didn't start")
	}

}
