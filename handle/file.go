package handle

import (
	"net/http"
)

//静态文件转接

//网站图标转移
//param w http.ResponseWriter 写入http句柄
//param r *http.Request 读取http句柄
func FileFavicon(w http.ResponseWriter, r *http.Request){
	goURL(w, r, "/assets/imgs/favicon.ico")
}
