package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type cacheLine struct{
	tag   uint  //标记位
	lru   uint  // lru count
	valid uint  // 有效位
}

var hits uint  //命中数
var miss uint  //未命中数
var evictions uint  //淘汰数

var s uint
var b uint
var E uint
var S uint
var verbose uint = 1
var eviction uint
var t string = "./yi.trace"

var  line []cacheLine = make([]cacheLine,E)
var  cache [][]cacheLine = make([][]cacheLine,S)

func initCache(){
	 cache  = make([][]cacheLine,S)
	for i:=0;i<len(cache);i++{
		cache[i] = make([]cacheLine,E)
	}
}

//解析参数
func  ParseArgument(){
	//var err error = nil
	for i,v:= range os.Args{
		if i==0{
			continue
		}
		arg:=strings.Split(v,"=")[0]
		p:=strings.Split(v,"=")[1]

		switch arg {
		case "-v":
			verbose = 1
			break
		case "-s":
			i,_:=strconv.Atoi(p)
			s = uint(i)
			S = 1<<s
			break
		case "-E":
			i,_:=strconv.Atoi(p)
			E = uint(i)
			break
		case "-b":
			i,_:=strconv.Atoi(p)
			b = uint(i)
			break
		case "-t":
			t = p
			break
		default:
			continue
		}
	}
}


//读取文件
func ReadFile(t string){
		file,err:=os.Open(t)
		defer file.Close()

		if err!=nil{
			panic("文件打开失败")
		}

		br:=bufio.NewReaderSize(file,20)

		for{
			line,_,err :=br.ReadLine()

			if err!=nil && err !=io.EOF{
				panic("文件读取错误")
			}

			if err == io.EOF{
				return
			}

			//解析指令
			ParseInstruction(line)
		}
}


//解析指令
func ParseInstruction(line []byte){
	var address uint
	var operation uint8
	var size uint

	fmt.Sscanf(string(line)," %c %x,%d",&operation,&address,&size)

	//fmt.Printf(" %c %x,%d\n",operation,address,size)
	//fmt.Println(address)
	var ret int
	if line[0] == 'I'{
		return
	}else{
		switch operation {
		case 'S':
			ret = VisitCache(address)
		case 'L':
			ret = VisitCache(address)
		case 'M':
			ret = VisitCache(address)
			hits++
		}

		if verbose==1{
			switch ret{
			case 0:
				fmt.Printf("%c %x,%d hit\n", operation, address, size)
			case 1:
				fmt.Printf("%c %x,%d miss\n", operation, address, size)
			case 2:
				fmt.Printf("%c %x,%d miss eviction\n",operation, address, size)
			}
		}
	}
}

//读取cache
func VisitCache(address uint) int{

  	//获取tag
	var tag uint =address>>(s+b)
	var index uint =address >> b & ((1 << s) - 1)
	cacheSet:=cache[index]
	evict:=0 //驱逐行初始化为0
	empty:=-1 //表示空行

	//判断是否命中索引
		for i,v:=range cacheSet{

			//缓存命中,lru置为1，表示最近访问过
			if v.valid == 1{
				if v.tag == tag {
					v.lru =1
					hits++ //命中数加一
					return 0
				}


				//当前行缓存没有命中 lru++
				v.lru++
				if v.lru >= cacheSet[evict].lru{
					//驱逐行变成当前下标
					evict = i
				}

			}else{
				empty = i
			}
		}
		//如果存在空行
		if empty!=-1{
			cacheSet[empty].tag = tag
			cacheSet[empty].valid =1
			cacheSet[empty].lru = 1
			miss++
			return 1
		}

		//最终没有命中缓存，需要驱逐出一行
		cacheSet[evict].tag = tag
		cacheSet[evict].lru = 1
		miss++
		eviction++
		return 2

}


func main(){
	ParseArgument()
	initCache()
	ReadFile(t)
	fmt.Println(hits,miss,evictions)
}