package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/qosmioo/ulr-shortener/internal/storage"
	"github.com/qosmioo/ulr-shortener/internal/usecase"
	"go.uber.org/zap"
)

type HTTPServer struct {
	usecase *usecase.URLShortenerService
	logger  *zap.Logger
}

func NewHTTPServer(usecase *usecase.URLShortenerService, logger *zap.Logger) *HTTPServer {
	return &HTTPServer{usecase: usecase, logger: logger}
}

func (s *HTTPServer) ApiEndpoints(router *mux.Router) {
	router.HandleFunc("/v1/urls", s.createShortURL).Methods("POST")
	router.HandleFunc("/v1/urls/{short_url}", s.getOriginalURL).Methods("GET")
}

func (s *HTTPServer) createShortURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OriginalURL string `json:"original_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Error("Invalid request", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	shortURL, err := s.usecase.CreateShortURL(req.OriginalURL)
	if err != nil {
		if err == storage.ErrURLExists {
			http.Error(w, "URL already exists", http.StatusConflict)
			return
		}
		s.logger.Error("Failed to create short URL", zap.Error(err))
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}
	s.logger.Info("Short URL created", zap.String("originalURL", req.OriginalURL), zap.String("shortURL", shortURL))
	json.NewEncoder(w).Encode(map[string]string{"short_url": shortURL})
}

func (s *HTTPServer) getOriginalURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["short_url"]
	originalURL, err := s.usecase.GetOriginalURL(shortURL)
	if err != nil {
		s.logger.Error("URL not found", zap.String("shortURL", shortURL), zap.Error(err))
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	s.logger.Info("Original URL retrieved", zap.String("shortURL", shortURL), zap.String("originalURL", originalURL))
	json.NewEncoder(w).Encode(map[string]string{"original_url": originalURL})
}
