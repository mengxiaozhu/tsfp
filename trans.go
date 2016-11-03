package tsfp

import (
	"bytes"
	"gopkg.in/macaron.v1"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"
)

var CookieJar *cookiejar.Jar
var HttpClient *http.Client

var lock *sync.Mutex

func Trans(ctx *macaron.Context, addr *url.URL) {
	lock.Lock()
	defer lock.Unlock()

	var resp *http.Response
	var err error

	HttpClient.Jar.SetCookies(addr, ctx.Req.Cookies())

	if ctx.Req.Method == "POST" {
		bs, err := ctx.Req.Body().Bytes()
		if err != nil {
			ctx.PlainText(400, []byte(err.Error()))
			return
		}
		resp, err = HttpClient.Post(addr.String() + ctx.Req.RequestURI, ctx.Req.Header.Get("Content-Type"), bytes.NewReader(bs))
	} else {
		resp, err = HttpClient.Get(addr.String() + ctx.Req.RequestURI)
	}

	if err != nil {
		ctx.PlainText(400, []byte(err.Error()))
		return
	}

	defer resp.Body.Close()

	io.Copy(ctx.Resp, resp.Body)
}

func KeepSession(server *Server) {
	for {
		lock.Lock()
		if server.Init != "" {
			HttpClient.Get(server.Init)
		}
		lock.Unlock()
		time.Sleep(3 * time.Second)
	}
}

func NewProxy(conf *Config) {
	m := macaron.Classic()
	m.Use(macaron.Renderer())

	// 静态文件
	m.Use(macaron.Static(conf.Static))

	// HttpClient
	CookieJar, _ = cookiejar.New(nil)
	HttpClient = &http.Client{
		Jar: CookieJar,
	}
	lock = &sync.Mutex{}

	// 保持会话
	for _, server := range conf.Servers {
		go KeepSession(&server)
	}

	// 路由转发
	for _, server := range conf.Servers {
		RegisterRoute(&server, m)
	}

	// 启动代理
	m.Run(conf.Port)
}

func RegisterRoute(server *Server, m *macaron.Macaron) {
	addrURL, err := url.Parse(server.Addr)
	if err != nil {
		return
	}

	m.Any(server.Pattern, func(ctx *macaron.Context) {
		Trans(ctx, addrURL)
	})
}