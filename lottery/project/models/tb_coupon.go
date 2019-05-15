package models

type TbCoupon struct {
	Id         int    `xorm:"not null pk autoincr INT(11)"`
	GiftId     int    `xorm:"comment('奖品ID，关联lt_gift表') INT(11)"`
	Code       string `xorm:"comment('虚拟券编码') VARCHAR(255)"`
	SysCreated int    `xorm:"comment('创建时间') INT(11)"`
	SysUpdated string `xorm:"comment('更新时间') VARCHAR(255)"`
	SysStatus  int    `xorm:"comment('状态，0：正常，1：作废，2：已发放') TINYINT(3)"`
}
