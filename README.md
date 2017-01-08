# atc
服务端调用HTTP API，并完成HTML模板渲染的工具。替代前端初始化时Ajax的调用方案，实现SEO友好的前后端分离。
在服务端调用API并做模板渲染，使用go template实现。
项目目的就是包装一些常用函数，供模板使用。

项目代码在 https://git.oschina.net/yinshuwei/atc

##下载与运行atc
	cd $GOPATH/src
    git clone https://git.oschina.net/yinshuwei/atc.git
    cd atc
    go build
    ./atc

用浏览器打开 http://localhost:8888/?name=man
结果
    Hello: man 
    Hello: b
说明安装成功

##config
    {
        "WebPath": "www",
        "Port": ":8888",
        "IsDev": true,
        "Page404": "/404.html",
        "Envs": {
            "a": "b"
        }
    }
WebPath 页面目录
Port 服务器端口
IsDev 是否为开发环境，开发环境关闭页面和模板缓存
Page404 404页面
Envs 环境变量，可以在页面模板中使用

##模板
所有的html文件都会被当做模板进行解析。
www中index.html
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <title>测试</title>
    </head>
    <body>

    Hello: {{.params.name}}
    <br/>
    Hello: {{.envs.a}}

    {{$data := .atc.GetAPI `http://www.weather.com.cn/data/cityinfo/101010100.html` }} 

    <table>
    {{range $key, $val := $data.weatherinfo}}
    <tr><td>{{$key}}: </td><td> {{$val}}</td></tr>
    {{end}}
    </table>

    </body>
    </html>
