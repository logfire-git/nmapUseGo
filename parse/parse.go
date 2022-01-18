package parse

import (
	"log"
	"regexp"
	"strings"
)

// 解析 ip 地址
func Parseip(line string)string{
	reg:=regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	result:=make([]string,1)
	result=reg.FindAllString(line,-1)
	result=append(result,"")
	if result[0]==""{
		log.Println("Can not parse Ip")
	}
	return result[0]
}

// 解析 host 信息,格式为序号 主机 ip
func Parsehost(line string)string{
	reg:=regexp.MustCompile(`\d+\t.+\t\d+\.\d+\.\d+\.\d+`)
	reg2:=regexp.MustCompile(`\t.+\t`)
	result:=make([]string,1)
	result=reg.FindAllString(line,-1)
	strs:=""
	for _,str:=range result{
		result2:=reg2.FindAllString(str,-1)
		for _,newstr:=range result2{
			newstr=strings.Replace(newstr,"\t","",-1)
			strs=strs+newstr
		}
	}
	if strs==""{
		log.Println("Can not parse Host")
	}
	return strs
}

// 解析 tag 信息,格式:时间 资源组 页码数
func Parsetag(line string)string{
	reg:=regexp.MustCompile(`:\d+\t.+\t`)
	reg2:=regexp.MustCompile(`\t.+\t`)
	result:=make([]string,1)
	result=reg.FindAllString(line,-1)
	strs:=""
	for _,str:=range result{
		result2:=reg2.FindAllString(str,-1)
		for _,str:=range result2{
			str=strings.Replace(str,"\t","",-1)
			strs=strs+str
		}
	}
	if strs==""{
		log.Println("Can not parse Tag")
	}
	return strs
}


// 解析端口信息,portbuffer 来自 masscan
func Parseport(portbuffer string)string{
	if portbuffer==""{
		log.Println("No scan result to parse!")
		return ""
	}
	reg:=regexp.MustCompile(`port \d+`)
	var result []string
	result=reg.FindAllString(portbuffer,-1)
	strs:=""
	for index,str:=range result{
		str=strings.Replace(str,"port ","",-1)
		if index==0{
			strs=strs+str
		}else{
			strs=strs+","+str
		}
	}
	if strs!=""{
		return strs
	}else{
// 解析不出端口信息返回空字符串
	log.Println("No port found in this ip")
		return ""
	}
}


// 解析扫描结果(带服务信息)，来自 Nmap ,把结果存入字符串切片，待处理
func ParseNmapresult(result string)([]string){
	var strs []string
	// 选出 raw 结果里面关键信息
	rega,_:=regexp.Compile(`PORT.*\n[\s\S]*(?:service unrecognized)`) // 有未识别的服务,1个
	results:=rega.FindAllString(result,-1)
	log.Println("Length of results:",len(results))
	if len(results)==0{
		regb,_:=regexp.Compile(`PORT.*\n[\s\S]*(?:services unrecognized)`) // 有未识别的服务,2个及以上
		results=regb.FindAllString(result,-1)
		if len(results)==0{
			reg,_:=regexp.Compile(`PORT.*\n[\s\S]*(?:Service detection)`) // 所有服务均可识别
			results=reg.FindAllString(result,-1)
		}
	}
	// 再次匹配
	reg2:=regexp.MustCompile("\n.*")
	// 匹配空格
	reg3:=regexp.MustCompile(`\s+`)
	for _,str:=range results{
		newstr:=reg2.FindAllString(str,-1)
		for _,str:=range newstr{
			// 把回车和空格替换为空和制表符
			str=strings.Replace(str,"\n","",-1)
			log.Println(str)
			str=reg3.ReplaceAllString(str,"\t")
			strs=append(strs,str)
			//log.Println(str)
		}
	}
	return strs
}


// 数据和操作分离,函数链式处理
// 函数结构必须相同才可以放入一块切片中
// 结构为临时指针，字符串，函数切片组，函数格式一致
func Stringprocesschain(tmp *string,prostr string,processchain[]func(string)string){
	for _,proc:=range processchain{
		prostr=proc(prostr)
	}
	*tmp=prostr
}
