package cli

import "fmt"
import "flag"

var Scan=flag.Bool("scan",false,"Scan all ports")
var Path=flag.String("ippath","","Ip lists path")
var Whitelist=flag.String("whitelist","","Whitelist path")
var Rate=flag.Int("rate",10,"Default scan threads 10 个，建议不要超过 50")
var Output=flag.String("output","","Output file path")

func init(){
	fmt.Println("Commands are collecting!")
	fmt.Println("==================================================")
	flag.Parse()
	fmt.Println("Begin scan all ports? ",* Scan)
	fmt.Println("Ippath is ",* Path)
	fmt.Println("Whitelist path is ",* Whitelist)
	fmt.Println("Output path is ",* Output)
	fmt.Println("Scan rate is",*Rate)
}


