# go-snowflake

go-snowflake是一个基于Twitter开源的snowflake算法而包装成[Go](https://golang.org)的包，具有如下特点

* 完全开源算法的序列号生成器
* 可以自定义生成器包括的每段长度，如服务器长度、进程数长度、毫秒内自增长度
* 默认使用毫秒作为高位时间戳，整体序列号自增


## Install

首先安装好golang对应基础环境，如果还没有安装好，请[查看这里](https://golang.org/doc/install)

```
go get github.com/ytf606/go-snowflake
```

## Example

```
package main

import (
    "fmt"
    "github.com/ytf606/go-snowflake"
)

func main() {
    var serverId int64 = 1
    arrayChan := make(chan int64, 20)
    for t := 0; t < 20; t++ {
        go func(ts int64, ch chan int64){
            s, _ := snowflake.NewProcessWork(serverId, ts)
            id, _ := s.Id()
            ch <- id
        }(int64(t), arrayChan)
    }
    arrayResult := [20]int64{0}
    for i := 0; i < 20; i++ {
        arrayResult[i] = <-arrayChan
    }
    fmt.Println(arrayResult)
}
```

## TODO

* 更加灵活的定义`毫秒时间戳`+`服务器`+`进程`+`自增`
* 批量获取多个自增式序列号
* 支持标准RESTful方式的接口化
* 多种数据格式
* ......
