var list = JSON.parse(localStorage.getItem("list_result"))
$(document).ready(function (){
//console.log("list is"+list[0]['HostIp'])
    for (var i=0;i<list.length;i++){
        console.log("hello,world\n");
        var str = "<tr><td>"+list[i]['RouteName']+"</td>" +
            "<td>"+list[i]['HostIp']+"</td>" +
            "<td>"+list[i]['Port']+"</td>" +
            "<td>"+list[i]['Uri']+"</td>" +
            "<td><div><button id='update_"+i+"' class='btn' type='button'>修改</button><button id='delete_"+i+"' class='btn delete-button' type='button' onclick='delete_route("+i+")'>删除</button></div></td>"+
            "</tr>"
        console.log(str);
        $("#nav-tbody").append(str);

        //给update按钮和delete按钮绑定事件
        var s="delete_"+i
        var a=$('#'+s)
        //a.click(del_route(list[i]['Id'],localStorage.getItem("access_token")))

    }

})


function delete_route(i){
    del_ajax(list[i]['Id'],localStorage.getItem("access_token"));
}

function add_route() {
    var name=$("#tab_route_name").val()
    var ip=$("#tab_route_ip").val()
    var port=$("#tab_route_port").val()
    var url=$("#tab_route_url").val()
    var auth=localStorage.getItem("access_token")
    console.log("name is",name)
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
            alert(err);
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
            alert(err);
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
            alert(err);
        }
    });
}
