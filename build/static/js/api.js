// 发送ajax请求获取数据
import 'jquery-3.6.js';
import 'jquery.cookie.js';
// 创建制定的axjx对象
(function(){
    function ax(option){
        return $.ajax({
            url: option.url,
            method: option.method || "post",
            data:option.data || {},
            headers:{
                "Authorization":$.cookie("token"),
            }
        })
    }
})
