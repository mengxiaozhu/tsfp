package main

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"gopkg.in/macaron.v1"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"qiniupkg.com/x/log.v7"
	"time"
)

type Server struct {
	Addr    string
	Pattern string
	Init    string
}
type Config struct {
	Port    int
	Init    string
	Addr    string
	Static  string
	Servers map[string]Server
}

var Conf Config
var CookieJar *cookiejar.Jar
var HttpClient *http.Client

func Trans(ctx *macaron.Context, addr *url.URL) {

	var resp *http.Response
	var err error
	log.Println(ctx.Req.Method, addr, ctx.Req.RequestURI)

	HttpClient.Jar.SetCookies(addr, ctx.Req.Cookies())
	if ctx.Req.Method == "POST" {
		bodyBytes, e := ctx.Req.Body().Bytes()
		if e != nil {
			ctx.PlainText(400, []byte(err.Error()))
		}
		bodyReader := bytes.NewReader(bodyBytes)
		resp, err = HttpClient.Post(addr.String()+ctx.Req.RequestURI, ctx.Req.Header.Get("Content-Type"), bodyReader)
	} else {
		resp, err = HttpClient.Get(addr.String() + ctx.Req.RequestURI)
	}
	if err != nil {
		ctx.PlainText(400, []byte(err.Error()))
	}
	defer resp.Body.Close()
	io.Copy(ctx.Resp, resp.Body)
}

func main() {
	Conf = Config{}
	if _, err := toml.DecodeFile(".tsfp.toml", &Conf); err != nil {
		log.Fatal(err.Error())
	}
	log.Println(Conf)

	CookieJar, _ = cookiejar.New(nil)
	HttpClient = &http.Client{
		Jar: CookieJar,
	}

	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Use(macaron.Static(Conf.Static))

	for _, server := range Conf.Servers {
		addrURL, err := url.Parse(server.Addr)
		if err != nil {
			continue
		}

		m.Any(server.Pattern, func(ctx *macaron.Context) {
			Trans(ctx, addrURL)
		})

	}

	go func() {
		for {
			for _, server := range Conf.Servers {
				log.Println(server.Init)
				HttpClient.Get(server.Init)
			}
			time.Sleep(3 * time.Second)
		}
	}()

	m.Run(Conf.Port)
}
