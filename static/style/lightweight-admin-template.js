

function initNavEvent(){
	var itemchild = $('.nav__item--child');
	itemchild.on('click', function(){
        var target = $(this);

        $(target).stop(true, true);
        $(target).siblings().removeClass('nav__item--active');
        $(target).addClass('nav__item--active');

        //un-comment if you won't redirect to other url
        //return false;
    });

    var item = $('.nav__item');
    item.on('click', function(){
        var target = $(this);

        $(target).addClass('nav__item--active');
		var siblings = $(target).siblings();

		$.each(siblings, function(i, el){
			if($(el).hasClass('nav__item--has-child')){
				$(el).find('.nav__wrapper--child').slideUp();
			}
	        $(el).removeClass('nav__item--active');
	    });        

        if($(target).hasClass('nav__item--has-child')){
        	$(target).find('.nav__item--child').removeClass('nav__item--active');
        	$(target).find('.nav__arrow').toggleClass('nav__arrow--active');        	
        	$(target).find('.nav__wrapper--child').slideToggle();
        }

        //un-comment if you won't redirect to other url
        //return false;
    });
}


function initAccountPopover(){
	var item = $('.account--has-login');
	item.hover(function(){
		var menu = $(this).find('.account__menu');		
        $(menu).stop(true, true).slideDown();
    }, function(){
        var menu = $(this).find('.account__menu');
        $(menu).stop(true, true).delay(1000).slideUp();
    });

    var menuEl = $('.account__menu');
    menuEl.hover(function(){
        $(this).stop(true, true);
    }, function(){
        $(this).stop(true, true).delay(1000).slideUp();
    });
}


function initOffCanvasMenu(){
    $('.header__nav-btn').on('click', function(){
        $(this).toggleClass('header__nav-btn--close');
        $('.nav').toggleClass('nav--opened');
    });
}

$(document).ready(function () {
	'use strict';   
	initNavEvent();
	initAccountPopover();
    initOffCanvasMenu();
});

$(window).load(function(){
    'use strict';    
});