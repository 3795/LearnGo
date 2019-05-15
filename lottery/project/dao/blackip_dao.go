package dao

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"project/lottery/project/models"
)

type BlackIpDao struct {
	engine *xorm.Engine
}

func NewBlackIpDao(engine *xorm.Engine) *BlackIpDao {
	return &BlackIpDao{
		engine: engine,
	}
}

func (d *BlackIpDao) Get(id int) *models.TbIpBlacklist {
	data := &models.TbIpBlacklist{Id:id}
	ok, err := d.engine.Get(data)
	if err == nil && ok {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *BlackIpDao) GetAll(page, size int) []models.TbIpBlacklist {
	offset := (page - 1) * size
	datalist := make([]models.TbIpBlacklist, 0)
	err := d.engine.Desc("id").Limit(size, offset).Find(&datalist)
	if err != nil {
		fmt.Print("BlackIpDao.GetAll error:", err)
	}
	return datalist
}

func (d *BlackIpDao) CountAll() int64 {
	num, err := d.engine.Count(&models.TbIpBlacklist{})
	if err != nil {
		return 0
	} else {
		return num
	}
}

func (d *BlackIpDao) Delete(id int) error {
	data := &models.TbIpBlacklist{Id:id}
	_, err := d.engine.Id(data.Id).Delete(data)
	return err
}

func (d *BlackIpDao) Update(data *models.TbIpBlacklist, columns []string) error {
	_, err := d.engine.Id(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *BlackIpDao) Create(data *models.TbIpBlacklist) error {
	_, err := d.engine.Insert(data)
	return err
}

func (d *BlackIpDao) GetByIp(ip string) *models.TbIpBlacklist {
	dataList := make([]models.TbIpBlacklist, 0)
	err := d.engine.Where("ip=?", ip).Desc("id").Limit(1).Find(&dataList)
	if err != nil || len(dataList) < 1 {
		return nil
	}
	return &dataList[0]
}

func (d *BlackIpDao) Search(ip string) []models.TbIpBlacklist {
	datalist := make([]models.TbIpBlacklist, 0)
	err := d.engine.Where("ip=?", ip).Desc("id").Find(&datalist)
	if err != nil {
		fmt.Print("BlackIpDao.Search:", err)
	}
	return datalist
}