package service

import "fmt"

type Service struct {
	db DB
	// cache Cache
}

type DB interface {
	HandleCallback() error
}

func New() *Service {

	return &Service{
		// db: db
	}
}

func (svc *Service) HandleCallback() error {
	fmt.Println("Handling callback in service layer")

	return nil
}
