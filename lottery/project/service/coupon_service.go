package service

import (
	"project/lottery/project/dao"
	"project/lottery/project/datasource"
	"project/lottery/project/models"
)

type CouponService interface {
	GetAll(page, size int) []models.TbCoupon
	CountAll() int64
	CountByGift(giftId int) int64
	Search(giftId int) []models.TbCoupon
	Get(id int) *models.TbCoupon
	Delete(id int) error
	Update(user *models.TbCoupon, columns []string) error
	Create(user *models.TbCoupon) error
	NextUsingCode(giftId, codeId int) *models.TbCoupon
	UpdateByCode(data *models.TbCoupon, columns []string) error
}

type codeService struct {
	dao *dao.CouponDao
}

func NewCodeService() CouponService {
	return &codeService{
		dao: dao.NewCouponDao(datasource.InstanceDBMaster()),
	}
}

func (s *codeService) GetAll(page, size int) []models.TbCoupon {
	return s.dao.GetAll(page, size)
}

func (s *codeService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *codeService) CountByGift(giftId int) int64 {
	return s.dao.CountByGift(giftId)
}

func (s *codeService) Search(giftId int) []models.TbCoupon {
	return s.dao.Search(giftId)
}

func (s *codeService) Get(id int) *models.TbCoupon {
	return s.dao.Get(id)
}

func (s *codeService) Delete(id int) error {
	return s.dao.Delete(id)
}

func (s *codeService) Update(data *models.TbCoupon, columns []string) error {
	return s.dao.Update(data, columns)
}

func (s *codeService) Create(data *models.TbCoupon) error {
	return s.dao.Create(data)
}

func (s *codeService) NextUsingCode(giftId, codeId int) *models.TbCoupon {
	return s.dao.NextUsingCode(giftId, codeId)
}

func (s *codeService) UpdateByCode(data *models.TbCoupon, columns []string) error {
	return s.dao.UpdateByCode(data, columns)
}
