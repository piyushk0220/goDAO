package repository

import (
	"context"
	"fmt"
	"src/DaoInterface/model"

	"database/sql"
)

type EmpRepo interface {
	Fetch(ctx context.Context) ([]model.Employee, error)
	GetByID(ctx context.Context, id int64) (*model.Employee, error)
	Create(ctx context.Context, e *model.Employee) (int64, error)
	Update(ctx context.Context, e *model.Employee) (*model.Employee, error)
	Delete(ctx context.Context, id int64) (bool, error)
}

// NewSQLPostRepo retunrs implement of EmpRepo repository interface
func NewSQLEmpRepo(Conn *sql.DB) *mysqlEmpRepo {
	return &mysqlEmpRepo{
		Conn: Conn,
	}
}

type mysqlEmpRepo struct {
	Conn *sql.DB
}

func (m *mysqlEmpRepo) fetch(ctx context.Context, query string, args ...interface{}) ([]model.Employee, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payload := []model.Employee{}
	for rows.Next() {
		data := model.Employee{}

		err := rows.Scan(
			&data.Id,
			&data.Name,
			&data.Skill,
		)
		if err != nil {
			fmt.Printf("%#v....\n", data)
			fmt.Printf("%#v....\n", err)
			return nil, err
		}
		payload = append(payload, data)
	}
	fmt.Printf("%#v\n", payload)
	return payload, nil
}

func (m *mysqlEmpRepo) Fetch(ctx context.Context) ([]model.Employee, error) {
	query := "Select id, name, skill From employees"

	return m.fetch(ctx, query)
}
func (m *mysqlEmpRepo) GetByID(ctx context.Context, id int64) (*model.Employee, error) {
	query := "Select id, name, skill From employees where id=?"

	rows, err := m.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	payload := model.Employee{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, model.ErrNotFound
	}

	return &payload, nil
}

func (m *mysqlEmpRepo) Create(ctx context.Context, p *model.Employee) (int64, error) {
	query := "Insert employees SET id=?, name=?, skill=?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}

	res, err := stmt.ExecContext(ctx, p.Id, p.Name, p.Skill)
	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func (m *mysqlEmpRepo) Update(ctx context.Context, p *model.Employee) (*model.Employee, error) {
	query := "Update employees set name=?, skill=? where id=?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(
		ctx,

		p.Name,
		p.Skill,
		p.Id,
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return p, nil
}

func (m *mysqlEmpRepo) Delete(ctx context.Context, id int64) (bool, error) {
	query := "Delete From employees Where id=?"
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}
