package dao

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"project/lottery/project/models"
)

type CouponDao struct {
	engine *xorm.Engine
}

func NewCouponDao(engine *xorm.Engine) *CouponDao {
	return &CouponDao{
		engine: engine,
	}
}

func (d *CouponDao) Get(id int) *models.TbCoupon {
	data := &models.TbCoupon{Id: id}
	ok, err := d.engine.Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *CouponDao) GetAll(page, size int) []models.TbCoupon {
	offset := (page - 1) * size
	datalist := make([]models.TbCoupon, 0)
	err := d.engine.Desc("id").Limit(size, offset).Find(&datalist)
	if err != nil {
		fmt.Print("CouponDao.GetAll error:", err)
	}
	return datalist
}

func (d *CouponDao) CountAll() int64 {
	num, err := d.engine.Count(&models.TbCoupon{})
	if err != nil {
		fmt.Print("CouponDao.CountAll error:", err)
	}
	return num
}

func (d *CouponDao) CountByGift(giftId int) int64 {
	num, err := d.engine.Where("gift_id=?", giftId).Count(&models.TbCoupon{})
	if err != nil {
		return 0
	} else {
		return num
	}
}

func (d *CouponDao) Search(giftId int) []models.TbCoupon {
	datalist := make([]models.TbCoupon, 0)
	err := d.engine.Where("gift_id=?", giftId).Desc("id").Find(&models.TbCoupon{})
	if err != nil {
		fmt.Print("CouponDao.Search error:", err)
	}
	return datalist
}

func (d *CouponDao) Delete(id int) error {
	data := &models.TbCoupon{Id: id, SysStatus: 1}
	_, err := d.engine.ID(data.Id).Update(data)
	return err
}

func (d *CouponDao) Update(data *models.TbCoupon, columns []string) error {
	_, err := d.engine.ID(data.Id).Cols(columns...).Update(data)
	return err
}

func (d *CouponDao) Create(data *models.TbCoupon) error {
	_, err := d.engine.Insert(data)
	return err
}

// 找到下一个可用的最小的优惠券
func (d *CouponDao) NextUsingCode(giftId, codeId int) *models.TbCoupon {
	datalist := make([]models.TbCoupon, 0)
	err := d.engine.
		Where("gift_id=?", giftId).
		Where("sys_status=?", 0).
		Where("id>?", codeId).
		Asc("id").
		Limit(1).
		Find(&datalist)

	if err != nil || len(datalist) < 1 {
		return nil
	} else {
		return &datalist[0]
	}
}

func (d *CouponDao) UpdateByCode(data *models.TbCoupon, columns []string) error {
	_, err := d.engine.Where("code=?", data.Code).MustCols(columns...).Update(data)
	return err
}
