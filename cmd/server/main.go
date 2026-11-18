package main

import (
	"fmt"
	"log"

	"app/internal/config"
	"app/internal/database"
	"app/internal/features/ai"
	"app/internal/features/auth"
	"app/internal/features/post"
	"app/internal/server"

	_ "app/docs"
)

// @title Delcom Blog API
// @version 1.0
// @description REST API untuk platform publikasi artikel Delcom Blog
// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Ketik "Bearer " diikuti token JWT
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Gagal inisialisasi basis data: %v", err)
	}

	// Inisialisasi Auth Feature
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg.AppKey)
	authHandler := auth.NewHandler(authService)

	// Inisialisasi Post Feature
	postRepo := post.NewRepository(db)
	postService := post.NewService(postRepo)
	postHandler := post.NewHandler(postService)

	// Inisialisasi AI Feature
	searchClient := ai.NewSearchClient(cfg.SearchEngineEndpoint)
	openaiClient := ai.NewOpenAIClient(cfg.OpenAIEndpoint, cfg.OpenAIKey, cfg.OpenAIModel)
	aiService := ai.NewService(searchClient, openaiClient)
	aiHandler := ai.NewHandler(aiService)

	// Inisialisasi Router
	r := server.SetupRouter(cfg, authHandler, postHandler, aiHandler)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server Delcom Blog berjalan pada port %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
