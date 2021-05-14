package config

import (
	"fmt"
	nacosRemote "github.com/kms9/go-viper-nacos"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/silenceper/log"
	"github.com/spf13/viper"
	"time"
)

var NacosConfig  *viper.Viper

var DefaultConfigMap = map[string]map[string]string{
	"208" : {
		"serverAdd" : "10.8.8.208",
		"dataId"    : "qq-config",
		"group"     : "DEFAULT_GROUP",
		"nameSpaceId" : "2baea186-51ba-459e-bfa8-4c222d16c308",
	},

	"201": {
		"serverAdd" : "192.168.31.201",
		"dataId"    : "qq-config",
		"group"     : "DEFAULT_GROUP",
		"nameSpaceId" : "4793f393-e460-43df-bace-99c8ba4cbe06",
	},
}

var  SChan chan string

func init()  {
	SChan = make(chan string, 1)
}


func ReadChan()  {
	for {
		tmp := <-SChan
		fmt.Println("已经执行了", tmp)
	}
}

func SendChan()  {
	SChan <- "TestKey:"+time.Now().String()
}

// StartLogger 初始日志
func StartNacosConfig() error {
	endpoint := "http://192.168.31.201:8848"
	path := ""
	tconfig:= DefaultConfigMap["208"]
	NacosConfig = viper.New()

	var (
		serverAdd 	= tconfig["serverAdd"]
		dataId    	= tconfig["dataId"]
		group     	= tconfig["group"]
		nameSpaceId = tconfig["nameSpaceId"]
		//port      = 8848
	)

	//config := viper.New()
	sc := []constant.ServerConfig{
		{
			IpAddr: serverAdd,
			Port:   8848,
			Scheme: "http",
			ContextPath: "/nacos",
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         nameSpaceId, //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./config/log",
		CacheDir:            "./config/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
		//ListenInterval:
	}

	params:=vo.NacosClientParam{
		ClientConfig:  &cc,
		ServerConfigs: sc,
	}

	nacosRemote.SetDataID(dataId)
	nacosRemote.SetGroup(group)
	nacosRemote.SetOnChangeCallback(SendChan)

	nacosRemote.SetNacosOptions(params)
	NacosConfig.SetConfigType("yaml")

	tc:=nacosRemote.GetAllRemoteConfig()
	_ = NacosConfig.MergeConfig(tc)


	//Logger.Logger.Info("nacosRemote set config")

	err:= NacosConfig.AddRemoteProvider("nacos", endpoint, path)
	if err!=nil {
		log.Info("NacosRConfig.AddRemoteProvider err: " + err.Error())
		return err
	}

	err = NacosConfig.ReadRemoteConfig()
	if err!=nil {
		log.Info("NacosRConfig.ReadRemoteConfig err: "+err.Error())
		return err
	}

	err = NacosConfig.WatchRemoteConfigOnChannel()
	if err!=nil {
		log.Info("NacosRConfig.WatchRemoteConfigOnChannel err: " + err.Error())
		return err
	}

	return nil
}



