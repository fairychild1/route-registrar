package controllers

import (
	"strconv"
	"fmt"
	"strings"
	"github.com/astaxie/beego"
	"encoding/json"
	//"code.cloudfoundry.org/lager"
	//"code.cloudfoundry.org/lager/lagerflags"
	"route-registrar/config"
	"route-registrar/utility"
	"route-registrar/models"
	"route-registrar/token"
	"net"
)

type RouteController struct {
	beego.Controller
}

func (this *RouteController) CheckUri() {
	req :=map[string]interface{}{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		fmt.Println(this.Ctx.Input.RequestBody)
		this.Ctx.Output.Body([]byte("can't use the kong interface!"))
		return
	}
}

func (this *RouteController) RouteRegister() {
	req :=map[string]interface{}{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("不能正常解析参数，请检查参数格式"))
		return
	}
	//提取用户名
	auth := this.Ctx.Request.Header["Authorization"]  //返回的是[]string类型
	username,_,_ := token.GetUserFromToken(auth[0])


	routeName := req["route_name"].(string)
	port := req["port"].(string)
	uri := req["uri"].(string)
	hostIp := req["host_ip"].(string)

	//检查参数是否合法
	if err := checkParam(routeName,port,uri,hostIp);err !=nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}

	//判断路由是否已经存在
	if flag,err := models.RouteInApp(uri);err !=nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}else if flag ==true {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("路由已经存在，被应用注册过了"))
		return
	}

	if flag,err := models.RouteInUserSet(uri);err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}else if flag ==true {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("路由已经存在，被用户注册过了"))
		return
	}

	if err := AddRoutes(routeName,port,uri,hostIp);err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("传入参数错误，hostip和port必须以分号分割，并且数量要一致"))
		return
	}

	//p,err := strconv.Atoi(port)
	//if err != nil {
	//	this.Ctx.Output.SetStatus(400)
	//	this.Ctx.Output.Body([]byte("port是用引号括起来的整数，请确认"))
	//	return
	//}
	//route := NewRoute(routeName,p,uri)
//
	////生成chan，并启动协程，注册route
	//utility.AddRouteChan(hostIp+":"+port)
	//go utility.RegisterRoute(hostIp,port,*route)

	//将注册信息写到数据库里面
	r := models.Route{RouteName:routeName,Port:port,Uri:uri,HostIp:hostIp,UserName:username}
	if err := models.AddRoute(&r);err !=nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("将路由的信息写到数据库失败"))
		return
	}

	this.Ctx.Output.SetStatus(200)
	this.Ctx.Output.Body([]byte("route注册成功"))
}

func AddRoutes(routename string,port string,uri string,hostip string) error{
	ha := strings.Split(hostip,";")
	pa := strings.Split(port,";")
	if (len(ha) != len(pa) ){
		fmt.Println("传入参数错误，hostip和port必须以分号分割，并且数量要一致")
		return fmt.Errorf("%s", "传入参数错误，hostip和port必须以分号分割，并且数量要一致")
	}
	for k,v := range ha {
		hostip := v
		port := pa[k]
		p,err := strconv.Atoi(port)
		if err != nil {
			fmt.Println("port是整数，请确认")
			return err
		}
		route := NewRoute(routename,p,uri)
		//生成chan，并启动协程，注册route
		utility.AddRouteChan(hostip+":"+port)
		go utility.RegisterRoute(hostip,port,*route)
	}
	return nil
}

func NewRoute(name string,port int,uri string) *config.Route{
	a := []string{}
	a = append(a,uri)
	return &config.Route{
		Name : name,
		Port : &port,
		URIs : a,
	}
}

func (this *RouteController) RouteDeregister() {
	i := this.Ctx.Input.Param(":route_id")
	id,err :=strconv.Atoi(i)
	if err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(fmt.Sprintf("请求中带的id是无效的整数")))
		return
	}
	err = delRoute(id)
	if err !=nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(fmt.Sprintf("数据库中找不到id %对应的route记录",id)))
	}

	if err = models.DelRoute(id);err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(fmt.Sprintf("删除id %d对应的路由失败",id)))
		return
	}

	this.Ctx.Output.SetStatus(200)
	this.Ctx.Output.Body([]byte(fmt.Sprintf("删除id %d对应的路由成功",id)))
}

func (this *RouteController) ListRoutes() {
	//提取用户名
	auth := this.Ctx.Request.Header["Authorization"]  //返回的是[]string类型
	username,_,_ := token.GetUserFromToken(auth[0])
	fmt.Printf("token中提取的用户名是%s\n",username)
	pr,err := models.ListRoutes(username)
	if err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(fmt.Sprintf("从数据库查询%s对应的路由失败\n",username)))

	}
	fmt.Printf("%s总共注册了%d条路由\n",username,len(*pr))

	this.Ctx.Output.SetStatus(200)
	this.Data["json"] = *pr
	this.ServeJSON()
	//this.TplName = "list.html"
	//this.Render()
	return
}

func (this *RouteController) GetRoute() {

}

func (this *RouteController) UpdateRoute() {
	req :=map[string]interface{}{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("不能正常解析参数，请检查参数格式"))
		return
	}

	//提取用户名
	auth := this.Ctx.Request.Header["Authorization"]  //返回的是[]string类型
	username,_,_ := token.GetUserFromToken(auth[0])

	routeName := req["route_name"].(string)
	port := req["port"].(string)
	uri := req["uri"].(string)
	hostIp := req["host_ip"].(string)
	i := req["id"].(string)
	id,err := strconv.Atoi(i)
	if err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("id是用引号括起来的整数，请确认"))
		return
	}

	//检查参数是否合法
	if err = checkParam(routeName,port,uri,hostIp);err !=nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	//判断路由是否已经存在
	if f := models.IfChangeuri(id,uri);f == true {
		if flag,err := models.RouteInApp(uri);err !=nil {
			this.Ctx.Output.SetStatus(400)
			this.Ctx.Output.Body([]byte(err.Error()))
			return
		}else if flag ==true {
			this.Ctx.Output.SetStatus(400)
			this.Ctx.Output.Body([]byte("路由已经存在，被应用注册过了"))
			return
		}

		if flag,err := models.RouteInUserSet(uri);err != nil {
			this.Ctx.Output.SetStatus(400)
			this.Ctx.Output.Body([]byte(err.Error()))
			return
		}else if flag ==true {
			this.Ctx.Output.SetStatus(400)
			this.Ctx.Output.Body([]byte("路由已经存在，被用户注册过了"))
			return
		}
	}

	//先停止该路由的注册，并从数据库删除相关记录
	err = delRoute(id)
	if err !=nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(fmt.Sprintf("数据库中找不到id %对应的route记录",id)))
		return
	}

	if err = models.DelRoute(id);err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(fmt.Sprintf("删除id %d对应的路由失败",id)))
		return
	}


	//重新注册路由

	if err := AddRoutes(routeName,port,uri,hostIp);err != nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("传入参数错误，hostip和port必须以分号分割，并且数量要一致"))
		return
	}

	//p,err := strconv.Atoi(port)
	//if err != nil {
	//	this.Ctx.Output.SetStatus(400)
	//	this.Ctx.Output.Body([]byte("port是用引号括起来的整数，请确认"))
	//	return
	//}
//
	////重新启动一个注册route的协程
	////生成chan，并启动协程，注册route
	//route := NewRoute(routeName,p,uri)
	//utility.AddRouteChan(hostIp+":"+port)
	//go utility.RegisterRoute(hostIp,port,*route)

	//更新数据库的route记录,将注册信息写到数据库里面
	r := models.Route{RouteName:routeName,Port:port,Uri:uri,HostIp:hostIp,UserName:username}
	if err := models.AddRoute(&r);err !=nil {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("将更新后的路由的信息写到数据库失败"))
		return
	}

	this.Ctx.Output.SetStatus(200)
	this.Ctx.Output.Body([]byte(fmt.Sprintf("更新id %d对应的路由成功",id)))
}

func delRoute(id int) error{
	//根据id从数据库查找对应的route记录
	r,err := models.GetRouteById(id)
	if err != nil {
		fmt.Printf("查找id %d对应的route记录失败\n",id)
		return err
	}

	err = setIpPortArrayByDb(r.RouteName,r.Port,r.Uri,r.HostIp)
	return err

	////向nats发送消息，注销route
	//p,err := strconv.Atoi(r.Port)
	//route := NewRoute(r.RouteName,p,r.Uri)
//
	//key := r.HostIp+":"+r.Port
	//c := utility.GetRouteChan(key)
	//c <- struct{}{} //停止注册route的协程
	//utility.DelRouteChan(key)
	//err = utility.DeregisterRoute(r.HostIp,r.Port,route)    //向nats发送注销route的消息
	//if err != nil {
	//	fmt.Printf("向nats发送注销%s对应的路由失败",key)
	//	return err
	//}
	//return nil
}

func setIpPortArrayByDb(routename string,port string,uri string,hostip string) error {
	ha := strings.Split(hostip,";")
	pa := strings.Split(port,";")
	if (len(ha) != len(pa) ){
		fmt.Println("传入参数错误，hostip和port必须以分号分割，并且数量要一致")
		return fmt.Errorf("%s", "传入参数错误，hostip和port必须以分号分割，并且数量要一致")
	}
	for k,ip := range ha {
		//向nats发送消息，注销route
		fmt.Printf("要删除的route的ip是%s port是%s\n",ip,pa[k])
		p,err := strconv.Atoi(pa[k])
		if err != nil {
			fmt.Printf("不能把%s转换成整数",pa[k])
		}
		route := NewRoute(routename,p,uri)
		key := ip+":"+pa[k]
		c := utility.GetRouteChan(key)
		c <- struct{}{} //停止注册route的协程
		utility.DelRouteChan(key)
		err = utility.DeregisterRoute(ip,pa[k],*route)    //向nats发送注销route的消息
		if err != nil {
			fmt.Printf("向nats发送注销%s对应的路由失败",key)
			return err
		}
	}
	return nil

}

func checkParam(routename string,port string,uri string,hostip string) error{
	ha := strings.Split(hostip,";")
	pa := strings.Split(port,";")
	if (len(ha) != len(pa) ){
		return fmt.Errorf("%s", "传入参数错误，hostip和port必须以分号分割，并且数量要一致")
	}

	for _,v := range pa {
		if v == "80" || v == "8080"{
			return fmt.Errorf("%s", "端口不能为80或者8080")
		}
		i,err := strconv.Atoi(v)
		if err !=nil {
			return fmt.Errorf("%s", "端口必须为整数")
		}
		if i<1 || i> 65535 {
			return fmt.Errorf("%s", "端口必须在1到65535之间")
		}
	}

	for _,v := range ha {
		if r :=net.ParseIP(v);r == nil {
			return fmt.Errorf("%s", "ip是非法的，请检查ip信息")
		}
	}
	return nil

}