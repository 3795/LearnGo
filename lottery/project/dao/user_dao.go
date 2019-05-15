package dao

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"project/lottery/project/models"
)

type UserDao struct {
	engine *xorm.Engine
}

func NewUserDao(engine *xorm.Engine) *UserDao {
	return &UserDao{
		engine: engine,
	}
}

func (d *UserDao) Get(id int) *models.TbUserBlacklist {
	data := &models.TbUserBlacklist{Id: id}
	ok, err := d.engine.Get(&data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *UserDao) GetAll(page, size int) []models.TbUserBlacklist {
	offset := (page - 1) * size
	datalist := make([]models.TbUserBlacklist, 0)
	err := d.engine.Desc("id").Limit(size, offset).Find(&datalist)
	if err != nil {
		fmt.Print("UserDao.GetAll error:", err)
	}
	return datalist
}

func (d *UserDao) CountAll() int64 {
	num, err := d.engine.Count(&models.TbUserBlacklist{})
	if err != nil {
		return 0
	} else {
		return num
	}
}

func (d *UserDao) Update(data *models.TbUserBlacklist, columns []string) error {
	_, err := d.engine.ID(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *UserDao) Create(data *models.TbUserBlacklist) error {
	_, err := d.engine.Insert(data)
	return err
}
