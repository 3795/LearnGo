package dao

import (
	"github.com/go-xorm/xorm"
	"log"
	"project/lottery/project/models"
)

type PrizeDao struct {
	engine *xorm.Engine
}

func NewPrizeDao(engine *xorm.Engine) *PrizeDao {
	return &PrizeDao{
		engine: engine,
	}
}

func (d *PrizeDao) Get(id int) *models.TbPrize {
	data := &models.TbPrize{Id:id}
	ok, err := d.engine.Get(data)
	if err == nil && ok {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *PrizeDao) GetAll() []models.TbPrize {
	dataList := make([]models.TbPrize, 0)
	err := d.engine.Asc("sys_status").Asc("displayorder").Find(&dataList)
	if err != nil {
		log.Println("prize_dao.GetAll error = ", err.Error())
		return dataList
	}
	return dataList
}

func (d *PrizeDao) CountAll() int64 {
	num, err := d.engine.Count(&models.TbPrize{})
	if err != nil {
		return 0
	} else {
		return num
	}
}

func (d *PrizeDao) Delete(id int) error {
	data := &models.TbPrize{Id:id, SysStatus:1}
	_, err := d.engine.Id(data.Id).Update(data)
	return err
}

func (d *PrizeDao) Update(data *models.TbPrize, columns []string) error {
	_, err := d.engine.Id(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *PrizeDao) Create(data *models.TbPrize) error {
	_, err := d.engine.Insert(data)
	return err
}