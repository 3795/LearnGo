package handler

import (
	"net/http"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/static/view/index.html", http.StatusFound)
		return
	}
}

// 秒传接口
func TryFastUpLoadHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	// 解析参数
	//username := r.Form.Get("username")
	//filehash := r.Form.Get("filehash")
	//filename := r.Form.Get("filename")
	//filesize, _ := strconv.Atoi(r.Form.Get("filesize"))
	//
	//// 从文件表中查询相同hash的文件记录
	//fileMeta, err := meta.GetFileMetaDB(filehash)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//
	//// 找不到对应记录，则秒传失败
	//if fileMeta == nil {
	//	resp := util.RespMsg{
	//		Code: -1,
	//		Msg: "秒传失败，请访问普通上传接口",
	//	}
	//	_, _ = w.Write(resp.JSONBytes())
	//	return
	//}
}
