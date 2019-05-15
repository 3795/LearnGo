package dao

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"project/lottery/project/models"
)

type UserTimeDao struct {
	engine *xorm.Engine
}

func NewUserdayDao(engine *xorm.Engine) *UserTimeDao {
	return &UserTimeDao{
		engine: engine,
	}
}

func (d *UserTimeDao) Get(id int) *models.TbUserTime {
	data := &models.TbUserTime{Id: id}
	ok, err := d.engine.Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *UserTimeDao) GetAll(page, size int) []models.TbUserTime {
	offset := (page - 1) * size
	datalist := make([]models.TbUserTime, 0)
	err := d.engine.Desc("id").Limit(size, offset).Find(&datalist)
	if err != nil {
		fmt.Print("UserTimeDao.GetAll error:", err)
	}
	return datalist
}

func (d *UserTimeDao) CountAll() int64 {
	num, err := d.engine.Count(&models.TbUserTime{})
	if err != nil {
		return 0
	} else {
		return num
	}
}

func (d *UserTimeDao) Search(uid, day int) []models.TbUserTime {
	datalist := make([]models.TbUserTime, 0)
	err := d.engine.Where("uid=?", uid).Where("day=?", day).Desc("id").Find(&datalist)
	if err != nil {
		fmt.Print("UserTimeDao.Search error:", err)
	}
	return datalist
}

func (d *UserTimeDao) Count(uid, day int) int {
	info := &models.TbUserTime{}
	ok, err := d.engine.Where("uid=?", uid).Where("day=?", day).Get(info)
	if !ok || err != nil {
		return 0
	} else {
		return info.Num
	}
}

func (d *UserTimeDao) Update(data *models.TbUserTime, columns []string) error {
	_, err := d.engine.ID(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *UserTimeDao) Create(data *models.TbUserTime) error {
	_, err := d.engine.Insert(data)
	return err
}
