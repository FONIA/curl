package main

import (
	"curl"
	"fmt"
	"io/ioutil"
)

//普通get
func Get() {
	cli := curl.CreateHttpClient()
	resp, err := cli.Get("http://49.232.145.118:20171/api/v1/portal/news?newsType=10&page=1&limit=50")
	if err == nil {
		txt, err := resp.GetContents()
		fmt.Println(txt, err)
	}
}

//带请求头get
func Get2() {
	// 创建 http 客户端的时候可以直接填充一些公共参数，后续请求会复用
	cli := curl.CreateHttpClient(curl.Options{
		Headers: map[string]interface{}{
			"Referer": "http://vip.stock.finance.sina.com.cn",
		},
		SetResCharset: "GB18030",
		BaseURI:       "",
	})
	resp, err := cli.Get("http://hq.sinajs.cn/list=sz002594")
	if err == nil {
		txt, err := resp.GetContents()
		fmt.Println(txt, err)
	}
}

//带请求头get
func Get3() {
	cli := curl.CreateHttpClient()
	//  cli.Get 切换成 cli.Post 就是 post 方式提交表单参数
	//resp, err := cli.Post("http://127.0.0.1:8091/postWithFormParams", goCurl.Options{
	resp, err := cli.Get("https://www.oschina.net/search", curl.Options{
		FormParams: map[string]interface{}{
			"random": 12345,
			"scope":  "project",
			"q":      "golang",
		},
		Headers: map[string]interface{}{
			"Content-Type": "application/x-www-form-urlencoded;charset=gb2312",
		},
	})

	if err == nil {
		txt, err := resp.GetContents()
		fmt.Println(txt, err)
	}
}

//发送中文
func Get4() {
	cli := curl.CreateHttpClient()
	resp, err := cli.Get("http://139.196.101.31:2080/test_json.php", curl.Options{
		FormParams: map[string]interface{}{
			//"user_name":"你好，该字段发送出去的数据为简体中文编码",  // 对方站点只接受 简体中文，这种不编码直接发出去就会报错
			"user_name": cli.Utf8ToSimpleChinese([]byte("该字段发送出去的数据为简体中文编码")), // 第二个参数：默认编码为 GB18030，（GBK 、GB18030 都是简体中文，go编码器中没有 gb2312）
		},
		//Headers: map[string]interface{}{
		//	"Content-Type": "application/x-www-form-urlencoded;charset=gb2312",
		//},
	})
	if err == nil {
		txt, err := resp.GetContents()
		fmt.Println(txt, err)
	}
}

//post 提交json
func Post() {
	cli := curl.CreateHttpClient()
	resp, err := cli.Post("http://127.0.0.1:8091/post-with-json", curl.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
		},
		JSON: struct {
			Code int      `json:"code"`
			Msg  string   `json:"msg"`
			Data []string `json:"data"`
		}{200, "OK", []string{"hello", "world"}},
	})
	if err == nil {
		txt, err := resp.GetContents()
		fmt.Println(txt, err)
	}
}

//post 提交xml (以表单参数形式提交x-www-form-urlencoded)
func Post2() {
	cli := curl.CreateHttpClient()
	resp, err := cli.Post("http://www.webxml.com.cn/WebServices/ChinaZipSearchWebService.asmx/getSupportCity", curl.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		FormParams: map[string]interface{}{
			"byProvinceName": "重庆", // 参数选项：上海、北京、天津、重庆 等。这个接口在postman测试有时候也是很稳定，可以更换参数多次测试
		},
		SetResCharset: "utf-8",
		Timeout:       10,
	})
	if err == nil {
		txt, err := resp.GetContents()
		fmt.Println(txt, err)
	}
}

//post 提交xml (以表单参数形式提交raw)
func Post3() {
	cli := curl.CreateHttpClient(curl.Options{
		SetResCharset: "utf-8",
	})

	// 需要提交的 xml 数据格式，发送前请转换为以下文本格式
	// 正式业务我们的参数是动态的
	// 那么就事先需要定义好go语言的结构体，最终将绑定好参数的结构体转为xml格式数据
	// 关于结构体转 xml 格式代码参见：https://blog.csdn.net/f363641380/article/details/87651427
	xml := `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <getSupportCity xmlns="http://WebXml.com.cn/">
      <byProvinceName>上海</byProvinceName>
    </getSupportCity>
  </soap:Body>
</soap:Envelope>
`

	resp, err := cli.Post("http://www.webxml.com.cn/WebServices/ChinaZipSearchWebService.asmx", curl.Options{
		Headers: map[string]interface{}{
			"Content-Type": "text/xml; charset=utf-8",
			"SOAPAction":   "http://WebXml.com.cn/getSupportCity", //  该参数按照业务方的具体要求传递
		},
		XML:     xml,
		Timeout: 20,
	})
	if err == nil {
		txt, err := resp.GetContents()
		fmt.Println(txt, err)
	}
}

// 设置代理ip访问目标站点
func Proxy() {
	cli := curl.CreateHttpClient()

	resp, err := cli.Get("http://myip.top/", curl.Options{
		Timeout: 5.0,
		Proxy:   "http://39.96.11.196:3211", // 该ip需要自己去申请每日免费试用
	})
	if err == nil {
		txt, err := resp.GetContents()
		fmt.Println(txt, err)
	}
}

// 文件下载
// 参数一 > 要下载的资源地址
// 参数二 > 指定下载路径（服务器最好指定绝对路径）
// 参数三 > 文件名，如果不设置，那么自动使用被下载的原始文件名
func Down() {
	cli := curl.CreateHttpClient()
	_, err := cli.Down("http://139.196.101.31:2080/GinSkeleton.jpg", "./", "ginskeleton.jpg", curl.Options{
		Timeout: 60.0,
	})
	if err == nil {
		fmt.Println("ok")
	}
}

//获取cookie
func cookie() {
	cli := curl.CreateHttpClient()
	resp, err := cli.Get(`https://www.baidu.com`)
	if err == nil {
		// 全量获取cookie
		for index, value := range resp.GetCookies() {
			fmt.Printf("序号：%d, %s\n", index, value.String())
		}
		// 根据键获取指定的 cookie
		fmt.Println(resp.GetCookie("BAIDUID"))
	}
}

//提交cookie
func cookie2() {
	cli := curl.CreateHttpClient()
	resp, err := cli.Post("http://127.0.0.1:8091/post-with-cookies", curl.Options{
		Cookies: "cookie1=value1;cookie2=value2",
	})
	if err == nil {
		txt, err := resp.GetContents()
		fmt.Println(txt, err)
	}
}

//提交cookie2 并从 body 体读取返回值
func cookie3() {
	cli := curl.CreateHttpClient()
	resp, err := cli.Post("http://127.0.0.1:8091/post-with-cookies", curl.Options{
		Cookies: map[string]string{
			"cookie1": "value1",
			"cookie2": "value2",
		},
	})
	if err == nil {
		body := resp.GetBody()
		defer func() {
			_ = body.Close()
		}()
		// 如果请求的返回结果是从body体读取的二进制数据，必须使用 body.Close()  函数关闭
		// 此外必须注意的是，该函数是直接从缓冲区获取的二进制，对方的编码类型如果有中文（gbk系列）就会是乱码,需要自己转换，转换代码参见 getContents（） 函数
		if bytes, err := ioutil.ReadAll(body); err == nil {
			fmt.Printf("%s", bytes)
		}
	}
}

//put
func Put() {
	cli := curl.CreateHttpClient()
	resp, err := cli.Put("http://127.0.0.1:8091/put")
	if err == nil {
		txt, err := resp.GetContents()
		fmt.Println(txt, err)
	}
}

//delete
func del() {
	cli := curl.CreateHttpClient()
	resp, err := cli.Delete("http://127.0.0.1:8091/delete")
	if err == nil {
		txt, err := resp.GetContents()
		fmt.Println(txt, err)
	}
}

func main() {
	Get()
	Get2()
	Get3()
	Get4()
	Post()
	Post2()
	Post3()
	Proxy()
	Down()
	cookie()
	cookie2()
	cookie3()
}
