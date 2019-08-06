package main

import (
	"Project/LearnGo/Learn/webcrawler/examples/finder/internal"
	"Project/LearnGo/Learn/webcrawler/examples/finder/monitor"
	"Project/LearnGo/Learn/webcrawler/helper/log"
	"Project/LearnGo/Learn/webcrawler/scheduler"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	firstURL string
	domains  string
	depth    uint
	dirPath  string
)

var logger = log.DLogger()

func init() {
	flag.StringVar(&firstURL, "first", "http://zhihu.sogou.com/zhihu?query=golang+logo", "The first URL which you want to access.")
	flag.StringVar(&domains, "domains", "zhihu.com", "The primary domains which you accepted. Please using comma-separated multiple domains.")
	flag.UintVar(&depth, "depth", 3, "The depth for crawling.")
	flag.StringVar(&dirPath, "dir", "./pictures", "The path which you want to save the image files.")
}

func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	_, _ = fmt.Fprintf(os.Stderr, "\tfinder [flags] \n")
	_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = Usage
	flag.Parse()
	// 调度器初始化参数
	sl := scheduler.NewScheduler()
	domainParts := strings.Split(domains, ",")
	var acceptedDomains []string
	for _, domain := range domainParts {
		domain = strings.TrimSpace(domain)
		if domain != "" {
			acceptedDomains = append(acceptedDomains, domain)
		}
	}
	requestArgs := scheduler.RequestArgs{
		AcceptedDomains: acceptedDomains,
		MaxDepth:        uint32(depth),
	}
	dataArgs := scheduler.DataArgs{
		ReqBufferCap:         50,
		ReqMaxBufferNumber:   1000,
		RespBufferCap:        50,
		RespMaxBufferNumber:  10,
		ItemBufferCap:        50,
		ItemMaxBufferNumber:  100,
		ErrorBufferCap:       50,
		ErrorMaxBufferNumber: 1,
	}
	downloaders, err := internal.GetDownloaders(1)
	if err != nil {
		logger.Fatalf("An error occurs when creating downloaders: %s", err)
	}

	analyzers, err := internal.GetAnalyzers(1)
	if err != nil {
		logger.Fatalf("An error occurs when creating analyzers: %s", err)
	}
	pipelines, err := internal.GetPipelines(1, dirPath)
	if err != nil {
		logger.Fatalf("An error occur when creating pipelines: %s", err)
	}

	moduleArgs := scheduler.ModuleArgs{
		Downloaders: downloaders,
		Analyzers:   analyzers,
		Pipelines:   pipelines,
	}

	// 初始化调度器
	err = sl.Init(requestArgs, dataArgs, moduleArgs)
	if err != nil {
		logger.Fatalf("An error occurs when initializing scheduler: %s", err)
	}
	// 准备监控参数
	checkInterval := time.Second
	summarizeInterval := 100 * time.Millisecond
	maxIdleCount := uint(5)
	// 开始监控
	checkCountChan := monitor.Monitor(sl, checkInterval, summarizeInterval, maxIdleCount, true, internal.Record)

	// 调度器启动参数
	firstHttpReq, err := http.NewRequest("GET", firstURL, nil)
	if err != nil {
		logger.Fatalln(err)
		return
	}

	// 开启调度器
	err = sl.Start(firstHttpReq)
	if err != nil {
		logger.Fatalf("An error occurs when starting scheduler: %s", err)
	}

	// 等待监控结果
	<-checkCountChan
}
