# zebra

![zebra](logo.png)

zebra 是一个追求简单实用，但同时功能特性完整、易扩展、编码灵活自由的Golang Web框架。
它不依赖任何第三方包。


## 如何使用zebra

使用zebra只需要创建一个zebra服务并启动它

``````go
z := zebra.New()
z.Run()
``````

zebra默认在`3000`端口运行，你也可以通过`z.RunOnAddr(":8888")`方法指定端口。另外，还
可以通过`http.Server`来配置完整的服务信息:

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


## 整合其它应用

zebra灵活的功能特性，允许你轻松整合`http.ServeMux`甚至其它第三方应用，如`martini`。

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

另外，一个已完成的zebra应用，也可以被轻松的整合到一个全新的zebra应用，构成它的一部分。

``````go
subApp := zebra.New()
subApp.SetName("bar")
// ...

app := zebra.New()
app.Use(subApp)
app.Run()
``````