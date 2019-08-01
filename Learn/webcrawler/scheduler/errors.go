package scheduler

import "Project/LearnGo/Learn/webcrawler/errors"

// 生成爬虫错误信息
func genError(errMsg string) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_SCHEDULER, errMsg)
}

// todo
