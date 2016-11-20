package coll

//煎蛋采集脚本
func CollJiandan() (bool, error) {
	//采集煎蛋网首页数据
	b, err := collJiandanIndex()
	if err != nil {
		CollPg.SendErrorLog(err)
	}
	if b == false {
		return false, nil
	}
	//采集分页页面数据
	b, err = collJiandanPage()
	if err != nil {
		CollPg.SendErrorLog(err)
	}
	if b == false {
		return false, nil
	}
	//返回成功
	return true, nil
}

//采集煎蛋网首页数据
func collJiandanIndex() (bool, error) {
	//source识别
	collSource := "jiandan-index"
	//解析页面内容
	CollPg.SendLog("开始采集煎蛋网首页...")
	doc,err := CollPg.GetHtmlGoquery("http://jandan.net")
	if err != nil {
		CollPg.SendLog("无法获取煎蛋网首页数据。")
		return false, err
	}
	//获取妹子图ID
	docList := doc.Find("#list-girl").Children().Children().Children().Find(".acv_comment").Find("p").Find("a")
	//遍历节点
	for i := range docList.Nodes {
		node := docList.Eq(i)
		v, b := node.Attr("href")
		if b == false {
			CollPg.SendLog("解析页面失败，无法获取页面的节点数据。")
			return false, nil
		}
		if v != "" {
			s, err := CollPg.AutoAddData(collSource, v, "", false)
			if err != nil || s == "" {
				continue
			}
		}
	}
	//返回结果
	return true, nil
}

//采集煎蛋网分页页面内容
func collJiandanPage() (bool, error) {
	//source识别
	collSource := "jiandan"
	//解析页面内容
	CollPg.SendLog("开始采集煎蛋网分页数据...")
	nextPageURL := "http://jandan.net/ooxx"
	doc,err := CollPg.GetHtmlGoquery(nextPageURL)
	if err != nil {
		CollPg.SendLog("无法获取煎蛋网分页首页的数据。")
		return false, err
	}
	//注意，这里有一个无限循环和三层嵌套
	//只有当解析发现10个以上节点重复或失败后、或下一个页面不存在的时候跳出
	var b bool
	var errNum int = 0
	for {
		//根据当前URL地址，获取页面内的节点列
		nodes := doc.Find(".commentlist").Children()
		//遍历所有子节点的子a元素
		for i := range nodes.Nodes {
			aNodes := nodes.Eq(i).Find(".text p").Children().Filter("a")
			//遍历子a节点，将其保存到数据库
			for j := range aNodes.Nodes {
				nodeURL, b := aNodes.Eq(j).Attr("href")
				if b == false {
					CollPg.SendLog("无法找到子节点数据，在获取href的时候发现节点不存在。")
					errNum += 1
					continue
				}
				s, err := CollPg.AutoAddData(collSource, nodeURL, "", false)
				if err != nil || s == "" {
					CollPg.SendLog("无法保存数据。")
					errNum += 1
					continue
				}
				//如果成功获取了节点及数据，则将错误计时器回拨为0
				errNum = 0
			}
		}
		//如果错误计时器超过10，则跳出
		if errNum > 10{
			CollPg.SendLog("发生10个以上节点不存在的问题，说明后续可能不存在数据，自动跳出该采集器。")
			break
		}
		//获取下一页URL
		nextPageURL, b = doc.Find(".previous-comment-page").Eq(0).Next().Attr("href")
		if b == false || nextPageURL == "" {
			CollPg.SendLog("无法获取煎蛋网分页下一个页面URL地址。")
			break
		}
		doc,err = CollPg.GetHtmlGoquery(nextPageURL)
		if err != nil {
			CollPg.SendLog("无法获取煎蛋网分页下一页的URL数据。")
			break
		}
		CollPg.SendLog("进入下一个页面，URL : " + nextPageURL)
	}
	//返回结果
	return true, nil
}
