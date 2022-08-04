package emp

import (
	"context"
	"database/sql"

	model "src/DaoInterface/model"
	eRepo "src/DaoInterface/repository"
)

// NewSQLPostRepo retunrs implement of EmpRepo repository interface
func NewSQLEmpRepo(Conn *sql.DB) eRepo.EmpRepo {
	return &mysqlEmpRepo{
		Conn: Conn,
	}
}

type mysqlEmpRepo struct {
	Conn *sql.DB
}

func (m *mysqlEmpRepo) fetch(ctx context.Context, query string, args ...interface{}) ([]*model.Employee, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payload := make([]*model.Employee, 0)
	for rows.Next() {
		data := new(model.Employee)

		err := rows.Scan(
			&data.Id,
			&data.Name,
			&data.Skill,
		)
		if err != nil {
			return nil, err
		}
		payload = append(payload, data)
	}
	return payload, nil
}

func (m *mysqlEmpRepo) Fetch(ctx context.Context, num int64) ([]*model.Employee, error) {
	query := "Select id, name, skill From employees limit ?"

	return m.fetch(ctx, query, num)
}
func (m *mysqlEmpRepo) GetByID(ctx context.Context, id int64) (*model.Employee, error) {
	query := "Select id, name, skill From employees where id=?"

	rows, err := m.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	payload := &model.Employee{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, model.ErrNotFound
	}

	return payload, nil
}

func (m *mysqlEmpRepo) Create(ctx context.Context, p *model.Employee) (int64, error) {
	query := "Insert employees SET name=?, skill=?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}

	res, err := stmt.ExecContext(ctx, p.Name, p.Skill)
	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func (m *mysqlEmpRepo) Update(ctx context.Context, p *model.Employee) (*model.Employee, error) {
	query := "Update employees set title=?, content=? where id=?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(
		ctx,
		p.Id,
		p.Name,
		p.Skill,
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
