# zebra - 一个黑白分明的 Go Web 框架

![zebra](logo.png)

zebra 是一个追求简单实用，但同时功能特性完整、易扩展、编码灵活自由的 Golang Web 框架。
它不依赖任何第三方包。


## 如何使用

使用 zebra 只需要创建一个 zebra 服务并启动它：

``````go
z := zebra.New()
z.Run()
``````

zebra默认在 `3000` 端口运行，你也可以通过 `z.RunOnAddr(":8888")` 方法指定端口。另外，
还可以通过`http.Server`来配置完整的服务信息，以便让zebra按照你自定义的方式运行：

``````go
s := &http.Server{
    Addr:           ":8080",
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
}
z := zebra.NewWithServer(s)
z.Run()
``````



## 路由
zebra的路由是满足 RESTful 规则的路由，但需要注意的是，它并不支持完整的HTTP Method,
我们只保留了一些常用的方法，如下：

    "GET"
    "POST"
    "PUT"
    "OPTIONS"
    "DELETE"
    "HEAD"

我们相信这些方法足够您实现一个完整的 RESTful Web 服务。但是如果它们真的不足以满足您的需求，
不要担心，zebra的路由中间件非常容易扩展。

### 创建 Router
使用 `zebra.NewRouter()` 方法即可创建一个Router.
``````go
r := zebra.NewRouter()
r.Get("/foo", func(c *zebra.Captain) {
    // do some thing..
})
app := zebra.New()
app.Use(r)
app.Run()
``````

### 路由规则
zebra的路由规则非常简单，使用它时你只需要记住一个能满足 80% 场景的规则 `:path`
即可，例如我们为用户信息创建如下Router：

``````go
// http GET:  localhost:8888/user/zhangsan
r.Get("/user/:name", func(c *zebra.Captain) {
    // 使用 c.Path("name") 得到请求中的路径参数 'zhangsan'
})
``````

我们在实践中发现，这种方式能解决绝大多数路由规则需求，但并不是万能的。
你总会遇到一些特殊的匹配规则，此时，只需要编写自己的正则即可。
zebra路由直接支持正则表达式，只需要将它放在一对花括号中：

``````go
// http GET:  localhost:8888/bar/badman123
r.Get("/bar/:userName{badman[\\d]+}", func(c *Captain) {
    // c.Path("userName")  --->  badman123
})
// 请求 localhost:8888/bar/goodboy1314 则不会被匹配
``````

## 日志

## favicon

## 静态服务

## 如何编写自己的中间件

## 整合其它应用

zebra灵活的功能特性，允许你轻松整合 `http.ServeMux` 甚至其它第三方应用，如 `martini` .

``````go
mux := http.NewServeMux()
mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hello zebra"))
})

z := zebra.New()
z.SetName("Zebra")
z.Use(mux)
z.Run()
``````

另外，一个已完成的zebra应用，也可以被轻松地整合到一个全新的zebra应用，构成它的一部分。

``````go
subApp := zebra.New()
subApp.SetName("bar")
// ...

app := zebra.New()
app.Use(subApp)
app.Run()
``````



## 联系我们

有问题可以直接创建一个issue

PS: zebra目前处于开发阶段，欢迎有同样想法的朋友加入。