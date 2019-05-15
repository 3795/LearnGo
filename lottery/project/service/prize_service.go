package service

import (
	"project/lottery/project/dao"
	"project/lottery/project/models"
)

type PrizeService interface {
	GetAll() []models.TbPrize
	CountAll() int64
	Get(id int) *models.TbPrize
	Delete(id int) error
	Update(data *models.TbPrize, columns []string) error
	Create(data *models.TbPrize) error
}

type prizeService struct {
	dao *dao.PrizeDao
}

func NewPrizeService() PrizeService {
	return &prizeService{
		dao: dao.NewPrizeDao(nil),
	}
}

func (p *prizeService) GetAll() []models.TbPrize {
	return p.dao.GetAll()
}

func (p *prizeService) CountAll() int64 {
	return p.dao.CountAll()
}

func (p *prizeService) Get(id int) *models.TbPrize {
	return p.dao.Get(id)
}

func (p *prizeService) Delete(id int) error {
	return p.dao.Delete(id)
}

func (p *prizeService) Update(data *models.TbPrize, columns []string) error {
	return p.dao.Update(data, columns)
}

func (p *prizeService) Create(data *models.TbPrize) error {
	return p.dao.Create(data)
}
