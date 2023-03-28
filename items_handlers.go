package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Item struct {
	Name string `json:"name"`
	MBTI string `json:"mbti"`
}

var items []Item

func getItemsHandler(w http.ResponseWriter, r *http.Request) {
	// 将“items”变量转换为 json
	itemsBytes, err := json.Marshal(items)
	// 如果有错误，打印到控制台，并返回一个服务器错误给用户
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 如果一切顺利，将名人 MBTI 的 JSON 列表写入响应
	w.Write(itemsBytes)
}

func createItemHandler(w http.ResponseWriter, r *http.Request) {
	// 创建 Item 实例
	item := Item{}

	// 我们使用 `ParseForm` 方法，解析请求发送过来的 HTML 表单数据
	err := r.ParseForm()
	// 如果出现任何错误，我们将返回错误给用户
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 从表单信息中获取对应字段的信息
	item.Name = r.Form.Get("name")
	item.MBTI = r.Form.Get("mbti")

	// 将新的 item 添加我们现有的 items 列表
	items = append(items, item)

	// 最后，我们使用 http 库的 `Redirect` 方法将用户重定向到原始的 HTMl 页面（位于 `/assets/`）
	http.Redirect(w, r, "/assets/", http.StatusFound)
}
