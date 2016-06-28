# tsfp
a static and proxy server

```
port = 4040 // port that proxy listen
static = "./" // static file path
[servers]
    [servers.baidu] // route baidu
    pattern = "/api/*" // route pattern
    addr = "http://www.baidu.com" // route target
    init = "http://www.baidu.com" // get the addr per 3s to keep session
    [servers.google] // route google
    pattern = "/g/*" // route pattern
    addr = "http://www.google.com" // route target
    // no need to keep session
```
