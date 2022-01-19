package main

import (
	"fmt"
	"log"
	"nmapwithgo/cli"
	"nmapwithgo/parse"
	"nmapwithgo/read"
	"nmapwithgo/scan"
	"os"
	"sync"
	"time"
	"strconv"
)

func main(){
	if !* cli.Scan{
		fmt.Println("No scan command !")
		os.Exit(0)
	}
	if !read.Isexist(* cli.Path){
		fmt.Println("Path error !")
		os.Exit(0)
	}
	if *cli.Output==""{
		fmt.Println("Output file is None")
	}else{
		fmt.Println("Output file is:",*cli.Output)
	}
	// 创建 maxroutine 值从命令行选项获取参数，cli.Rate 指针就是这货
	var maxroutinenum int=*cli.Rate
	// 创建扫描协程的状态通知信道,最大缓存为 maxroutinenum 值
	var ch chan int=make(chan int,maxroutinenum)
	// 父 routine 要等待子 routine 完结
	var waitgp sync.WaitGroup
	eachline:=read.Readline(*cli.Path)
	// 生成对 ip 文件的闭包函数
	var writeto func(string,string,string,string,[]string)
	var file *os.File
	start:=time.Now()
	if *cli.Output!=""{
		timestamp:=strconv.FormatInt(time.Now().Unix(),10)
		fileaddr:=timestamp+*cli.Output
		file,writeto=read.Writetooutput(fileaddr)
		log.Println("待解析数",read.Lineamount(*cli.Path))
		defer file.Close()
		for i:=0;i<read.Lineamount(*cli.Path);i++{
			// 添加守护
			waitgp.Add(1)
			// 开始，写入信道
			ch <- 1
			// 启动协程
			go func(j int){
				defer waitgp.Done()
				time.Sleep(1*time.Second)
			// parse.Parseport 是从 masscan 扫描结果中提取端口
			// scan.MasscanLine 是把 ip 拿去扫描
			// parse.Parseip 是把文件中的 ip 地址逐行提取
				host:=parse.Parsehost(eachline(j))
				tag:=parse.Parsetag(eachline(j))
				ip:=parse.Parseip(eachline(j))
				if ip!=""{
					massresult:=scan.MasscanLine(ip)
					port:=parse.Parseport(massresult)
					//log.Printf("\nIp : %v\nHost :%v\nTag :%v",ip,host,tag)
					//log.Println()
					if port!=""{
						result:=scan.Nmapopen(ip,port)
						stuff:=parse.ParseNmapresult(result)
						writeto(host,tag,ip,port,stuff)
					}
				// 运行结束，读取信道
				<-ch
			}
		}(i)
	}
		// 等待所有 go 协程走完
		waitgp.Wait()
		fmt.Println("Finish scanning!")
		// 看一下耗时
		fmt.Println("Cost:",time.Now().Sub(start))
		fmt.Println(fileaddr)
	}
}
