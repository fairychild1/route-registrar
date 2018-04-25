package controllers

import (
	"github.com/astaxie/beego"
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Login() {
	req :=map[string]interface{}{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &req); err != nil {
		this.Ctx.Output.SetStatus(400)
		fmt.Println(this.Ctx.Input.RequestBody)
		this.Ctx.Output.Body([]byte("没有指定参数!"))
		return
	}
	user := req["client_id"].(string)
	passwd := req["client_secret"].(string)
	t := req["grant_type"].(string)

	//向ucenter发送用户名和密码，验证是否正确
	uc := beego.AppConfig.String("usercenter")
	url := "http://"+uc+"/oauth/token?"+"username="+user+"&"+"password="+passwd+"&"+"grant_type="+t
	fmt.Println("url is:",url)
	r, _ := http.NewRequest("POST", url,nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Authorization", "Basic Y2Y6")
	r.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp,err := client.Do(r)
	//defer resp.Body.Close()
	if err != nil {
		fmt.Println("发送用户名和密码的验证请求失败.")
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte("不能正常解析参数，请检查参数格式"))
		return
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		tempStru :=  map[string] interface{}{}
		if err = json.Unmarshal(body, &tempStru); err != nil {
			fmt.Println(err)
			this.Ctx.Output.SetStatus(400)
			this.Ctx.Output.Body([]byte("ucenter返回的数据不能正常解析"))
			return
		}

		this.Data["json"]=tempStru
		this.ServeJSON()
		resp.Body.Close()
		return
	}
}