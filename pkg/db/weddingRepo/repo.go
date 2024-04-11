package weddingRepo

import (
	"context"
	"github.com/babafemi99/WR/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// UTIL

func (r Repository) RunInTx(ctx context.Context, fn func() error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	var done bool

	defer func() {
		if !done {
			tx.Rollback(ctx) // Rollback if transaction hasn't been committed
		}
		if p := recover(); p != nil {
			tx.Rollback(ctx) // Rollback if panic occurred
			panic(p)         // Re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback(ctx) // Rollback if error occurred
		} else {
			err = tx.Commit(ctx) // Commit the transaction if no error occurred
		}
	}()

	if err = fn(); err != nil {
		return err
	}

	return nil
}

// WEDDING

func (r Repository) PersistWedding() (string, string, error) {
	panic("implement me ")
}

func (r Repository) ToggleWeddingLink(key string) error {
	panic("implement me ")
}

func (r Repository) GetLinkByKey(key string) (model.Wedding, error) {
	panic("implement me ")
}

func (r Repository) IdExists(id string) (bool, error) { panic("implement me") }

func (r Repository) IdToday(id string) (bool, error) { panic("implement me") }

func (r Repository) AddMembers(req model.Member) error { panic("implement me") }

func (r Repository) GetMembers(id string, offset, limit int) ([]model.Member, error) {
	panic("implement me")
}

func (r Repository) RemoveMember(email, id string) error { panic("implement me") }

// ADMIN

func (r Repository) EmailExist(email, userType string) (bool, error) {
	panic("implement me")
}

func (r Repository) IdExist(id, userType string) (bool, error) {
	panic("implement me")
}

func (r Repository) PersistAdmin(admin model.Admin) error {
	panic("implement me")
}

func (r Repository) PersistStaff(staff model.Staff) error {
	panic("implement me")
}

func (r Repository) FindAdminByEmail(email string) (model.Admin, error) {
	panic("implement me")
}

func (r Repository) FindStaffByEmail(email string) (model.Staff, error) {
	panic("implement me")
}

func (r Repository) FindAdminById(id string) (model.Admin, error) {
	panic("implement me")
}

func (r Repository) FindStaffById(id string) (model.Staff, error) {
	panic("implement me")
}

func (r Repository) BLockStaff(id string) error {
	panic("implement me")
}

func (r Repository) BLockAdmin(id string) error {
	panic("implement me")
}

func (r Repository) DeleteAdmin(id string) error {
	panic("implement me")
}

func (r Repository) DeleteStaff(id string) error {
	panic("implement me")
}

func (r Repository) UpdateStaffPassword(id string, password string) error {
	panic("implement me ")
}
func (r Repository) UpdateAdminPassword(id string, password string) error {
	panic("implement me ")
}
