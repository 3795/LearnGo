package models

type TbPrize struct {
	Id           int    `xorm:"not null pk autoincr INT(11)"`
	Title        string `xorm:"comment('奖品名称') VARCHAR(255)"`
	PrizeNum     int    `xorm:"comment('奖品数量') INT(11)"`
	LeftNum      int    `xorm:"comment('奖品剩余数量') INT(11)"`
	PrizeCode    string `xorm:"comment('0-9999表示100%，0-0表示万分之一的中奖概率') VARCHAR(50)"`
	PrizeTime    int    `xorm:"comment('发奖周期，单位为天') INT(11)"`
	Img          string `xorm:"comment('奖品图片') VARCHAR(255)"`
	Displayorder int    `xorm:"comment('位置序号，小的排在前面') INT(11)"`
	Gtype        int    `xorm:"comment('奖品类型，0：虚拟币，1：虚拟券，2：实物小奖，3：实物大奖') INT(11)"`
	Gdata        string `xorm:"comment('扩展数据，如虚拟币数量') VARCHAR(255)"`
	TimeBegin    int    `xorm:"comment('开始时间') INT(11)"`
	TimeEnd      int    `xorm:"comment('结束时间') INT(10)"`
	PrizeData    string `xorm:"comment('发奖计划') MEDIUMTEXT"`
	PrizeBegin   int    `xorm:"comment('发奖周期开始') INT(11)"`
	PrizeEnd     int    `xorm:"comment('发奖周期结束') INT(11)"`
	SysStatus    int    `xorm:"comment('0：正常，1：删除') TINYINT(3)"`
}
