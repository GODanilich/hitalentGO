package main

import (
	"log"
	"net/http"

	"GODanilich/hitalentGO/internal/config"
	"GODanilich/hitalentGO/internal/db"
	"GODanilich/hitalentGO/internal/http/handlers"
	"GODanilich/hitalentGO/internal/http/middleware"
	"GODanilich/hitalentGO/internal/http/router"
	"GODanilich/hitalentGO/internal/repo"
	"GODanilich/hitalentGO/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	gdb, err := db.Open(cfg.DatabaseURL, cfg.GormLog)
	if err != nil {
		log.Fatal(err)
	}

	// repos
	depRepo := repo.NewDepartmentRepo(gdb)
	empRepo := repo.NewEmployeeRepo(gdb)

	// services
	depSvc := service.NewDepartmentService(depRepo, empRepo)

	// handlers
	depH := handlers.NewDepartmentsHandler(depSvc)
	healthH := handlers.NewHealthHandler(gdb)

	// router
	rt := router.New()
	rt.Handle(http.MethodGet, "/health", http.HandlerFunc(healthH.Health))
	rt.Handle(http.MethodPost, "/departments", http.HandlerFunc(depH.CreateDepartment))
	rt.Handle(http.MethodPost, "/departments/{id}/employees", http.HandlerFunc(depH.CreateEmployee))

	// middleware
	h := middleware.Chain(rt, middleware.Recover(), middleware.Logger())

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: h,
	}

	log.Printf("listening on %s", cfg.HTTPAddr)
	log.Fatal(srv.ListenAndServe())
}
