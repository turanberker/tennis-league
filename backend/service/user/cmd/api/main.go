package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tennis-league/common/http/router"
	"tennis-league/common/lib/cache"
	"tennis-league/common/lib/database"
	"tennis-league/common/security/authmiddleware"
	"tennis-league/common/security/repository"
	"tennis-league/user-interface/grpc/pb/playerpb"
	"tennis-league/user-service/internal/controller/authhandler"
	"tennis-league/user-service/internal/controller/dashboardhandler"
	grpcctrl "tennis-league/user-service/internal/controller/grpc"
	"tennis-league/user-service/internal/controller/playerhandler"
	"tennis-league/user-service/internal/controller/userhandler"
	"tennis-league/user-service/internal/repository/postgres"
	"tennis-league/user-service/internal/repository/redis"
	"tennis-league/user-service/internal/service/auth"
	"tennis-league/user-service/internal/service/player"
	"tennis-league/user-service/internal/service/token"
	"tennis-league/user-service/internal/service/user"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	serverConfig := router.LoadServerConfig()

	db, err := database.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}

	redisClient, err := cache.NewRedis()
	if err != nil {
		log.Fatal(err)
	}
	sessionGetterRepository := repository.NewSessionGetterRepositoryImpl(redisClient)

	sessionRepository := redis.NewSessionRepository(sessionGetterRepository, redisClient)
	transactionManager := database.NewTransactionManager(db)
	userRepo := postgres.NewUserRepository(db)
	playerRepository := postgres.NewPlayerRepository(db)

	tokenService := token.NewTokenService("tennis", sessionRepository, serverConfig)
	userUC := user.NewUsecase(transactionManager, userRepo)
	authUC := auth.NewUsecase(db, userRepo, sessionRepository)
	playerUc := player.NewUsecase(transactionManager, playerRepository)

	dashboardHandler := dashboardhandler.NewDashboardHandler(playerUc)
	playerHandler := playerhandler.NewPlayerHandler(playerUc)
	authHandler := authhandler.NewAuthHandler(authUC, tokenService)
	userHandler := userhandler.NewUserHandler(userUC)

	r := router.NewRouter(serverConfig,
		authmiddleware.NewAuthMiddleware("tennis", sessionRepository),
		dashboardHandler,
		playerHandler,
		authHandler,
		userHandler)

	httpServer := &http.Server{
		Addr:    ":" + serverConfig.Port,
		Handler: r,
	}

	// 3. gRPC Sunucu Kurulumu ve Port Tanımı
	grpcPort := ":50051"
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("gRPC dinleme hatası: %v", err)
	}

	grpcServer := grpc.NewServer()

	// gRPC controller'ı kaydediyoruz
	grpcctrl.NewPlayerGrpcServer(playerUc)

	playerGrpcImpl := grpcctrl.NewPlayerGrpcServer(playerUc)
	playerpb.RegisterPlayerServiceServer(grpcServer, playerGrpcImpl)

	// 4. Eş Zamanlı Çalıştırma ve Graceful Shutdown Yönetimi
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)

	// HTTP (Gin) Başlatılıyor
	g.Go(func() error {
		log.Printf("HTTP (Gin) Sunucusu çalışıyor: %s", serverConfig.Port)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	// gRPC Başlatılıyor
	g.Go(func() error {
		log.Printf("gRPC Sunucusu çalışıyor: %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			return err
		}
		return nil
	})

	// Sinyal Geldiğinde Kapatma Yönetimi
	g.Go(func() error {
		<-gCtx.Done() // Sinyal gelene kadar bloke olur

		log.Println("Kapatma sinyali alındı, servisler durduruluyor...")

		// HTTP sunucusunu kapatmak için 5 saniye tolerans tanıyalım
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP shutdown hatası: %v", err)
		} else {
			log.Println("HTTP sunucusu kapatıldı.")
		}

		// gRPC sunucusunu kibarca kapatıyoruz
		grpcServer.GracefulStop()
		log.Println("gRPC sunucusu kapatıldı.")

		return nil
	})

	// Herhangi bir sunucuda hata olursa ya da graceful shutdown bittiğinde burası çözülür
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("Uygulama beklenmedik bir hata nedeniyle kapandı: %v", err)
	}

	log.Println("Sunucu tamamen sonlandırıldı.")
}
