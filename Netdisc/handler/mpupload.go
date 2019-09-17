package handler

import (
	"LearnGo/Netdisc/cache/redis"
	"LearnGo/Netdisc/util"
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

/*
分块上传的处理器
*/

// 初始化信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

// 分块上传时初始化相关信息
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户的请求参数
	_ = r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		_, _ = w.Write(util.NewRespMsg(-1, "param invalid", nil).JSONBytes())
		return
	}

	redisConn := redis.RedisPool().Get()
	defer redisConn.Close()

	// 2. 生成分块上传的初始化信息
	upInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024,                                       // 分块文件大小为5M
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))), // 需要分多少块
	}

	// 3. 将初始化信息写入Redis中
	_, _ = redisConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkCount", upInfo.ChunkCount)
	_, _ = redisConn.Do("HSET", "MP_"+upInfo.UploadID, "fileHash", upInfo.FileHash)
	_, _ = redisConn.Do("HSET", "MP_"+upInfo.UploadID, "fileSize", upInfo.FileSize)

	// 4. 将响应信息返回给客户端
	_, _ = w.Write(util.NewRespMsg(0, "OK", upInfo).JSONBytes())
}

// 分块上传文件的具体处理方式
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	redisConn := redis.RedisPool().Get()
	defer redisConn.Close()

	// 2. 获取文件句柄，用户存储分块内容
	fPath := "/data/" + uploadID + "/" + chunkIndex
	_ = os.MkdirAll(path.Dir(fPath), 0744) // 创建目标文件夹
	fd, err := os.Create(fPath)            // 创建目标文件
	if err != nil {
		_, _ = w.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes())
		return
	}

	defer fd.Close()

	buf := make([]byte, 1024)
	for {
		n, err := r.Body.Read(buf)
		_, _ = fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 3. 更新Redis中的信息
	_, _ = redisConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 返回处理结果
	_, _ = w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}
