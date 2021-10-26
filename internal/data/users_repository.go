package data

import (
	"context"
	"log"
	"time"

	pb "github.com/IgorDevCuemby/crudGrpcMysqlMicroservice/users/userpb"
)

type UsersRepository struct {
	Data *Data
}

func (pr *UsersRepository) GetOne(ctx context.Context, id string) (*pb.CreateUserRequest, error) {
	q := `SELECT id, name, email FROM users WHERE id = $1;`

	row := pr.Data.DB.QueryRowContext(ctx, q, id)
	var usr pb.CreateUserRequest
	err := row.Scan(&usr.Id, &usr.Name, &usr.Email)

	if err != nil {
		log.Printf("Error %v", err.Error())
		return nil, err
	}
	return &usr, nil
}
func (pr *UsersRepository) Verify(ctx context.Context, id uint32) (*pb.User, error) {
	q := `SELECT id, name, email FROM users WHERE id = $1;`

	row := pr.Data.DB.QueryRowContext(ctx, q, id)
	var usr pb.User
	err := row.Scan(&usr.Id, &usr.Name, &usr.Email)

	if err != nil {
		log.Printf("Error %v", err.Error())
		return nil, err
	}
	return &usr, nil
}

func (pr *UsersRepository) GetAll(ctx context.Context) ([]*pb.User, error) {
	q := `SELECT id, name, email, created_at, updated_at FROM users; `

	rows, err := pr.Data.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*pb.User
	for rows.Next() {
		var p pb.User
		rows.Scan(&p.Id, &p.Name, &p.Email, &p.CreatedAt, &p.UpdatedAt)
		users = append(users, &p)
	}
	return users, nil
}

func (pr *UsersRepository) Create(ctx context.Context, p *pb.CreateUserRequest) (*pb.CreateUserRequest, error) {
	q := `
    INSERT INTO Users (name, email, password, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id;
    `
	stmt, err := pr.Data.DB.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	p.Password, _ = EncryptPassword(p.Password)
	row := stmt.QueryRowContext(ctx, p.Name, p.Email, p.Password, time.Now(), time.Now())
	err = row.Scan(&p.Id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (pr UsersRepository) Update(ctx context.Context, id string, p *pb.CreateUserRequest) (*pb.CreateUserRequest, error) {
	q := `UPDATE users set name=$1, email=$2, password=$3, updated_at=$4 WHERE id=$5; `

	stmt, err := pr.Data.DB.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(
		ctx, p.Name, p.Email, p.Password, time.Now(), id,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (pr *UsersRepository) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM Users WHERE id=$1;`

	stmt, err := pr.Data.DB.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
