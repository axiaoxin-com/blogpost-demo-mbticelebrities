// main_test.go

package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	// 这里，我们创建一个新的HTTP请求实例。这个请求将被传递给 handler。
	// 第一个参数是请求方法，第二个参数是路由(我们暂时把它留空)，第三个参数是请求体，在本例中没有。
	req, err := http.NewRequest("GET", "", nil)
	// 如果在创建请求实例时出现错误，我们会失败并停止测试
	if err != nil {
		t.Fatal(err)
	}

	// 我们使用 Go 的 httptest 库来创建一个 http 记录器。
	// 这个记录器将作为我们 http 请求的目标（你可以把它想象成一个迷你浏览器，它将接受我们发出的http请求的结果）。
	recorder := httptest.NewRecorder()

	// 使用我们的 handler 创建一个 HTTP 处理程序。 “handler”是我们要测试的 main.go 文件中定义的处理函数
	hf := http.HandlerFunc(handler)

	// 将HTTP请求发送到我们的记录器。这一行实际上执行了我们想要测试的 handler 处理程序
	hf.ServeHTTP(recorder, req)

	// 检查状态码是否符合我们的预期。
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 检查响应 body 是否符合我们的预期。
	expected := `Hello World!`
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestRouter(t *testing.T) {
	// 使用前面定义的构造函数实例化路由器
	r := newRouter()

	// 使用“httptest”库的 `NewServer` 方法创建一个新服务器
	// 文档：https://golang.org/pkg/net/http/httptest/#NewServer
	mockServer := httptest.NewServer(r)

	// mockServer 运行了一个服务器，通过它的 URL 属性对外暴露了服务器地址
	// 我们向路由器中定义的“hello”路由发起一个GET请求
	resp, err := http.Get(mockServer.URL + "/hello")
	// 处理任何意外错误
	if err != nil {
		t.Fatal(err)
	}

	// 我们希望状态码是 200 (ok)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be ok, got %d", resp.StatusCode)
	}

	// 在接下来的几行中，读取响应体，并将其转换为字符串
	defer resp.Body.Close()
	// 将body读入一串字节中(b)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	// 将 bytes 转换为 string
	respString := string(b)
	expected := "Hello World!"

	// 我们希望我们的响应与我们的处理程序中定义的响应相匹配。
	// 如果它碰巧是“Hello world!”，那么确认路由是正确的
	if respString != expected {
		t.Errorf("Response should be %s, got %s", expected, respString)
	}
}

func TestRouterForNonExistentRoute(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)
	// 大部分代码都是相似的。唯一的区别是现在我们向一个没有定义的路由发出请求，比如 POST /hello 路由。
	resp, err := http.Post(mockServer.URL+"/hello", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 我们希望我们的状态为 405（方法不允许）
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Status should be 405, got %d", resp.StatusCode)
	}

	// 测试 body 的代码也基本相同，但这次，我们期望得到的是一个空 body
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	respString := string(b)
	expected := ""

	if respString != expected {
		t.Errorf("Response should be %s, got %s", expected, respString)
	}
}

func TestStaticFileServer(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)

	// 我们想要点击 `GET /assets/` 路由来获取 index.html 文件响应
	resp, err := http.Get(mockServer.URL + "/assets/")
	if err != nil {
		t.Fatal(err)
	}

	// 希望状态码是 200 (ok)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be 200, got %d", resp.StatusCode)
	}

	// 测试 HTML 文件的全部内容是不明智的。
	// 我们测试 content-type 标头是“text/html; charset=utf-8”，这样我们就知道已经提供了一个 html 文件
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("Wrong content type, expected %s, got %s", expectedContentType, contentType)
	}
}
