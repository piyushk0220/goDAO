package repository

import (
	"context"
	"src/DaoInterface/model"
)

type EmpRepo interface {
	Fetch(ctx context.Context, num int64) ([]*model.Employee, error)
	GetByID(ctx context.Context, id int64) (*model.Employee, error)
	Create(ctx context.Context, e *model.Employee) (int64, error)
	Update(ctx context.Context, e *model.Employee) (*model.Employee, error)
	Delete(ctx context.Context, id int64) (bool, error)
}
