package cron

import (
	"caa-test/internal/client"
	"caa-test/internal/config"
	"caa-test/internal/qismo"
	"caa-test/internal/room"
	"context"
	"time"

	"caa-test/internal/postgres"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func NewServer() *Server {
	cfg := config.Load()

	db := postgres.NewGORM(cfg.Database)
	postgres.Migrate(db)

	client := client.New()
	qismo := qismo.New(client, cfg.Qiscus.Omnichannel.URL, cfg.Qiscus.AppID, cfg.Qiscus.SecretKey)

	roomRepo := room.NewRepository(db)
	roomSvc := room.NewService(roomRepo, qismo, cfg)

	return &Server{
		svc: roomSvc,
	}
}

type Server struct {
	svc *room.Service
}

// Run starts the cron job and schedules it to execute every minute.
func (c *Server) Run() {
	log.Info().Msg("cron is started")

	s := gocron.NewScheduler(time.UTC)
	s.Every(60).Second().Do(func() {
		reqID := uuid.New().String()
		ctx := log.With().Str("request_id", reqID).Logger().WithContext(context.Background())

		err := c.svc.AssignAgentEvent(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("error handle assigned room to agent : %s", err.Error())
		}
	})

	s.StartBlocking()
}