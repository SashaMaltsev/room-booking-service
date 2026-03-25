package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	jwtauth "github.com/SashaMaltsev/room-booking-service/internal/auth"
	"github.com/SashaMaltsev/room-booking-service/internal/config"
	postgresrepo "github.com/SashaMaltsev/room-booking-service/internal/repository/postgres"
	authservice "github.com/SashaMaltsev/room-booking-service/internal/service/auth"
	bookingservice "github.com/SashaMaltsev/room-booking-service/internal/service/booking"
	roomservice "github.com/SashaMaltsev/room-booking-service/internal/service/room"
	scheduleservice "github.com/SashaMaltsev/room-booking-service/internal/service/schedule"
	slotservice "github.com/SashaMaltsev/room-booking-service/internal/service/slot"
	httptransport "github.com/SashaMaltsev/room-booking-service/internal/transport/http"
)

type API struct {
	Config  config.Config
	DB      *sql.DB
	Handler http.Handler
	Server  *http.Server
}

func NewAPI(ctx context.Context, cfg config.Config) (*API, error) {
	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	userRepository := postgresrepo.NewUserRepository(db)
	roomRepository := postgresrepo.NewRoomRepository(db)
	scheduleRepository := postgresrepo.NewScheduleRepository(db)
	slotRepository := postgresrepo.NewSlotRepository(db)
	bookingRepository := postgresrepo.NewBookingRepository(db)

	tokenManager := jwtauth.NewManager(cfg.JWT.Secret, cfg.JWT.TTL)
	authSvc := authservice.New(userRepository, tokenManager)

	roomSvc := roomservice.New(roomRepository)
	scheduleSvc := scheduleservice.New(scheduleRepository, roomRepository, slotRepository, nil, 30)
	slotSvc := slotservice.New(roomRepository, scheduleRepository, slotRepository, nil)
	bookingSvc := bookingservice.New(bookingRepository, slotRepository, nil)

	handler := httptransport.NewHandler(httptransport.Dependencies{
		Auth:          authSvc,
		Rooms:         roomSvc,
		Schedules:     scheduleSvc,
		Slots:         slotSvc,
		Bookings:      bookingSvc,
		TokenVerifier: tokenManager,
	})

	server := &http.Server{
		Addr:              ":" + cfg.HTTP.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &API{
		Config:  cfg,
		DB:      db,
		Handler: handler,
		Server:  server,
	}, nil
}

func (a *API) Close() error {
	if a.DB == nil {
		return nil
	}

	return a.DB.Close()
}
