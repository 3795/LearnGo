package master

import (
	"LearnGo/go-crontab/crontab/common"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

// 任务的HTTP接口
type ApiServer struct {
	httpServer *http.Server
}

// 创建单例对象
var (
	G_apiServer *ApiServer
)

// 保存任务接口
// POST job = {"name": "job1", "command": "echo hello", "cronExpr": "*/5 * * * * * *"}
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)

	// 解析Post表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 获取传入的Json
	postJob = req.PostForm.Get("job")
	// 反序列化
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}

	// 保存到etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	// 返回成功的响应
	bytes, _ = common.BuildResponse(0, "成功", oldJob)
	_, _ = resp.Write(bytes)
	return

ERR:
	handleResponse(resp, err)
}

// 删除任务接口
// POST /job/delete
func handleJobDelete(resp http.ResponseWriter, req *http.Request) {
	var (
		err    error
		name   string
		oldJob *common.Job
		bytes  []byte
	)

	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	name = req.PostForm.Get("name")

	if oldJob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "成功", oldJob); err != nil {
		goto ERR
	}

	_, _ = resp.Write(bytes)
	return
ERR:
	handleResponse(resp, err)
}

// 列举所有crontab任务
func handleJobList(resp http.ResponseWriter, req *http.Request) {
	var (
		jobList []*common.Job
		bytes   []byte
		err     error
	)

	// 获取列表
	if jobList, err = G_jobMgr.ListJobs(); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", jobList); err != nil {
		goto ERR
	}
	_, _ = resp.Write(bytes)

	return

ERR:
	handleResponse(resp, err)
}

// 强制杀死某个任务
// POST /job/kill  name=job1
func handleJobKill(resp http.ResponseWriter, req *http.Request) {
	var (
		err   error
		name  string
		bytes []byte
	)

	// 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 要杀死的任务名
	name = req.PostForm.Get("name")

	// 杀死任务
	if err = G_jobMgr.KillJob(name); err != nil {
		goto ERR
	}

	// 正常应答
	if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
		_, _ = resp.Write(bytes)
	}
	return

ERR:
	handleResponse(resp, err)
}

// 查询任务日志
func handleJobLog(resp http.ResponseWriter, req *http.Request) {
	var (
		err        error
		name       string // 任务名字
		skipParam  string // 从第几条开始
		limitParam string // 返回多少条
		skip       int
		limit      int
		logArr     []*common.JobLog
		bytes      []byte
	)

	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	name = req.Form.Get("name")
	skipParam = req.Form.Get("skip")
	limitParam = req.Form.Get("limit")
	if skip, err = strconv.Atoi(skipParam); err != nil {
		skip = 0
	}
	if limit, err = strconv.Atoi(limitParam); err != nil {
		limit = 20
	}

	if logArr, err = G_logMgr.ListLog(name, skip, limit); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", logArr); err == nil {
		_, _ = resp.Write(bytes)
	}

ERR:
	handleResponse(resp, err)

}

// 获取健康的worker节点列表
func handleWorkerList(resp http.ResponseWriter, req *http.Request) {
	var (
		workerArr []string
		err       error
		bytes     []byte
	)

	if workerArr, err = G_workerMgr.ListWorkers(); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", workerArr); err == nil {
		_, _ = resp.Write(bytes)
	}

ERR:
	handleResponse(resp, err)
}

// 处理业务中的异常
func handleResponse(resp http.ResponseWriter, err error) {
	if bytes, err := common.BuildResponse(-1, err.Error(), nil); err == nil {
		_, _ = resp.Write(bytes)
	}
}

func InitApiServer() (err error) {
	var (
		mux           *http.ServeMux
		listener      net.Listener
		httpServer    *http.Server
		staticDir     http.Dir     // 静态文件根目录
		staticHandler http.Handler // 静态文件的HTTP回调
	)
	// 配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)
	mux.HandleFunc("/job/log", handleJobLog)
	mux.HandleFunc("/worker/list", handleWorkerList)

	// 静态文件目录
	staticDir = http.Dir(G_config.WebRoot)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))

	// 启动TCP监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	// 创建一个Http服务
	httpServer = &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	// 单例赋值
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	go httpServer.Serve(listener)

	time.Sleep(100 * time.Millisecond)
	fmt.Printf("服务已启动，端口：%d", G_config.ApiPort)

	return
}
