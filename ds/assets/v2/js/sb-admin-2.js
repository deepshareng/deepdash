$(function() {

    $('#side-menu').metisMenu();

});

//Loads the correct sidebar on window load,
//collapses the sidebar on window resize.
// Sets the min-height of #page-wrapper to window size
$(function() {
    $(window).bind("load resize", function() {
        var topOffset = 50;
        var width = (this.window.innerWidth > 0) ? this.window.innerWidth : this.screen.width;
        if (width < 768) {
            $('div.navbar-collapse').addClass('collapse');
            topOffset = 100; // 2-row-menu
        } else {
            $('div.navbar-collapse').removeClass('collapse');
        }

        var height = ((this.window.innerHeight > 0) ? this.window.innerHeight : this.screen.height) - 1;
        height = height - topOffset;
        if (height < 1) height = 1;
        if (height > topOffset) {
            $("#page-wrapper").css("min-height", (height) + "px");
        }
    });

    var url = window.location.href.toString();
    $('.deepdash-sidebar ul.nav a').filter(function() {
        return url.indexOf(this.href) >= 0;
    }).addClass('active');

    $(window).on("scroll", function(){
        var $headerContainer = $(".ds-integrate-steps-header-container");
        var $headerContainerFixed = $(".ds-integrate-steps-header-container-fixed");
        if ($headerContainer.length > 0 && $headerContainerFixed.length > 0) {
            var headerOffset = $headerContainer.offset();
            if($(window).scrollTop() > headerOffset.top) {
                $headerContainerFixed.addClass("fixed-top");
            }else{
                $headerContainerFixed.removeClass("fixed-top");
            }
        }
    });
});
