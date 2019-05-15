package service

import (
	"project/lottery/project/dao"
	"project/lottery/project/datasource"
	"project/lottery/project/models"
)

type LotteryRecordService interface {
	GetAll(page, size int) []models.TbLotteryRecord
	CountAll() int64
	GetNewPrize(size int, giftIds []int) []models.TbLotteryRecord
	SearchByGift(giftId, page, size int) []models.TbLotteryRecord
	SearchByUser(uid, page, size int) []models.TbLotteryRecord
	CountByGift(giftId int) int64
	CountByUser(uid int) int64
	Get(id int) *models.TbLotteryRecord
	Delete(id int) error
	Update(user *models.TbLotteryRecord, columns []string) error
	Create(user *models.TbLotteryRecord) error
}

type resultService struct {
	dao *dao.LotteryRecord
}

func NewResultService() LotteryRecordService {
	return &resultService{
		dao: dao.NewResultDao(datasource.InstanceDBMaster()),
	}
}

func (s *resultService) GetAll(page, size int) []models.TbLotteryRecord {
	return s.dao.GetAll(page, size)
}

func (s *resultService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *resultService) GetNewPrize(size int, giftIds []int) []models.TbLotteryRecord {
	return s.dao.GetNewPrize(size, giftIds)
}

func (s *resultService) SearchByGift(giftId, page, size int) []models.TbLotteryRecord {
	return s.dao.SearchByGift(giftId, page, size)
}

func (s *resultService) SearchByUser(uid, page, size int) []models.TbLotteryRecord {
	return s.dao.SearchByUser(uid, page, size)
}

func (s *resultService) CountByGift(giftId int) int64 {
	return s.dao.CountByGift(giftId)
}

func (s *resultService) CountByUser(uid int) int64 {
	return s.dao.CountByUser(uid)
}

func (s *resultService) Get(id int) *models.TbLotteryRecord {
	return s.dao.Get(id)
}

func (s *resultService) Delete(id int) error {
	return s.dao.Delete(id)
}

func (s *resultService) Update(data *models.TbLotteryRecord, columns []string) error {
	return s.dao.Update(data, columns)
}

func (s *resultService) Create(data *models.TbLotteryRecord) error {
	return s.dao.Create(data)
}
