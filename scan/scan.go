package scan

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"time"
)

// 执行 nmap 和 masscan 封装起来
// 一次性读取
func Masscan(ip string){
	start:=time.Now()
	cmd:=exec.Command("bash","-c","masscan "+ip+" -p1-65535 --rate 10000")
	log.Println(ip)
	output,err:=cmd.CombinedOutput()
	if err!=nil{
		log.Println(err)
	}
	log.Println(string(output))
	log.Println(time.Now().Sub(start))
}

// 逐行实时读取
func MasscanLine(ip string)(portbuffer string){
	start:=time.Now()
	if ip==""{
		return ""
	}
	cmd:=exec.Command("bash","-c","masscan "+ip+" -p1-65535 --rate 5000")
	stdout,err:=cmd.StdoutPipe()
	if err!=nil{
		log.Println(err)
	}
	cmd.Start()
	reader:=bufio.NewReader(stdout)
	log.Printf("Scanning %v !\n",ip)
	for {
		line,err2:=reader.ReadString('\n')
		if err2!=nil||io.EOF==err2{
			break
		}
		portbuffer=portbuffer+" "+line
	}
	log.Println(time.Now().Sub(start))
	cmd.Wait()
	return
}

// 数据传送给 Nmap 使用,一次性读取
func Nmapall(ip string,port string)string{
	if port!=""{
	cmd:=exec.Command("bash","-c","nmap -sS "+ip+" -p "+port+" -Pn  -n -sV")
	output,err:=cmd.CombinedOutput()
	if err!=nil{
		log.Println(err)
	}
	log.Println(string(output))
	return string(output)
}else{
// 解析不出信息则打印该提示，打印空字符串
	log.Println(ip,"has no port to scan,port:",port)
	return ""
}
}

// 只扫开放端口
func Nmapopen(ip string,port string)string{
	if port!=""{
		cmd:=exec.Command("bash","-c","nmap -sS "+ip+" -p "+port+" -Pn --open -n -sV")
		output,err:=cmd.CombinedOutput()
	if err!=nil{
		log.Println(err)
	}
	//log.Println(string(output))
	return string(output)
}else{
// 解析不出信息则打印该提示，打印空字符串
	log.Println(ip,"has no port to scan,port:",port)
	return ""
}
}

// 进攻性扫描
func Attackscan(ip string,port string){
	cmd:=exec.Command("bash","-c","nmap -sS "+ip+" -A"+"-T4 -Pn -n")
	output,err:=cmd.CombinedOutput()
	if err!=nil{
		log.Println(err)
	}
	log.Println(string(output))
}
