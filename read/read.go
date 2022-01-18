package read

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	_	"time"
	_	"strconv"
)

//判断文件是否存在
func Isexist(fileaddr string)bool{
	_,err:=os.Stat(fileaddr)
	if err!=nil{
		log.Println(err)
		if os.IsExist(err){
			return true
		}
		log.Println("Ippath not  exist!")
		return false
	}
	log.Println("Ippath exist!")
	return true
}

//读取文件某行
func Readline(fileaddr string)func(int)string{
	// 设定匿名函数做闭包
	// 设定耗时
	var Handlefile *os.File
	Handlefile,err:=os.Open(fileaddr)
	defer Handlefile.Close()
	if err!=nil{
		log.Println(err)
		os.Exit(0)
	}
	// 延迟关闭文件
	// 逐行读取 linereader 里面的信息
	var i int=0
	var content []string
	Linereader:=bufio.NewReader(Handlefile)
	for{
		line,_,err:=Linereader.ReadLine()
		if err==io.EOF{
			break
		}
		content=append(content,string(line))
		i=i+1
		}
	return func(linenum int)string{
		return content[linenum]
	}
}


// 统计文件行数
func Lineamount(fileaddr string)int{
	var Handlefile *os.File
	Handlefile,err:=os.Open(fileaddr)
	defer Handlefile.Close()
	if err!=nil{
		log.Println(err)
		os.Exit(0)
	}
	linereader:=bufio.NewReader(Handlefile)
	i:=0
	for {
		// 循环按行读取
		_,_,err:=linereader.ReadLine()
		if err==io.EOF{
			break
		}
		i=i+1
	}
	return i
}

// 读取整个文件
func Readall(fileaddr string){
	bytes,err:=ioutil.ReadFile(fileaddr)
	if err!=nil{
		log.Fatal(err)
	}
	log.Println("Bytes read is :",len(bytes))
	log.Println("String read is :",string(bytes))
}

// 写入文件
func Writetooutput(fileaddr string)(*os.File,func(string,string,string,string,[]string)){
//	timestamp:=strconv.FormatInt(time.Now().Unix(),10)
//	fileaddr=timestamp+fileaddr
	var writefile *os.File
	writefile,err:=os.OpenFile(fileaddr,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	if err!=nil{
		log.Println(err)
	}
	_,err=writefile.WriteString("host\ttag\tip\tport/protocal\tstatus\tservice\tremark/version\n")
	if err!=nil{
		log.Panicln(err)
	}
	return writefile,func(host string,tag string,ip string,port string,result[]string){
		writefile.WriteString(host+"\t"+tag+"\t"+ip+"\t"+"\t"+"\t"+"\t"+port+"\n")
		for _,value:=range result{
			if value!=""{
				writefile.WriteString(host+"\t"+tag+"\t"+ip+"\t"+value+"\n")
	//			log.Println(value)
			}
		}
	}
}
