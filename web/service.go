package web

import (
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/surajjain36/log_server/infra"
	"github.com/surajjain36/log_server/misc"
)

// Service HTTP server info
type Service struct {
	shutdownChan chan bool
	domain       string

	router *gin.Engine
	wg     sync.WaitGroup
	mdb    *infra.Mongo
	//mysql     *infra.Mysql
	//queue     *infra.QueueConfig
	AppName   string
	Version   string
	BuildTime string
}

// NewService Create a new service
func NewService(conf *misc.Config) (*Service, error) {
	mdb, err := infra.NewMongo(&conf.Mongo)
	if err != nil {
		log.WithError(err).Error("Failed to connect to MongoDB")
		return nil, err
	}

	s := &Service{
		router:       gin.New(),
		mdb:          mdb,
		shutdownChan: make(chan bool),
		domain:       conf.HTTP.Domain,
	}

	s.router.Use(gin.Logger())

	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "OPTIONS", "POST"},
		AllowHeaders:     []string{"origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s.router.GET("/", s.index)
	v1 := s.router.Group("/v1")
	{
		v1.POST("/log", s.WriteLog)
		v1.GET("/log", s.ReadLog)
	}

	return s, nil
}

// Start the web service
func (s *Service) Start(address string) error {
	return s.router.Run(address)
}

// Close all threads and free up resources
func (s *Service) Close() {
	close(s.shutdownChan)

	s.wg.Wait()

	//s.rc.Close()
}
