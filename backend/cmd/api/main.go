package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/ryohighbridge/learn-github-copilot/backend/internal/handler"
	"github.com/ryohighbridge/learn-github-copilot/backend/internal/repository"
	"github.com/ryohighbridge/learn-github-copilot/backend/internal/service"
)

func main() {
	// データベース接続
	db, err := repository.NewPostgresDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// リポジトリとサービスの初期化
	eventRepo := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepo)
	calendarService := service.NewCalendarService()

	// ハンドラーの初期化
	eventHandler := handler.NewEventHandler(eventService)
	calendarHandler := handler.NewCalendarHandler(calendarService)

	// ルーターの設定
	r := mux.NewRouter()

	// カレンダーAPI
	r.HandleFunc("/api/calendar/{year:[0-9]+}/{month:[0-9]+}", calendarHandler.GetCalendar).Methods("GET")
	r.HandleFunc("/api/holidays/{year:[0-9]+}", calendarHandler.GetHolidays).Methods("GET")

	// イベントAPI
	r.HandleFunc("/api/events", eventHandler.GetEvents).Methods("GET")
	r.HandleFunc("/api/events", eventHandler.CreateEvent).Methods("POST")
	r.HandleFunc("/api/events/{id:[0-9]+}", eventHandler.GetEvent).Methods("GET")
	r.HandleFunc("/api/events/{id:[0-9]+}", eventHandler.UpdateEvent).Methods("PUT")
	r.HandleFunc("/api/events/{id:[0-9]+}", eventHandler.DeleteEvent).Methods("DELETE")

	// ヘルスチェック
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// CORS設定
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
