package main

import (
	"fmt"
	"mongosync"
	"time"
)

func main() {
	mongosync.Run("D:\\mongosync\\src\\config.yaml")

	fmt.Println("进程启动...")
	sum := 0
	for {
		sum++
		fmt.Println("sum:", sum)
		time.Sleep(time.Second)
	}

}
