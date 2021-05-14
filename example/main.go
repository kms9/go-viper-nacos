package main

import (
	"fmt"
	"github.com/kms9/go-viper-nacos/example/config"
	"github.com/silenceper/log"
	"os"
	"os/signal"
	"time"
)

func main()  {
	err:= config.StartNacosConfig()
	if err!=nil{
		log.Error("StartNacosConfig err:"+ err.Error())
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt,os.Kill)
	go TestViper()
	go config.ReadChan()
	s := <-c
	fmt.Println("stop,signal:",s)
}

func TestViper()  {
	for {
		time.Sleep(time.Second*5)
		fmt.Println(config.NacosConfig.GetString("yc.cashApi.scene"), time.Now().String())
	}
}
