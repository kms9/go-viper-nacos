package main

import (
	"github.com/silenceper/log"
	"os"
	"os/signal"
	"time"
	"fmt"
)

func main()  {
	err:=StartNacosConfig()
	if err!=nil{
		log.Error("StartNacosConfig err:"+ err.Error())
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt,os.Kill)
	go TestViper()
	s := <-c
	fmt.Println("stop,signal:",s)
}

func TestViper()  {
	for {
		time.Sleep(time.Second*5)
		fmt.Println(NacosConfig.GetString("yc.cashApi.scene"))
	}
}
