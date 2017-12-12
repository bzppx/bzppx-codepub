 $("#login-button").click(function(event){
		 event.preventDefault();
	 
	 $('form').fadeOut(500,function(){
		 setTimeout(function(){
			  location="index.html";
		 },500);
	 });
	 $('.wrapper').addClass('form-success');
	
});