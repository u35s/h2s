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
