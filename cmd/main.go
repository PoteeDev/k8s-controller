package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/PoteeDev/k8s-controller/internal/handlers"
	"github.com/zitadel/oidc/v3/pkg/client/rs"
	http_mw "github.com/zitadel/zitadel-go/v2/pkg/api/middleware/http"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
)

type key int

const (
	requestIDKey key = 0
)

func DefaultEnv(name string, value string) string {
	if env := os.Getenv(name); env != "" {
		return env
	}
	return value
}

var (
	issuer = flag.String("issuer", DefaultEnv("ZITADEL_HOST", "localhost:8080"), "issuer of your ZITADEL instance (in the form: https://<instance>.zitadel.cloud or https://<yourdomain>)")
	port   = flag.String("port", DefaultEnv("PORT", "80"), "")
)

var (
	healthy int32
)

func main() {
	flag.Parse()

	introspection, err := http_mw.NewIntrospectionInterceptor(*issuer, middleware.OSKeyPath())
	if err != nil {
		log.Fatal("apikey error:", err)
	}
	log.Printf("use zitadel issuer: %s\n", *issuer)

	listenAddr := fmt.Sprintf("0.0.0.0:%s", *port)
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Printf("Server is starting...")

	provider, err := rs.NewResourceServerFromKeyFile(context.TODO(), *issuer, middleware.OSKeyPath())
	if err != nil {
		logger.Fatalf("error creating token source %s", err.Error())
	}

	s := handlers.InitServer(provider)

	router := http.NewServeMux()
	router.HandleFunc("/public", s.Ping)
	router.HandleFunc("/protected", introspection.HandlerFunc(s.Ping))
	router.HandleFunc("/stand/deploy", introspection.HandlerFunc(s.DeployStand))
	router.HandleFunc("/stand/destroy", introspection.HandlerFunc(s.DestroyStand))
	router.HandleFunc("/stand/info", introspection.HandlerFunc(s.InfoHandler))

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      tracing(nextRequestID)(logging(logger)(router)),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready to handle requests at", listenAddr)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Println("Server stopped")
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
