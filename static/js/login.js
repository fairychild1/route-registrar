/**
 * Created by LENOVO on 2018/3/2.
 */

$(function (){
    var a=$('form.form-horizontal button');
    //var a=$('#login_button');
    a.click(function () {
        /* alert('login in.');*/
        var jqxhr= $.ajax('/login',{
            /*dataType: 'json',*/
            type: 'post',
            data: JSON.stringify(
                {
                    "client_id":$('#username').val(),
                    "client_secret":$('#inputPassword').val(),
                    "grant_type":"password"
                }
            ),
            contentType: "application/json; charset=utf-8",
            success: function(data) {
                if (data.error_description !== undefined) {
                    alert(data.error_description)
                }
                if (data.access_token === undefined) {
                    alert("用户名或密码错误")
                }else {
                    localStorage.setItem('access_token',data.access_token)  //将access_token存入web浏览器的localStorage
                    to_list_page(data.access_token)
                }
            },
            error: function(err) {
                alert("err");
            }
        });
    })
});

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