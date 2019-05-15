package service

import (
	"fmt"
	"project/lottery/project/dao"
	"project/lottery/project/datasource"
	"project/lottery/project/models"
	"strconv"
	"time"
)

type UserdayService interface {
	GetAll(page, size int) []models.TbUserTime
	CountAll() int64
	Search(uid, day int) []models.TbUserTime
	Count(uid, day int) int
	Get(id int) *models.TbUserTime
	//Delete(id int) error
	Update(user *models.TbUserTime, columns []string) error
	Create(user *models.TbUserTime) error
	GetUserToday(uid int) *models.TbUserTime
}

type userdayService struct {
	dao *dao.UserTimeDao
}

func NewUserdayService() UserdayService {
	return &userdayService{
		dao: dao.NewUserdayDao(datasource.InstanceDBMaster()),
	}
}

func (s *userdayService) GetAll(page, size int) []models.TbUserTime {
	return s.dao.GetAll(page, size)
}

func (s *userdayService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *userdayService) Search(uid, day int) []models.TbUserTime {
	return s.dao.Search(uid, day)
}

func (s *userdayService) Count(uid, day int) int {
	return s.dao.Count(uid, day)
}

func (s *userdayService) Get(id int) *models.TbUserTime {
	return s.dao.Get(id)
}

func (s *userdayService) Update(data *models.TbUserTime, columns []string) error {
	return s.dao.Update(data, columns)
}

func (s *userdayService) Create(data *models.TbUserTime) error {
	return s.dao.Create(data)
}

func (s *userdayService) GetUserToday(uid int) *models.TbUserTime {
	y, m, d := time.Now().Date()
	strDay := fmt.Sprintf("%d%02d%02d", y, m, d)
	day, _ := strconv.Atoi(strDay)
	list := s.dao.Search(uid, day)
	if list != nil && len(list) > 0 {
		return &list[0]
	} else {
		return nil
	}
}
