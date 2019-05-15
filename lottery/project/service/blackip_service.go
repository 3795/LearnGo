package service

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"project/lottery/project/common"
	"project/lottery/project/dao"
	"project/lottery/project/datasource"
	"project/lottery/project/models"
	"sync"
)

var cachedBlackipList = make(map[string]*models.TbIpBlacklist)
var cachedBlackipLock = sync.Mutex{}

type BlackipService interface {
	Get(id int) *models.TbIpBlacklist
	GetByIp(ip string) *models.TbIpBlacklist
	GetAll(page, size int) []models.TbIpBlacklist
	CountAll() int64
	Search(ip string) []models.TbIpBlacklist
	Update(user *models.TbIpBlacklist, columns []string) error
	Create(user *models.TbIpBlacklist) error
}

type blackipService struct {
	dao *dao.BlackIpDao
}

func NewBlackipService() BlackipService {
	return &blackipService{
		dao: dao.NewBlackIpDao(datasource.InstanceDBMaster()),
	}
}

func (s *blackipService) Get(id int) *models.TbIpBlacklist {
	return s.dao.Get(id)
}

func (s *blackipService) GetByIp(ip string) *models.TbIpBlacklist {
	data := s.getByCache(ip)
	if data == nil || data.Ip == "" {
		data = s.dao.GetByIp(ip)
		if data == nil || data.Ip == "" {
			data = &models.TbIpBlacklist{Ip: ip}
		}
		s.setByCache(data)
	}
	return data
}

func (s *blackipService) GetAll(page, size int) []models.TbIpBlacklist {
	return s.dao.GetAll(page, size)
}

func (s *blackipService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *blackipService) Search(ip string) []models.TbIpBlacklist {
	return s.dao.Search(ip)
}

func (s *blackipService) Update(user *models.TbIpBlacklist, columns []string) error {
	s.updateByCache(user, columns)
	return s.dao.Update(user, columns)
}

func (s *blackipService) Create(user *models.TbIpBlacklist) error {
	return s.dao.Create(user)
}

func (s *blackipService) getByCache(ip string) *models.TbIpBlacklist {
	key := fmt.Sprintf("info_blackip_%s", ip)
	rds := datasource.InstanceCache()
	dataMap, err := redis.StringMap(rds.Do("HGETALL", key))
	if err != nil {
		log.Println("blockip_service.getByCache HGETALL key=", key, ", error=", err)
		return nil
	}
	dataIp := common.GetStringFromStringMap(dataMap, "Ip", "")
	if dataIp == "" {
		return nil
	}
	data := &models.TbIpBlacklist{
		Id:         int(common.GetInt64FromStringMap(dataMap, "Id", 0)),
		Ip:         dataIp,
		Blacktime:  int(common.GetInt64FromStringMap(dataMap, "Blacktime", 0)),
		SysCreated: int(common.GetInt64FromStringMap(dataMap, "SysCreated", 0)),
		SysUpdated: int(common.GetInt64FromStringMap(dataMap, "SysUpdated", 0)),
	}
	return data
}

func (s *blackipService) setByCache(data *models.TbIpBlacklist) {
	if data == nil || data.Ip == "" {
		return
	}

	key := fmt.Sprintf("info_blackip_%s", data.Ip)
	rds := datasource.InstanceCache()

	params := []interface{}{key}
	params = append(params, "Ip", data.Ip)
	if data.Id > 0 {
		params = append(params, "Blacktime", data.Blacktime)
		params = append(params, "SysCreated", data.SysCreated)
		params = append(params, "SysUpdated", data.SysUpdated)
	}

	_, err := rds.Do("HMSET", params...)
	if err != nil {
		log.Println("blackip_service.setByCache HMSET params=", params, ", error=", err)
	}
}

func (s *blackipService) updateByCache(data *models.TbIpBlacklist, columns []string) {
	if data == nil || data.Ip == "" {
		return
	}

	key := fmt.Sprintf("info_blackip_%s", data.Ip)
	rds := datasource.InstanceCache()
	_, _ = rds.Do("DEL", key)
}

