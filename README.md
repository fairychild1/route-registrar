cloudfoundry的路由注册系统
===============

允许向cloudfoundry的gorouter注册路由，已实现不在cloudfoundry平台部署应用，依然也可以使用cloudfoundry提供的负载均衡和路由分发服务。不管是虚机还是物理机，只要将ip+port+域名提交到路由注册系统里，就可以通过cloudfoundry的haproxy访问到。

# 依赖
mysql数据库环境，将用户名、密码、数据库名在conf/app.conf配置文件里。
go+beego

## 原理
本程序借助cloudfoundry的route-registrar模块，通过nats消息中间件，将要注册的路由发送到gorouter.