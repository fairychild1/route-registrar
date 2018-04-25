package utility

import (
	"fmt"
	"route-registrar/config"
	"code.cloudfoundry.org/lager"
	"time"
	"route-registrar/messagebus"
	"github.com/nu7hatch/gouuid"
	"strconv"
	"github.com/astaxie/beego"
)

var Logger lager.Logger
var instanceId string
var Msgbus messagebus.MessageBus
var Config config.Config
func init() {
	aUUID, _ := uuid.NewV4()
	instanceId = aUUID.String()
	AllRoutesChan = make(map[string] chan struct{})
}

//每个route由一个协程运行，协程周期性的像nats广播route的信息。协程有一个chan，当向这个chan发送一个struct{}的时候，该协程会退出。一个route的ip地址加端口租车一个string，这个string作为key，chan作为value，这样一个key value对存放在AllRoutesChan里面。
var AllRoutesChan map[string] chan struct{}




func GetRouteChan(key string) chan struct{}{
	return AllRoutesChan[key]
}

func AddRouteChan(key string) {
	c := make(chan struct{})
	AllRoutesChan[key] = c
}

func DelRouteChan(key string) {
	delete(AllRoutesChan,key)
}

func RegisterRoute(hostIp string,port string,route config.Route) error {
	key := hostIp + ":" + port
	c := GetRouteChan(key)
	tc := beego.AppConfig.String("broadCastTimeCycle")
	t, _ := strconv.Atoi(tc)
	count := 0
	for {
		select {
		case <-c:
			Logger.Info("收到停止信号,即将停止注册", lager.Data{"route": route})
			return nil
		default:
			if count == 0 {
				if err := Msgbus.SendMessage("router.register", hostIp, route, instanceId); err != nil {
					Logger.Info("向nats发送route注册的信息失败")
					return err
				}
				fmt.Printf("向nats广播路由注册信息%s成功，%d秒之后即将再次发送注册信息\n",key,t)
				count =t
			}

			time.Sleep(time.Duration(1) * time.Second)
			count = count -1
		}
	}
}

func DeregisterRoute(hostIp string,port string,route config.Route) error{
	key := hostIp + ":" + port
	if err := Msgbus.SendMessage("router.unregister", hostIp, route, instanceId); err != nil {
		fmt.Printf("向nats发送注销路由信息%s失败\n",key)
		return err
	}else {
		fmt.Printf("向nats发送注销路由信息%s成功\n",key)
	}
	return nil
}

func Int64ToInt(a int64) int{
	return int(a)
}

func Int64ToString(a int64) string {
	return strconv.FormatInt(a,10)
}

type Res struct{
	result string
	msg string
}

func GenerateRes(result string,msg string) *Res {
	r := Res{result:result, msg:msg}
	return &r
}