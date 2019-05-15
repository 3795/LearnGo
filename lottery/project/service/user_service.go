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


// 用户信息，可以缓存(本地或者redis)，有更新的时候，可以直接清除缓存或者根据具体情况更新缓存
var cachedUserList = make(map[int]*models.TbUserBlacklist)
var cachedUserLock = sync.Mutex{}

type UserService interface {
	GetAll(page, size int) []models.TbUserBlacklist
	CountAll() int64
	//Search(country string) []models.TbUserBlacklist
	Get(id int) *models.TbUserBlacklist
	//Delete(id int) error
	Update(user *models.TbUserBlacklist, columns []string) error
	Create(user *models.TbUserBlacklist) error
}

type userService struct {
	dao *dao.UserDao
}

func NewUserService() UserService {
	return &userService{
		dao: dao.NewUserDao(datasource.InstanceDBMaster()),
	}
}

func (s *userService) Get(id int) *models.TbUserBlacklist {
	return nil
}

func (s *userService) GetAll(page, size int) []models.TbUserBlacklist {
	return s.dao.GetAll(page, size)
}

func (s *userService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *userService) Update(user *models.TbUserBlacklist, columns []string) error {
	return nil
}

func (s *userService) Create(user *models.TbUserBlacklist) error {
	return nil
}

func (s *userService) getByCache(id int) *models.TbUserBlacklist {
	key := fmt.Sprintf("info_user_%d", id)
	rds := datasource.InstanceCache()
	dataMap, err := redis.StringMap(rds.Do("HGETALL", key))
	if err != nil {
		log.Println("user_service.getByCache HGETALL key=", key, ", error=", err)
		return nil
	}

	dataId := common.GetInt64FromStringMap(dataMap, "Id", 0)
	if dataId <= 0 {
		return nil
	}

	data := &models.TbUserBlacklist{
		Id:         int(dataId),
		Username:   common.GetStringFromStringMap(dataMap, "Username", ""),
		Blacktime:  int(common.GetInt64FromStringMap(dataMap, "Blacktime", 0)),
		Realname:   common.GetStringFromStringMap(dataMap, "Realname", ""),
		Mobile:     common.GetStringFromStringMap(dataMap, "Mobile", ""),
		Address:    common.GetStringFromStringMap(dataMap, "Address", ""),
		SysCreated: int(common.GetInt64FromStringMap(dataMap, "SysCreated", 0)),
		SysUpdated: int(common.GetInt64FromStringMap(dataMap, "SysUpdated", 0)),
		SysIp:      common.GetStringFromStringMap(dataMap, "SysIp", ""),
	}
	return data
}

func (s *userService) setByCache(data *models.TbUserBlacklist) {
	if data == nil || data.Id <= 0 {
		return
	}
	id := data.Id
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_user_%d", id)
	rds := datasource.InstanceCache()
	// 数据更新到redis缓存
	params := []interface{}{key}
	params = append(params, "Id", id)
	if data.Username != "" {
		params = append(params, "Username", data.Username)
		params = append(params, "Blacktime", data.Blacktime)
		params = append(params, "Realname", data.Realname)
		params = append(params, "Mobile", data.Mobile)
		params = append(params, "Address", data.Address)
		params = append(params, "SysCreated", data.SysCreated)
		params = append(params, "SysUpdated", data.SysUpdated)
		params = append(params, "SysIp", data.SysIp)
	}
	_, err := rds.Do("HMSET", params...)
	if err != nil {
		log.Println("user_service.setByCache HMSET params=", params, ", error=", err)
	}
}

func (s *userService) updateByCache(data *models.TbUserBlacklist, columns []string) {
	if data == nil || data.Id <= 0 {
		return
	}
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_user_%d", data.Id)
	rds := datasource.InstanceCache()
	// 删除redis中的缓存
	_, _ = rds.Do("DEL", key)
}
