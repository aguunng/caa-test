package api

import (
	"context"
	"fmt"
	"caa-test/internal/api/resp"
	"caa-test/internal/client"
	"caa-test/internal/config"
	"caa-test/internal/qismo"
	"caa-test/internal/room"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func NewServer() *Server {
	cfg := config.Load()

	client := client.New()
	qismo := qismo.New(client, cfg.Qiscus.Omnichannel.URL, cfg.Qiscus.AppID, cfg.Qiscus.SecretKey)

	roomSvc := room.NewService(qismo)
	roomHandler := room.NewHttpHandler(roomSvc)

	r := http.NewServeMux()
	r.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			resp.WriteJSON(w, http.StatusNotFound, resp.HTTPError{
				StatusCode: 404,
				Message:    "Not Found",
			})
			return
		}

		resp.WriteJSON(w, http.StatusOK, resp.HTTPSuccess{
			Message: "Service Running",
		})
	}))

	r.Handle("POST /api/v1/caa", http.HandlerFunc(roomHandler.WebhookCaa))
	r.Handle("POST /api/v1/mark_as_resolved", http.HandlerFunc(roomHandler.WebhookMarkResolved))
	r.Handle("GET /api/v1/rooms", http.HandlerFunc(roomHandler.GetRooms))
	r.Handle("GET /api/v1/first_room", http.HandlerFunc(roomHandler.FirstCustomerRoom))
	r.Handle("POST /api/v1/update-max-customer", http.HandlerFunc(roomHandler.UpdateMaxCustomerHandler))

	return &Server{router: r}
}

type Server struct {
	router *http.ServeMux
}

// Run method of the Server struct runs the HTTP server on the specified port. It initializes
// a new HTTP server instance with the specified port and the server's router.
func (s *Server) Run(port int) {
	addr := fmt.Sprintf(":%d", port)

	h := chainMiddleware(
		s.router,
		recoverHandler,
		loggerHandler(func(w http.ResponseWriter, r *http.Request) bool { return r.URL.Path == "/" }),
		realIPHandler,
		requestIDHandler,
		corsHandler,
	)

	httpSrv := http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info().Msg("server is shuting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		httpSrv.SetKeepAlivesEnabled(false)
		if err := httpSrv.Shutdown(ctx); err != nil {
			log.Fatal().Err(err).Msg("could not gracefully shutdown the server")
		}
		close(done)
	}()

	log.Info().Msgf("server serving on port %d", port)
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msgf("could not listen on %s", addr)
	}

	<-done
	log.Info().Msg("server stopped")

}
