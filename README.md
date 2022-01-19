# nmapUseGo

[toc]

## 功能
* 联动 masscan 和 nmap 扫描大量 ip
* masscan 扫到的有效信息传给 nmap
* 每一条 ip 扫描参数根据 masscan 结果配置
* 利用 goroutine 和 chan 控制

## 用法
  -ippath [ip 清单路径]
  -output [输出路径]
  -rate [线程数]
    	Default scan threads 10 个，建议不要超过 50 (default 10)
  -scan [可选参数]
    	Scan all ports
  -whitelist [可选参数]
    	Whitelist path


## 其他
* 使用原生 nmap ，需要执行
```shell
# ubuntu
sudo apt install nmap
# centos
yum install nmap
# or
dnf install nmap
```
* linux 下使用