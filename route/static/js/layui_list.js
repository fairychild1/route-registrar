layui.use(['jquery', 'layer','form'], function() {
    var $ = layui.$
    var layer = layui.layer
    var form = layui.form
    var list = JSON.parse(localStorage.getItem("list_result"))

    $(document).ready(function (){
//console.log("list is"+list[0]['HostIp'])
        if (list === null) {
        }else {
            for (var i=0;i<list.length;i++){
                console.log("layui\n");
                var str = "<tr><td>"+list[i]['RouteName']+"</td>" +
                    "<td>"+list[i]['HostIp']+"</td>" +
                    "<td>"+list[i]['Port']+"</td>" +
                    "<td>"+list[i]['Uri']+"</td>" +
                    "<td><div><button id='update_"+i+"' class='btn' type='button' onclick='update_route("+i+")'>修改</button><button id='delete_"+i+"' class='btn delete-button' type='button' onclick='delete_route("+i+")'>删除</button></div></td>"+
                    "</tr>"
                console.log(str);
                $("#nav-tbody").append(str);

                //给update按钮和delete按钮绑定事件
                var s="delete_"+i
                var a=$('#'+s)
                //a.click(del_route(list[i]['Id'],localStorage.getItem("access_token")))

            }
        }


        var str = "<form id='tab_route' class='layui-form' action=''>"+
            "<div class='layui-form-item'>"+
            "<label class='layui-form-label'>名称</label>"+
            "<div class='layui-input-block'>"+
            "<input id='tab_route_name' type='text' name='title' required  lay-verify='required' placeholder='例如:route-test' autocomplete='off' class='layui-input'>"+
            "</div>"+
            "</div>"+
            "<div class='layui-form-item'>"+
            "<label class='layui-form-label'>域名</label>"+
            "<div class='layui-input-block'>"+
            "<input id='tab_route_url' type='text' name='title' required  lay-verify='required' placeholder='例如：www.sina.com' autocomplete='off' class='layui-input'>"+
            "</div>"+
            "</div>"+
            "<div class='layui-form-item'>"+
            "<label class='layui-form-label'>ip</label>"+
            "<div class='layui-input-block'>"+
            "<input id='tab_route_ip' type='text' name='title' required  lay-verify='required' placeholder='可以是多个，以冒号分隔，并与端口顺序对应' autocomplete='off' class='layui-input'>"+
            "</div>"+
            "</div>"+
            "<div class='layui-form-item'>"+
            "<label class='layui-form-label'>port</label>"+
            "<div class='layui-input-block'>"+
            "<input id='tab_route_port' type='text' name='title' required  lay-verify='required' placeholder='可以是多个，以冒号分隔，不能是80，8080' autocomplete='off' class='layui-input'>"+
            "</div>"+
            "</div>"+
            "</form>"
        $("#add_route").click(function(){
            layer.open({
                type: 1,
                area: ['500px', '300px'],
                content: str,
                btn:['提交','取消'],
                btn1: function(index,layero) {
                  add_route()
                },
                btn2: function(index,layero) {
                },
            })

        })
        $("#exit_route_registrar").click(function(){
            localStorage.clear()
            window.location.href = '/';
        })



    })

})

function delete_route(i){
    var list = JSON.parse(localStorage.getItem("list_result"))
    del_ajax(list[i]['Id'],localStorage.getItem("access_token"));
}

function update_route(i) {
    layui.use(['jquery', 'layer','form'], function() {
        var $ = layui.$
        var layer = layui.layer
        var form = layui.form
        var list = JSON.parse(localStorage.getItem("list_result"))
        var id = list[i]['Id']
        var name = list[i]['RouteName']
        var ip = list[i]['HostIp']
        var port = list[i]['Port']
        var url = list[i]['Uri']
        var str = "<form id='tab_route' class='layui-form' action=''>"+
            "<div class='layui-form-item'>"+
            "<label class='layui-form-label'>名称</label>"+
            "<div class='layui-input-block'>"+
            "<input id='tab_route_name' type='text' name='title' required  lay-verify='required' placeholder='例如:route-test' autocomplete='off' class='layui-input' value='"+name+"'>"+
            "</div>"+
            "</div>"+
            "<div class='layui-form-item'>"+
            "<label class='layui-form-label'>域名</label>"+
            "<div class='layui-input-block'>"+
            "<input id='tab_route_url' type='text' name='title' required  lay-verify='required' placeholder='例如：www.sina.com' autocomplete='off' class='layui-input' value='"+url+"'>"+
            "</div>"+
            "</div>"+
            "<div class='layui-form-item'>"+
            "<label class='layui-form-label'>ip</label>"+
            "<div class='layui-input-block'>"+
            "<input id='tab_route_ip' type='text' name='title' required  lay-verify='required' placeholder='可以是多个，以冒号分隔，并与端口顺序对应' autocomplete='off' class='layui-input' value='"+ip+"'>"+
            "</div>"+
            "</div>"+
            "<div class='layui-form-item'>"+
            "<label class='layui-form-label'>port</label>"+
            "<div class='layui-input-block'>"+
            "<input id='tab_route_port' type='text' name='title' required  lay-verify='required' placeholder='可以是多个，以冒号分隔，不能是80，8080' autocomplete='off' class='layui-input' value='"+port+"'>"+
            "</div>"+
            "</div>"+
            "</form>"
        layer.open({
            type: 1,
            area: ['500px', '300px'],
            content: str,
            btn:['提交','取消'],
            btn1: function(index,layero) {
                ajax_update_route(id)
            },
            btn2: function(index,layero) {
                console.log("提交的ip是"+$("#tab_route_ip").val())
            },
        })

    })
}

function ajax_update_route(id) {

    var name=$("#tab_route_name").val()
    var ip=$("#tab_route_ip").val()
    var port=$("#tab_route_port").val()
    var url=$("#tab_route_url").val()
    var auth=localStorage.getItem("access_token")
    var i=String(id)
    var a = $.ajax('/v1/route/',{
        type: 'put',
        data: JSON.stringify(
            {
                "id":i,
                "route_name":name,
                "port":port,
                "uri":url,
                "host_ip":ip
            }
        ),
        headers:{"Authorization":auth},
        contentType: "application/json; charset=utf-8",
        success: function(data) {
            //localStorage.setItem('list_result',data)
            alert("修改成功");
            to_list_page(auth);
            //window.location.href = '/views/list.html';
            //$("ss").innerHTML = data;

        },
        error: function(err) {
            alert(err.responseText);
        }
    });
}

function add_route() {
    var name=$("#tab_route_name").val()
    var ip=$("#tab_route_ip").val()
    var port=$("#tab_route_port").val()
    var url=$("#tab_route_url").val()
    var auth=localStorage.getItem("access_token")
    var a = $.ajax('/v1/route_register/',{
        type: 'post',
        data: JSON.stringify(
            {
                "route_name":name,
                "port":port,
                "uri":url,
                "host_ip":ip
            }
        ),
        headers:{"Authorization":auth},
        contentType: "application/json; charset=utf-8",
        success: function(data) {
            //localStorage.setItem('list_result',data)
            alert("添加成功");
            to_list_page(auth);
            //window.location.href = '/views/list.html';
            //$("ss").innerHTML = data;

        },
        error: function(err) {
            alert(err.responseText);
        }
    });
}


function to_list_page(auth) {
    var a = $.ajax('/v1/list',{
        type: 'get',
        headers:{"Authorization":auth},
        contentType: "application/json; charset=utf-8",
        success: function(data) {
            //localStorage.setItem('list_result',data)
            localStorage.setItem("list_result",JSON.stringify(data));
            window.location.href = '/views/list.html';
            //$("ss").innerHTML = data;

        },
        error: function(err) {
            alert(err.responseText);
        }
    });
}

function del_ajax(id,auth) {
    var a = $.ajax('/v1/route_deregister/'+id,{
        type: 'delete',
        headers:{"Authorization":auth},
        //contentType: "application/json; charset=utf-8",
        success: function(data) {
            //localStorage.setItem('list_result',data)
            alert("删除成功");
            to_list_page(auth);
            //window.location.href = '/views/list.html';
            //$("ss").innerHTML = data;

        },
        error: function(err) {
            alert(err.responseText);
        }
    });
}
