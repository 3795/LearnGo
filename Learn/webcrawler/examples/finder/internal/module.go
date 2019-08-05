package internal

import "Project/LearnGo/Learn/webcrawler/module"

// 序列号生成器
var snGen = module.NewSNGenerator(1, 0)

// 获取下载器列表
func GetDownloaders(number uint8) ([]module.Downloader, error) {
	return nil, nil
}
