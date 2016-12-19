//初始化
$(document).ready(function() {
    //更新消息列队
    postStatus = $('.main').attr('data-status')
    if(postStatus != 'no') {
        if (postStatus == "error" || postStatus == "has") {
            sendMessage("失败", "修改用户信息失败，请检查您提交的昵称或密码是否超过限定。", "negative");
        } else if (postStatus == "ok"){
            sendMessage("成功", "修改用户昵称和密码成功", "positive");
        }
    }
    //提交按钮
    $('a[href="#edit-ok"]').click(function(){
        //确保有内容
        nicename = $('[name="nicename"]').val();
        password = $('[name="password"]').val();
        if(!nicename || !password){
            sendMessage("填写不够完整", "请填写完整的信息，之后再尝试提交。", "negative");
            return false;
        }
        //提交前换算password hex
        passwordSha1 = hex_sha1(password);
        $('[name="password"]').val(passwordSha1);
        $('.form').submit();
    });
});