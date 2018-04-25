package models
import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
)

func init() {
	orm.RegisterDriver("postgres", orm.DRPostgres)
	//orm.RegisterDataBase("default", "postgres", "ccadmin:c1oudc0w@tcp(192.168.1.178:5524)/ccdb?charset=utf8")
	ccdbname := beego.AppConfig.String("ccdbname")
	ccuser := beego.AppConfig.String("ccuser")
	ccpasswd := beego.AppConfig.String("ccpasswd")
	cchost := beego.AppConfig.String("cchost")
	ccport := beego.AppConfig.String("ccport")
	conn := "user="+ccuser+" password="+ccpasswd+" host="+cchost+" port="+ccport+" dbname="+ccdbname+" sslmode=disable"
	orm.RegisterDataBase("appdb","postgres",conn)
}

func GetAllAppRoutes() *[]orm.Params{
	o := orm.NewOrm()
	o.Using("appdb") // 默认使用 default，你可以指定为其他数据库

	var maps []orm.Params
	_,err := o.Raw("select r.host||'.'||d.name route from routes r,domains d where r.domain_id = d.id").Values(&maps)
	if err == nil  {
		return &maps
	}else {
		fmt.Println(err.Error())
		return nil
	}
}

func RouteInApp(uri string) (bool,error) {
	m := GetAllAppRoutes()
	if m == nil {
		return true,fmt.Errorf("%s", "查询cc的数据库失败")
	}
	maps := *m
	for k,_ := range maps {
		if maps[k]["route"] == uri {
			return true,nil
		}
	}
	return false,nil
}


