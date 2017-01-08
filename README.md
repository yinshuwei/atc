# atc
服务端调用HTTP JSON API，并完成HTML模板渲染的工具。替代前端初始化时Ajax的调用方案，实现SEO友好的前后端分离。

使用go template实现。

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
    -------------	-------------
    context:	客户 签收人: 朱代签 已签收 感谢使用圆通速递，期待再次为您服务
    ftime:	2017-01-02 17:29:15
    location:	
    time:	2017-01-02 17:29:15
    -------------	-------------
    context:	上海市浦东新区南汇公司(点击查询电话)薄** 派件中 派件员电话13918345153
    ftime:	2017-01-02 08:26:45
    location:	
    time:	2017-01-02 08:26:45
    ......

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

    {{$data := .atc.GetAPI `http://www.kuaidi100.com/query?type=yuantong&postid=560513101584` }} 

    <table>
    <tr><td>-------------</td><td>-------------</td></tr>
    {{range $data.data}}
    {{range $key, $val := .}}
    <tr><td>{{$key}}: </td><td> {{$val}}</td></tr>
    {{end}}
    <tr><td>-------------</td><td>-------------</td></tr>
    {{end}}
    </table>

    </body>
    </html>

.params.PARAM_NAME 用来获取页面URL上的参数
.envs.ENV_NAME 用来获取atc.json中配置的变量
.atc.GetAPI 通过GET方式获得API数据(这里是一个快递接口的调用)
后面是使用go template的语法使用数据

##函数
.atc.PostAPI 通过POST方式，获得body内容，并通过json解码，生成一个map

    arg1 url
    arg2 data json字符串
    result map[string]interface{}

.atc.GetAPI 通过URL获得body内容，并通过json解码，生成一个map

    arg1 url
    result map[string]interface{}

.atc.GetBody 通过URL获得body内容

    arg1 url
    result template.HTML

.atc.Set 设置值到页面环境上

    arg1 key
    arg2 value
    result nil

.atc.Get 从页面环境中取值

    arg1 key
    result value

.atc.Add 在原来的value(页面环境中)上加上新的value（value为int时）

    arg1 key
    arg2 value
    result nil

.atc.IsEnd 在行布局中，是否在行尾（如：10，5，true）

    arg1 index
    arg2 width
    result bool

.atc.Others 在行布局中，最后一行会缺一部分，用缺的个数作为长度，创建一个0值数组

    arg1 array
    arg2 width
    result array

.atc.Cut 数组切成两半，前一半试图长度为width

    arg1 array
    arg2 width
    result ArrPair（First，Second）

.atc.Ter 三目，bool为true,返回value1,否则value2

    arg1 bool
    arg2 value1
    arg3 value2
    result value1 or value2

.atc.At 用下标取array元素

    arg1 array
    arg2 index
    result value

.atc.F2i float64 to int

    arg1 float64
    result int

.atc.Len 数组长度

    arg1 array
    result length

.atc.Arr 可变长的一组数据转成数组

    arg1 interface{} ... 可变长数据
    result array

.atc.SetTo 往一个map变量上设置值

    arg1 map
    arg2 key
    arg3 value
    result nil

.atc.Ref 引用一个文件

    arg1 path
    result template.HTML 

.ints.Add 加法

    arg1 int
    arg2 int
    result int

.ints.Sub 减法

    arg1 int
    arg2 int
    result int

.ints.Mod 取模

    arg1 int
    arg2 int
    result int

.ints.Arr 创建指定长度的int数组，值全为0

    arg1 length
    result array

.ints.Int 转换为int

    arg1 interface{}
    result int

.strs.Substr 子字符串

    arg1 str
    arg2 start
    arg3 end
    result string

.strs.Split 分割字符串

    arg1 str
    arg2 sep
    result str array

.strs.QueryEscape URL编码

    arg1 str
    result string
