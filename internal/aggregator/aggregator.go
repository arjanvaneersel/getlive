package aggregator

import (
	"github.com/arjanvaneersel/getlive/internal/entry"
	"log"
)

type Aggregator interface {
	Aggregate(chan *entry.Entry) error
}

type Server struct {
	logger      *log.Logger
	Aggregators []Aggregator
	stopChan    chan struct{}
}

func New(logger *log.Logger, a ...Aggregator) *Server {
	return &Server{
		logger:      logger,
		Aggregators: a,
		stopChan:    make(chan struct{}, 1),
	}
}

func (s *Server) Run() {
	entryChan := make(chan *entry.Entry)
	for _, a := range s.Aggregators {
		go func() {
			a.Aggregate(entryChan)
		}()
	}

	s.logger.Printf("aggregator : running")

	for {
		select {
		case entry := <-entryChan:
			s.logger.Printf("aggregator : received entry %v", entry)
		case <-s.stopChan:
			break
		}
	}
}

func (s *Server) Shutdown() {
	s.stopChan <- struct{}{}
}
