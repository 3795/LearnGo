package errors

import (
	"bytes"
	"fmt"
	"strings"
)

type ErrorType string

// 错误类型
const (
	// 下载器错误
	ERROR_TYPE_DOWNLOADER ErrorType = "downloader error"
	// 分析器错误
	ERROR_TYPE_ANALYZER ErrorType = "analyzer error"
	// 条目处理管道错误
	ERROR_TYPE_PIPELINE ErrorType = "pipeline error"
	// 调度器错误
	ERROR_TYPE_SCHEDULER ErrorType = "scheduler error"
)

type CrawlerError interface {
	// 错误类型
	Type() ErrorType
	// 错误提示信息
	Error() string
}

type myCrawlerError struct {
	// 错误类型
	errType ErrorType
	// 简要错误信息
	errMsg string
	// 完整的错误信息
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

// 生成错误提示信息
func (ce *myCrawlerError) genFullErrMsg() {
	var buffer bytes.Buffer
	buffer.WriteString("crawler error: ")
	if ce.errType != "" {
		buffer.WriteString(string(ce.errType))
		buffer.WriteString(": ")
	}
	buffer.WriteString(ce.errMsg)
	ce.fullErrMsg = fmt.Sprintf("%s", buffer.String())
}

// 非法参数错误类型
type IllegalParameterError struct {
	msg string
}

func NewIllegalParameterError(errMsg string) IllegalParameterError {
	return IllegalParameterError{msg: fmt.Sprintf("illegal parameter: %s", strings.TrimSpace(errMsg))}
}

func (ipe IllegalParameterError) Error() string {
	return ipe.msg
}
