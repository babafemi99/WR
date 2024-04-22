package cockroachDB

import (
	"context"
	"errors"
	"fmt"
	"github.com/babafemi99/WR/internal/logger"
	"github.com/babafemi99/WR/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(dsn string) *Repository {
	pool, err := pgxpool.New(context.TODO(), dsn)
	if err != nil {
		logger.Log.Error(fmt.Errorf("[DAL]: unable to connect: %v", err.Error()).Error())
		os.Exit(1)
	}

	err = pool.Ping(context.TODO())
	if err != nil {
		logger.Log.Error(fmt.Errorf("[DAL]: unable to ping: %v", err.Error()).Error())
		os.Exit(1)
	}
	return &Repository{pool: pool}
}

// UTIL

func (r Repository) Shutdown(ctx context.Context) {
	if r.pool != nil {
		r.pool.Close()
	}
}

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

func (r Repository) PersistWedding(ctx context.Context, req model.NewWeddingReq) error {
	stmt := `INSERT INTO weddings (
                      couple_name, 
                      couple_id, 
                      state, 
                      link, 
                      guest_link, 
                      wedding_id, 
                      wedding_date, 
                      created_at 
                      ) values ($1, $2, $3, $4, $5, $6, $7, $8)`

	exec, err := r.pool.Exec(ctx, stmt,
		req.CoupleName, req.CoupleId, req.State, req.Link, req.GuestLink, req.WeddingId, req.WeddingDate, req.CreatedAt)
	if err != nil {
		return err
	}
	log.Println("EXEC COMMAND ", exec)
	return nil
}

func (r Repository) GetWeddingById(ctx context.Context, id string) (*model.WeddingIdRes, error) {
	stmt := `
		SELECT 
			   couple_name, 
			   state, 
			   location, 
			   screen, 
			   status 
		FROM weddings 
		WHERE wedding_id = $1
`

	var wedding model.WeddingIdRes
	row := r.pool.QueryRow(ctx, stmt, id)
	err := row.Scan(
		&wedding.CoupleName,
		&wedding.State,
		&wedding.Location,
		&wedding.Screen,
		&wedding.Status,
	)
	if err != nil {
		return nil, err
	}

	return &wedding, nil

}

func (r Repository) ToggleWeddingLink(ctx context.Context, req model.ToggleWeddingReq) error {
	stmt := `
		UPDATE weddings 
		SET 
		    status = 'live',
		    screen = $1,
		    location = $2,
		    modified_by = $3,
		    modified_at = $4
		WHERE wedding_id = $5
    `
	exec, err := r.pool.Exec(ctx, stmt, req.Screen, req.Registry, req.TogglerId, req.ModifiedAt, req.WeddingId)
	if err != nil {
		return err
	}
	log.Println("EXEC COMMAND ", exec)
	return nil
}
func (r Repository) ToggleWeddingLinkOff(ctx context.Context, id string) error {
	log.Println("id", id)
	stmt := `
		UPDATE weddings 
		SET 
		    status = 'done'
		WHERE wedding_id = $1
    `
	exec, err := r.pool.Exec(ctx, stmt, id)
	if err != nil {
		return err
	}
	log.Println("EXEC COMMAND ", exec)
	return nil
}

func (r Repository) GetLinkByKey(key string) (model.Wedding, error) { panic("implement me ") }

func (r Repository) IdExists(id string) (bool, error) { panic("implement me") }

func (r Repository) IdToday(id string) (bool, error) { panic("implement me") }

func (r Repository) AddMembers(ctx context.Context, req model.Member) error {
	stmt := `INSERT INTO  wedding_members(wedding_id, member_email, member_code) VALUES ($1, $2, $3 )`

	exec, err := r.pool.Exec(ctx, stmt, req.WeddingId, req.MemberEmail, req.MemberCode)
	if err != nil {
		return err
	}

	log.Println("EXEC COMMAND ", exec)
	return nil
}

func (r Repository) MemberExist(ctx context.Context, req model.Member) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS (SELECT 1 FROM wedding_members WHERE member_email = $1 AND wedding_id = $2) AS exists`

	err := r.pool.QueryRow(ctx, stmt, req.MemberEmail, req.WeddingId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil

}

func (r Repository) MemberCodeExist(ctx context.Context, code, weddingId string) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS (SELECT 1 FROM wedding_members WHERE member_code = $1 AND wedding_id = $2) AS exists `

	err := r.pool.QueryRow(ctx, stmt, code, weddingId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil

}

func (r Repository) GetMembers(ctx context.Context, id string, offset, limit int) ([]model.Member, error) {
	stmt := ` 
		SELECT 
		     member_email, member_code
		    FROM wedding_members 
		    WHERE wedding_id = $1
		OFFSET $2
		LIMIT $3
		    `

	var members []model.Member
	row, err := r.pool.Query(ctx, stmt, id, offset, limit)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var member model.Member
		err = row.Scan(
			&member.MemberEmail,
			&member.MemberCode,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, err

}

func (r Repository) RemoveMember(ctx context.Context, email, id string) error {
	stmt := `DELETE FROM wedding_members WHERE member_email = $1 AND  wedding_id = $2`
	exec, err := r.pool.Exec(ctx, stmt, email, id)
	if err != nil {
		return err
	}

	log.Println("EXEC COMMAND", exec)
	return nil
}

// ADMIN

func (r Repository) EmailExist(ctx context.Context, email, userType string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM "

	switch userType {
	case "admin":
		query += "admin WHERE email = $1) AS exists"
	case "staff":
		query += "staff WHERE email = $1) AS exists"
	default:
		return false, errors.New("invalid user type")
	}

	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r Repository) IdExist(ctx context.Context, id, userType string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM "

	switch userType {
	case "admin":
		query += "admin WHERE id = $1) AS exists"
	case "staff":
		query += "staff WHERE id = $1) AS exists"
	default:
		return false, errors.New("invalid user type")
	}

	err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r Repository) PersistAdmin(ctx context.Context, admin model.Admin) error {
	stmt := `INSERT INTO admin (
                   id, 
                   first_name, 
                   last_name, 
                   email, 
                   hash_password, 
                   role, 
                   status, 
                   created_at, 
                   updated_at
                   )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	exec, err := r.pool.Exec(ctx, stmt, admin.Id, admin.FirstName, admin.LastName, admin.Email, admin.HashPassword, admin.Role,
		admin.Status, admin.CreatedAt, admin.UpdatedAt)
	if err != nil {
		return err
	}

	log.Println("EXEC COMMAND ", exec)
	return nil
}

func (r Repository) PersistStaff(ctx context.Context, staff model.Staff) error {
	stmt := `
		INSERT INTO staff(
                  id,
                  first_name,
                  last_name,
                  email,
                  hash_password,
                  status,
                  state,
                  created_at,
                  updated_at
                  ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	exec, err := r.pool.Exec(ctx, stmt, staff.Id, staff.FirstName, staff.LastName, staff.Email, staff.HashPassword, staff.Status, staff.State, staff.CreatedAt, staff.UpdatedAt)
	if err != nil {
		return err
	}

	log.Println("EXEC COMMAND ", exec)
	return nil
}

func (r Repository) FindAdminByEmail(ctx context.Context, email string) (model.Admin, error) {
	stmt := `
	SELECT id, first_name, last_name, email, hash_password, role, status FROM admin WHERE email = $1
`
	var res model.Admin
	row := r.pool.QueryRow(ctx, stmt, email)
	err := row.Scan(
		&res.Id,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.HashPassword,
		&res.Role,
		&res.Status,
	)
	if err != nil {
		return model.Admin{}, err
	}

	return res, nil
}

func (r Repository) FindSuperAdminByEmail(email string) (model.Admin, error) {
	panic("implement me")
}

func (r Repository) FindStaffByEmail(ctx context.Context, email string) (model.Staff, error) {
	stmt := `
	SELECT id, first_name, last_name, email, hash_password, role, status, state FROM staff WHERE email = $1
`
	var res model.Staff
	row := r.pool.QueryRow(ctx, stmt, email)
	err := row.Scan(
		&res.Id,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.HashPassword,
		&res.Role,
		&res.Status,
		&res.State,
	)
	if err != nil {
		return model.Staff{}, err
	}

	return res, nil
}

func (r Repository) FindAdminById(ctx context.Context, id string) (model.Admin, error) {
	stmt := `
	SELECT id, first_name, last_name, email, hash_password, role, status FROM admin WHERE id = $1
`
	var res model.Admin
	row := r.pool.QueryRow(ctx, stmt, id)
	err := row.Scan(
		&res.Id,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.HashPassword,
		&res.Role,
		&res.Status,
	)
	if err != nil {
		return model.Admin{}, err
	}

	return res, nil
}

func (r Repository) FindStaffById(ctx context.Context, id string) (model.Staff, error) {
	stmt := `
	SELECT id, first_name, last_name, email, hash_password, role, status, state FROM staff WHERE id = $1
`
	var res model.Staff
	row := r.pool.QueryRow(ctx, stmt, id)
	err := row.Scan(
		&res.Id,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.HashPassword,
		&res.Role,
		&res.Status,
		&res.State,
	)
	if err != nil {
		return model.Staff{}, err
	}

	return res, nil
}

func (r Repository) BLockStaff(ctx context.Context, id string) error {
	stmt := `UPDATE staff SET status = 'blocked' WHERE id = $1`
	exec, err := r.pool.Exec(ctx, stmt, id)
	if err != nil {
		return err
	}
	log.Println("EXEC COMMAND ", exec)
	return nil
}

func (r Repository) BLockAdmin(ctx context.Context, id string) error {
	stmt := `UPDATE admin SET status = 'blocked' WHERE id = $1`
	exec, err := r.pool.Exec(ctx, stmt, id)
	if err != nil {
		return err
	}
	log.Println("EXEC COMMAND ", exec)
	return nil
}

func (r Repository) DeleteAdmin(ctx context.Context, id string) error {
	stmt := `DELETE FROM admin WHERE id = $1`
	exec, err := r.pool.Exec(ctx, stmt, id)
	if err != nil {
		return err
	}
	log.Println("EXEC COMMAND ", exec)
	return nil

}

func (r Repository) DeleteStaff(ctx context.Context, id string) error {
	stmt := `DELETE FROM staff WHERE id = $1`
	exec, err := r.pool.Exec(ctx, stmt, id)
	if err != nil {
		return err
	}
	log.Println("EXEC COMMAND ", exec)
	return nil
}

func (r Repository) UpdateStaffPassword(ctx context.Context, id string, password string) error {
	stmt := `
	UPDATE staff 
	SET 
	    hash_password = $1, 
	    updated_at = now() 
	WHERE id = $2`

	exec, err := r.pool.Exec(ctx, stmt, password, id)
	if err != nil {
		return err
	}
	log.Println("EXEC COMMAND ", exec)
	return nil
}
func (r Repository) UpdateAdminPassword(ctx context.Context, id string, password string) error {
	stmt := `
	UPDATE admin 
	SET 
	    hash_password = $1, 
	    updated_at = now() 
	WHERE id = $2`

	exec, err := r.pool.Exec(ctx, stmt, password, id)
	if err != nil {
		return err
	}
	log.Println("EXEC COMMAND ", exec)
	return nil

}
