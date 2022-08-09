package service

import (
	"context"
	"src/DaoInterface/model"
	"src/DaoInterface/repository"
)

type EmpService struct {
	repo repository.EmpRepo
}

func NewEmpService(r repository.EmpRepo) EmpService {
	return EmpService{repo: r}
}

func (service EmpService) Fetch() ([]model.Employee, error) {
	return service.repo.Fetch(context.Background())
}
func (service EmpService) GetByID(id int64) (*model.Employee, error) {
	return service.repo.GetByID(context.Background(), id)
}
func (service EmpService) Create(e *model.Employee) (int64, error) {
	return service.repo.Create(context.Background(), e)
}
func (service EmpService) Update(e *model.Employee) (*model.Employee, error) {
	return service.repo.Update(context.Background(), e)
}
func (service EmpService) Delete(id int64) (bool, error) {
	return service.repo.Delete(context.Background(), id)
}
