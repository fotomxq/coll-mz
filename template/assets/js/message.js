//页面必须设定唯一的message作为组件
//调用函数即可使用
//param title string 标题
//param message string 消息内容
//param type string 消息类型，eg:warning
function sendMessage(title,message,type){
    $('#message').attr('class','ui message '+type);
    $('#message .header').html(title);
    $('#message .content').html(message);
    $('#message').show();
    setTimeout(function(){
        $('#message').hide();
    },10000)
}

function sendMessageBool(data,trueTitle,trueMessage,falseTitle,falseMessage){
    if(!data){
        sendMessage("无响应","无法连接到服务器。","warning")
        return false;
    }
    if(data['status']){
        sendMessage(trueTitle,trueMessage,"positive")
    }else{
        sendMessage(falseTitle,falseMessage,"negative")
    }
}