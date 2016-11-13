$(document)
    .ready(function() {
      $('.ui.form')
        .form({
          fields: {
            email: {
              identifier  : 'email',
              rules: [
                {
                  type   : 'empty',
                  prompt : '请输入您的用户名'
                },
                {
                  type   : 'email',
                  prompt : '请输入合适的email'
                }
              ]
            },
            password: {
              identifier  : 'password',
              rules: [
                {
                  type   : 'empty',
                  prompt : '密码不能留空'
                },
                {
                  type   : 'length[6]',
                  prompt : '密码必须大于6位数'
                }
              ]
            }
          }
        })
      ;
    })
  ;