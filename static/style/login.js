 $("#login-button").click(function(event) {
     event.preventDefault();
     var u = $("[name='username']").val();
     var p = $("[name='password']").val();
     if (!u || !p) {
         return
     }
     $.post("", {
         "username": u,
         "password": p,
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
         }
     }, "json");
     return


 });