$(function() {
    if (parent != self) {
        parent.location = location
    }

    function errorMsg(msg) {
        $(".error-box h1").text("登录出错：" + msg);
        $(".error-box").show();
    }
    var refreshCode = function() {
        $(this)[0].src = "/login/captcha?_=" + (Math.random() + "").substring(1)
    };
    $(".captcha").click(refreshCode);
    $(".login-button").click(function(event) {
        event.preventDefault();
        var $form = $(this).parent().parent();
        var u = $form.find("[name='username']").val();
        var p = $form.find("[name='password']").val();
        var c = $form.find("[name='captcha']").val();
        var i = $form.find("[name='api_auth_id']").val();
        if (!u) {
            errorMsg("用户名不能为空！");
            return
        }
        if (!p) {
            errorMsg("密码不能为空！");
            return
        }
        if (!c) {
            errorMsg("验证码不能为空！");
            return
        }
        $.post("", {
            "username": u,
            "password": p,
            "captcha": c,
            "api_auth_id": i,
        }, function(data) {
            if (data.code) {
                $(".error-box").hide();
                $form.slideUp(500, function() {
                    $('.welcome-text').show().addClass("animated").addClass("bounceInUp");
                    setTimeout(function() {
                        location = data.redirect.url;
                    }, data.redirect.sleep * 3);
                });
            } else {
                errorMsg(data.message);
                // alert(data.message);
                $form.find("[name='captcha']").val("");
                refreshCode();
                $(".captcha")[0].click();
            }
        }, "json");
        return
    });
    $(".form-tab li").on("click", function() {
        var index = $(this).index();
        $(this).addClass("cur").siblings().removeClass("cur");
        $(".form-con>div").hide().eq(index).show();
        if (index == 0) {
            $(".captcha")[0].src = "/login/captcha?_=" + (Math.random() + "").substring(1);
            $(".form-foot").hide();
        } else {
            $(".captcha")[1].src = "/login/captcha?_=" + (Math.random() + "").substring(1);
            $(".form-foot").show();
        }
        $(".form-error").hide();
    });
    $('form').on("keydown", function(event) {
        if (event.which == 13) {
            $(this).find(".login-button").click();
        }
    });
});