package repository

import (
	"github.com/jmoiron/sqlx"
	"petstore/internal/model"
)

type PetRepository struct {
	db *sqlx.DB
}

func NewPetRepository(db *sqlx.DB) *PetRepository {
	return &PetRepository{db: db}
}

func (r *PetRepository) Create(pet *model.Pet) error {
	query := `
		INSERT INTO pets (name, status, category_id)
		VALUES ($1, $2, $3)
		RETURNING id`

	return r.db.QueryRow(query, pet.Name, pet.Status, pet.CategoryID).Scan(&pet.ID)
}

func (r *PetRepository) GetByID(id int64) (*model.Pet, error) {
	pet := &model.Pet{}
	query := `SELECT id, name, status, category_id FROM pets WHERE id = $1`

	err := r.db.Get(pet, query, id)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

func (r *PetRepository) Update(pet *model.Pet) error {
	query := `
		UPDATE pets 
		SET name = $1, status = $2, category_id = $3
		WHERE id = $4`

	_, err := r.db.Exec(query, pet.Name, pet.Status, pet.CategoryID, pet.ID)
	return err
}

func (r *PetRepository) Delete(id int64) error {
	query := `DELETE FROM pets WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
