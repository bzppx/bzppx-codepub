$(function() {
    var refreshCode = function() {
        $(".captcha")[0].src = "/login/captcha?_=" + (Math.random() + "").substring(1)
    }
    $(".captcha").click(refreshCode);
    $("#login-button").click(function(event) {
        event.preventDefault();
        var u = $("[name='username']").val();
        var p = $("[name='password']").val();
        var c = $("[name='captcha']").val();
        if (!u || !p || !c) {
            return
        }
        $.post("", {
            "username": u,
            "password": p,
            "captcha": c,
        }, function(data) {
            if (data.code) {
                $('form').fadeOut(500, function() {
                    setTimeout(function() {
                        location = data.redirect.url;
                    }, data.redirect.sleep);
                });
                $('.wrapper').addClass('form-success');
            } else {
                alert(data.message)
                $("[name='captcha']").val("");
                refreshCode();
            }
        }, "json");
        return
    });
});