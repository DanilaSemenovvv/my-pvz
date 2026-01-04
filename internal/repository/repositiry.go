package repository

import (
	"encoding/json"
	"os"

	"github.com/DanilaSemenovvv/my-pvz/internal/models"
)

const (
	invalidIndex    = -1
	filePermissions = 0644
)

type Repository interface {
	FindAll() ([]models.Order, error)
	SaveAll([]models.Order) error
}

type FileRepository struct {
	filename string
}

func NewFileRepository(filename string) Repository {
	return &FileRepository{filename: filename}
}

func (r *FileRepository) FindAll() ([]models.Order, error) {
	data, err := os.ReadFile(r.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Order{}, nil
		}
		return nil, err
	}

	var orderers []models.Order
	json.Unmarshal(data, &orderers)
	return orderers, err
}

func (r *FileRepository) SaveAll(orders []models.Order) error {
	data, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filename, data, filePermissions)

}
