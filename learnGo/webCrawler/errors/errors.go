package errors

import (
	"bytes"
	"fmt"
	"strings"
)

type ErrorType string

const (
	// 下载器错误
	ERROR_TYPE_DOWNLOADER ErrorType = "Downloader Error"
	// 分析器错误
	ERROR_TYPE_ANALYZER ErrorType = "Analyzer Error"
	// 条目处理管道错误
	ERROR_TYPE_PIPELINE ErrorType = "Pipeline Error"
	// 调度器错误
	ERROR_TYPE_SCHEDULER ErrorType = "Scheduler Error"
)

// 爬虫错误的接口类型
type CrawlerError interface {
	// 错误类型
	Type() ErrorType
	// 错误提示信息
	Error() string
}

type myCrawlerError struct {
	errType ErrorType
	errMsg string
	fullErrMsg string
}

func NewCrawlerError(errType ErrorType, errMsg string) CrawlerError {
	return &myCrawlerError{
		errType: errType,
		errMsg:  strings.TrimSpace(errMsg),
	}
}

func NewCrawlerErrorBy(errType ErrorType, err error) CrawlerError {
	return NewCrawlerError(errType, err.Error())
}

func (ce *myCrawlerError) Type() ErrorType {
	return ce.errType
}

func (ce *myCrawlerError) Error() string {
	if ce.fullErrMsg == "" {
		ce.genFullErrMsg()
	}
	return ce.fullErrMsg
}

func (ce *myCrawlerError) genFullErrMsg() {
	var buffer bytes.Buffer
	buffer.WriteString("crawler error:")
	if ce.Type() != "" {
		buffer.WriteString(string(ce.Type()))
		buffer.WriteString(":")
	}
	buffer.WriteString(ce.errMsg)
	ce.fullErrMsg = fmt.Sprintf("%s", buffer.String())
}

// IllegalParameterError 代表非法的参数的错误类型。
type IllegalParameterError struct {
	msg string
}

// NewIllegalParameterError 会创建一个IllegalParameterError类型的实例。
func NewIllegalParameterError(errMsg string) IllegalParameterError {
	return IllegalParameterError{
		msg: fmt.Sprintf("illegal parameter: %s",
			strings.TrimSpace(errMsg)),
	}
}

func (ipe IllegalParameterError) Error() string {
	return ipe.msg
}