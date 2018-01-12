[![Build Status](https://travis-ci.org/u35s/h2s.svg?branch=master)](https://travis-ci.org/u35s/h2s)


# h2s
h2s是一个把socks5代理转换为http,https代理的小工具,同时支持socks5,http,https

# 安装
```golang
go get -v github.com/u35s/h2s
```

# 启动
h2s默认会读取当前文件夹下的config.json配置,也可以用-c指定配置文件

```
h2s -c example.config.json
```

h2s也可以通过命令行参数指定相关参数

```
h2s -s "socks5.com:8388" -P "0.0.0.0:8088" -b "0.0.0.0" -l 1080 
```

* -P 本地http,https代理地址 eg: 0.0.0.0:8088
* -s socks5服务器地址,当此地址包含端口时会忽略-p参数
* -p socks5服务器端口
* -b 本地socks5代理地址 eg: 127.0.0.1
* -l 本地socks5代理地址端口 eg: 1080
* -m 加密方法 default: aes-256-cfb
* -d 开启调试日志

# 测试
终端下执行

```bash
root@hostname https_proxy=127.0.0.1:8088 http_proxy=127.0.0.1:8088 curl -v ip.cn
* Rebuilt URL to: ip.cn/
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 8088 (#0)
> GET http://ip.cn/ HTTP/1.1
> Host: ip.cn
> User-Agent: curl/7.54.0
> Accept: */*
> Proxy-Connection: Keep-Alive
>
< HTTP/1.1 200 OK
< Server: nginx/1.8.0
< Date: Fri, 12 Jan 2018 12:58:41 GMT
< Content-Type: text/html; charset=UTF-8
< Transfer-Encoding: chunked
< Connection: keep-alive
< Vary: Accept-Encoding
< X-Powered-By: PHP/5.6.32-1~dotdeb+7.1
<
当前 IP：47.89.180.x 来自：美国 阿里云
* Connection #0 to host 127.0.0.1 left intact
``` 