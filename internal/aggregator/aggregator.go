package aggregator

import (
	"context"
	"github.com/arjanvaneersel/getlive/internal/entry"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
	"time"
)

type Aggregator interface {
	Aggregate(chan entry.NewEntry) error
}

type Server struct {
	db          *sqlx.DB
	logger      *log.Logger
	stopChan    chan struct{}
	domains     []string
	Aggregators []Aggregator
}

func New(db *sqlx.DB, logger *log.Logger, domains []string, a []Aggregator) *Server {
	return &Server{
		db:          db,
		logger:      logger,
		stopChan:    make(chan struct{}, 1),
		domains:     domains,
		Aggregators: a,
	}
}

func (s *Server) IsAllowedURL(url string) bool {
	if len(s.domains) == 0 {
		return true
	}

	for _, d := range s.domains {
		if strings.Contains(strings.ToLower(url), strings.ToLower(d)) {
			return true
		}
	}

	return false
}

func (s *Server) Run() {
	entryChan := make(chan entry.NewEntry)
	for _, a := range s.Aggregators {
		go func() {
			a.Aggregate(entryChan)
		}()
	}

	s.logger.Printf("running")

	for {
		select {
		case e := <-entryChan:
			if !s.IsAllowedURL(e.URL) {
				continue
			}

			se, err := entry.Create(context.Background(), s.db, e, time.Now())
			if err != nil {
				s.logger.Printf("Couldn't save entry to database: %v", err)
				continue
			}
			s.logger.Printf("saved entry %s to database: %s", se.ID, se.Title)
		case <-s.stopChan:
			break
		}
	}
}

func (s *Server) Shutdown() {
	s.stopChan <- struct{}{}
}
