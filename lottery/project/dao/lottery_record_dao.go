package dao

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"project/lottery/project/models"
)

type LotteryRecord struct {
	engine *xorm.Engine
}

func NewResultDao(engine *xorm.Engine) *LotteryRecord {
	return &LotteryRecord{
		engine: engine,
	}
}

func (d *LotteryRecord) Get(id int) *models.TbLotteryRecord {
	data := &models.TbLotteryRecord{Id: id}
	ok, err := d.engine.Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *LotteryRecord) GetAll(page, size int) []models.TbLotteryRecord {
	offset := (page - 1) * size
	datalist := make([]models.TbLotteryRecord, 0)
	err := d.engine.Desc("id").Limit(size, offset).Find(&datalist)
	if err != nil {
		fmt.Print("LotteryRecord.GetAll error:", err)
	}
	return datalist
}

func (d *LotteryRecord) CountAll() int64 {
	num, err := d.engine.Count(&models.TbLotteryRecord{})
	if err != nil {
		return 0
	} else {
		return num
	}
}

func (d *LotteryRecord) GetNewPrize(size int, giftIds []int) []models.TbLotteryRecord {
	datalist := make([]models.TbLotteryRecord, 0)
	err := d.engine.In("gift_id", giftIds).Desc("id").Limit(size).Find(&datalist)
	if err != nil {
		fmt.Print("LotteryRecord.GetNewPrize error:", err)
	}
	return datalist
}

func (d *LotteryRecord) SearchByGift(giftId, size, page int) []models.TbLotteryRecord {
	datalist := make([]models.TbLotteryRecord, 0)
	offset := (page - 1) * size
	err := d.engine.Where("gift_id=?", giftId).Desc("id").Limit(size, offset).Find(&datalist)
	if err != nil {
		fmt.Print("LotteryRecord.SearchByGift error:", err)
	}
	return datalist
}

func (d *LotteryRecord) SearchByUser(uid, size, page int) []models.TbLotteryRecord {
	datalist := make([]models.TbLotteryRecord, 0)
	offset := (page - 1) * size
	err := d.engine.Where("uid=?", uid).Desc("id").Limit(size, offset).Find(&datalist)
	if err != nil {
		fmt.Print("LotteryRecord.SearchByUser error:", err)
	}
	return datalist
}

func (d *LotteryRecord) CountByGift(giftId int) int64 {
	num, err := d.engine.Where("gift_id=?", giftId).Count(&models.TbLotteryRecord{})
	if err != nil {
		return 0
	} else {
		return num
	}
}

func (d *LotteryRecord) CountByUser(uid int) int64 {
	num, err := d.engine.Where("uid=?", uid).Count(&models.TbLotteryRecord{})
	if err != nil {
		return 0
	} else {
		return num
	}
}

func (d *LotteryRecord) Delete(id int) error {
	data := &models.TbLotteryRecord{Id: id, SysStatus: 1}
	_, err := d.engine.ID(id).Update(data)
	return err
}

func (d *LotteryRecord) Update(data *models.TbLotteryRecord, columns []string) error {
	_, err := d.engine.ID(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *LotteryRecord) Create(data *models.TbLotteryRecord) error {
	_, err := d.engine.Insert(data)
	return err
}

