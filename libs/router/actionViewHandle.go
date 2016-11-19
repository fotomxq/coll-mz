package router

import "net/http"

//查看数据句柄
func actionViewHandle(w http.ResponseWriter, r *http.Request){

}

//查看结构
type View struct {

}

//获取数据列表
func (view *View) GetList(page int,max int,sortField string,descOn bool) (map[int]map[string]string,error){
	return nil,err
}

//将文件反馈到浏览器
func (view *View) GetFileToHeader(){

}

//删除文件，但保留采集记录
func (view *View) DeleteFile(){

}