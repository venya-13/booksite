package repo

import "fmt"

type (
	MockRepo struct{}

	Config struct{}
)

func NewMockRepo(_ Config) *MockRepo {
	return &MockRepo{}
}

func (r *MockRepo) SaveOrUpdate(user map[string]interface{}) error {
	// There will be a real DB
	fmt.Println("💾 Saving user to DB:", user)
	return nil
}
