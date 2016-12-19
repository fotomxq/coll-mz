//之前提交的内容，避免反复提交
var oldUsername = '';
var oldPassword = '';
//当前动作类型
var nowAction = 'login-form';

//开始登录
function submitLogin() {
    //如果正在提交数据
    if (nowAction == 'action-login') {
        sendMessage('正在提交数据，请稍等片刻...', 'spinner loading');
        return false;
    }
    //获取数据
    var username = $('[name="username"]').val();
    var password = $('[name="password"]').val();
    //比对内容，如果没有改变则不提交
    if (oldUsername == username && password == oldPassword) {
        sendMessage('请不要反复提交！', 'warning circle');
        return false;
    }
    //如果没有内容，则返回
    if (username == '' || password == '') {
        sendMessage('请输入用户名和密码！', 'warning circle');
        return false;
    }
    //验证内容是否合法
    if (!checkStr(/^[a-zA-z0-9]\w{3,20}$/, username)) {
        sendMessage('请输入正确的用户名！', 'warning circle');
        return false;
    }
    if (!checkStr(/^[a-zA-z0-9]\w{3,20}$/, password)) {
        sendMessage('请输入正确的密码！', 'warning circle');
        return false;
    }
    //计算password sha1
    var passwordSha1 = hex_sha1(password);
    //正在提交提示
    sendMessage('等待服务器响应。', 'spinner loading');
    //覆盖旧的用户和密码信息，避免重复提交
    oldUsername = username;
    oldPassword = password;
    nowAction = 'action-login';
    //提交到服务器，等待确认
    $.post('/action-login', {
        'username': username,
        'password': passwordSha1
    }, function(data) {
        nowAction = 'login-form';
        if (!data) {
            sendMessage('服务器没有响应！', 'ban');
            return false;
        }
        if (!data['status']) {
            sendMessage('服务器没有响应！', 'ban');
            return false;
        }
        if (data['data'] == 'success') {
            sendMessage('登录成功，正在跳转！', 'alarm outline');
            window.location.href = '/center';
            return false;
        } else {
            sendMessage('登录失败，请检查用户名或密码是否正确？', 'ban');
            return false;
        }
    }, 'json');
}

//发送消息
function sendMessage(msg, ico) {
    $('#message').dimmer('show');
    $('#message-content').html('<i class="' + ico + ' big icon"></i> ' + msg);
}

//验证正则表达式
function checkStr(m, c) {
    var reg = new RegExp(m);
    return reg.test(c)
}

//登陆组件控制
$(document).ready(function() {
    //按钮提交表单
    $('a[href="#login"]').click(function() {
        submitLogin();
    });
    //回车提交表单
    $('body').keyup(function(event) {
        if (event.keyCode === 13) {
            submitLogin();
        }
    });
});
