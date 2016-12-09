$(document).ready(function() {
    $('.submit').click(function() {
        $('form').submit()
    });
    $('.ui.form').form({
        fields: {
            email: {
                identifier: 'email',
                rules: [{
                    type: 'empty',
                    prompt: '请输入用户名'
                }, {
                    type: 'email',
                    prompt: '用户名输入有误，必须是email'
                }]
            },
            password: {
                identifier: 'password',
                rules: [{
                    type: 'empty',
                    prompt: '请输入密码'
                }, {
                    type: 'length[6]',
                    prompt: '密码必须大于6位'
                }]
            }
        }
    });
});
