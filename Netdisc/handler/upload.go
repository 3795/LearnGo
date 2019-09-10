package handler

import (
	"LearnGo/Netdisc/db"
	"LearnGo/Netdisc/meta"
	"LearnGo/Netdisc/util"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/static/view/index.html", http.StatusFound)
		return
	}
}

func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	fileSha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(fileSha1)

	// 删除文件
	_ = os.Remove(fMeta.Location)
	// 删除文件元信息
	meta.RemoveFileMeta(fileSha1)

	// todo 删除表文件信息
	w.WriteHeader(http.StatusOK)
}

// 秒传接口
func TryFastUpLoadHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	// 解析参数
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	// 从文件表中查询相同hash的文件记录
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 找不到对应记录，则秒传失败
	if fileMeta == nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		_, _ = w.Write(resp.JSONBytes())
		return
	}

	// 上传过则将文件信息写入用户文件表，返回成功
	suc := db.OnUserFileUploadFinished(username, filehash, filename, int64(filesize))
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		_, _ = w.Write(resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: -2,
		Msg:  "秒传失败，请稍后重试",
	}
	_, _ = w.Write(resp.JSONBytes())
	return
}
