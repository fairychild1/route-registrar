package models

import(
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"time"
	"github.com/astaxie/beego"
)

func init() {
	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpassword")
	db := beego.AppConfig.String("db")
	conn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/"+db+"?charset=utf8"
	//fmt.Println("传入参数是:"+conn)
	//orm.RegisterDataBase("default", "mysql", "root:admin@tcp(10.0.0.46:3306)/test?charset=utf8")
	orm.RegisterDataBase("default", "mysql", conn, 30)

	// register model
	orm.RegisterModel(new(Route))

	// create table
	orm.RunSyncdb("default", false, true)
}

type Route struct {
	Id int
	UserName string `orm:"size(50)"`
	RouteName string `orm:"size(50)"`
	Port string `orm:"size(120)"`
	//Port int
	Uri string `orm:"size(100)"`
	HostIp string `orm:"size(320)"`
	Created time.Time `orm:"auto_now;type(datetime)"`
}

//添加route
func AddRoute(route *Route) error {
	o := orm.NewOrm()
	_,err := o.Insert(route)
	if err !=nil {
		fmt.Println("添加route的信息到数据库失败")
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//更新route
func UpdateRoute(route *Route) error {
	o := orm.NewOrm()
	_, err := o.Update(route)
	if err != nil {
		fmt.Println("在数据库更新route的信息失败")
		return err
	}
	return nil
}

//删除route
func DelRoute(id int) error{
	o := orm.NewOrm()
	_, err := o.Delete(&Route{Id: id})
	if err != nil {
		fmt.Println("在数据库删除route的信息失败")
	}
	return nil
}

//根据UserName查询对应的route的所有记录信息
func ListRoutes(username string) (*[]orm.Params,error) {
	o :=orm.NewOrm()
	var routes []orm.Params
	_, err := o.QueryTable("route").Filter("user_name", username).Values(&routes)
	return &routes,err
}

//根据id查找route的记录
func GetRouteById(id int) (*Route,error){
	o :=orm.NewOrm()
	route :=Route{Id : id}
	err := o.Read(&route)
	if err != nil {
		if err == orm.ErrNoRows {
			fmt.Printf("在数据库里，查询不到id %d对应的记录\n",id)
		} else if err == orm.ErrMissPK {
			fmt.Printf("id %d是无效的主键\n",id)
		}else{
			fmt.Printf("根据id %d查询route的记录，发现未知的数据库访问错误\n",id)
		}
		return &route,err
	}else {
		fmt.Printf("查询到id %d对应的数据库记录\n",id)
		return &route,err
	}
}

//是否修改了路由,如果修改了，返回true，否则返回false
func IfChangeuri(id int,uri string) bool {
	if r,err := GetRouteById(id);err != nil {
		return false
	}else {
		if r.Uri == uri {
			return false
		}else {
			return true
		}
	}
}

func RouteInUserSet(uri string) (bool,error) {
	o :=orm.NewOrm()
	var routes []orm.Params
	_, err := o.QueryTable("route").Values(&routes)
	if err != nil {
		fmt.Printf("查询数据库的路由信息失败\n")
		return false,err
	}
	for _,v := range routes {
		fmt.Printf("uri是%s\n",v["Uri"])
		if v["Uri"] == uri {
			return true,nil
		}
	}
	return false,nil
}

//加载所有的route信息
func LoadAllRoutes() (*[]orm.Params,error) {
	o :=orm.NewOrm()
	var routes []orm.Params
	_, err := o.QueryTable("route").Values(&routes)
	return &routes,err
}