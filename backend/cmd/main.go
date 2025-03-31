package main

import (
	"log"
	"net/http"
	"tourism/pkg/routes"
	"tourism/pkg/tasks"

	"github.com/gorilla/mux"
	"github.com/hibiken/asynq"
)

func startWorker() {
	redisConn := asynq.RedisClientOpt{Addr: "redis:6379"}
	scheduler := asynq.NewScheduler(redisConn, nil)

	srv := asynq.NewServer(redisConn, asynq.Config{
		Concurrency: 10,
	})

	mux := asynq.NewServeMux()
	_, err := scheduler.Register("@every 10m", asynq.NewTask(tasks.CheckResrvations, nil))
	if err != nil {
		log.Fatalf("Could not schedule task: %v", err)
	}

	log.Println("Starting periodic task scheduler...")
	if err := scheduler.Start(); err != nil {
		log.Fatalf("Scheduler failed: %v", err)
	}
	mux.HandleFunc(tasks.CheckResrvations, tasks.NewPeriodicTask)

	log.Println("Worker started...")

	if err := srv.Run(mux); err != nil {
		log.Fatalf("Worker error: %v", err)
	}
}

func main() {
	go startWorker()
	r := mux.NewRouter()
	routes.RegisterRoutes(r)
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err == nil {
			log.Println("Registered route:", path)
		} else {
			log.Println("Error registering route:", err)
		}
		return nil
	})
	log.Fatal(http.ListenAndServe(":8000", r))
}