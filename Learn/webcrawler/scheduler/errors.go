package scheduler

import (
	"Project/LearnGo/Learn/webcrawler/errors"
)

// 生成爬虫错误信息
func genError(errMsg string) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_SCHEDULER, errMsg)
}

func genErrorByError(err error) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_SCHEDULER, err.Error())
}

func genParameterError(errMsg string) error {
	return errors.NewCrawlerErrorBy(errors.ERROR_TYPE_SCHEDULER, errors.NewIllegalParameterError(errMsg))
}

// todo 向缓冲池中发送错误信息
//func sendError(err error, mid module.MID, error)
