(function ($) {
    "use strict" // Start of use strict

    function setup() {
        selectPageSection("dashboard")


        $('#lightswitch').click(function () {
            sendLight()
        })
    }

    setup()

    $.getJSON("/api", function (data) {
        console.log("Json resp:", data)
        window.Profile = data
        console.log(JSON.stringify(data, " ", " "))

        onProfileLoaded()
    })

    function onLightButtonClick() {
        console.log("ok")
    }

    function sendCommand(cmd) {
        $.getJSON("/api?cmd=" + cmd)
    }

    var high = true

    function sendLight() {
        high = !high
        var cmd = "DW 2 " + (high ? "LOW" : "HIGH")
        console.log("Sending:", cmd)
        sendCommand(cmd)


    }

    function onProfileLoaded() {
        var ul = $('<ul class="list-group"/>')
        Profile.devices.forEach(function (device) {
            console.log("device:", device)
            ul.append(
                '<li class="list-group-item">' + device.name + '</li>\n'
            )
        })
        $('.page-section-devices').append(ul)
    }

    function selectPageSection(section, collapse) {
        $('.page-section').hide()
        var sec = $('.page-section-' + section).show()
        var tit = $('li[data-nav-target=' + section + ']').text()
        $('.navbar-brand').text(tit)


        if (collapse) {
            $('.navbar-toggler').addClass('collapsed')
            $('#navbarResponsive').removeClass("show")
            $('.navbar-toggler').attr('aria-expanded', false)
            $('.tooltip.navbar-sidenav-tooltip.fade.bs-tooltip-right.show').hide()
        }
    }

    $('li[data-nav-target]').click(function (e) {
        e.preventDefault()
        selectPageSection($(this).attr('data-nav-target'), true)
    })


    // Configure tooltips for collapsed side navigation
    $('.navbar-sidenav [data-toggle="tooltip"]').tooltip({
        template: '<div class="tooltip navbar-sidenav-tooltip" role="tooltip"><div class="arrow"></div><div class="tooltip-inner"></div></div>'
    })
    // Toggle the side navigation
    $("#sidenavToggler").click(function (e) {
        e.preventDefault();
        $("body").toggleClass("sidenav-toggled");
        $(".navbar-sidenav .nav-link-collapse").addClass("collapsed");
        $(".navbar-sidenav .sidenav-second-level, .navbar-sidenav .sidenav-third-level").removeClass("show");
    });
    // Force the toggled class to be removed when a collapsible nav link is clicked
    $(".navbar-sidenav .nav-link-collapse").click(function (e) {
        e.preventDefault();
        $("body").removeClass("sidenav-toggled");
    });
    // Prevent the content wrapper from scrolling when the fixed side navigation hovered over
    $('body.fixed-nav .navbar-sidenav, body.fixed-nav .sidenav-toggler, body.fixed-nav .navbar-collapse').on('mousewheel DOMMouseScroll', function (e) {
        var e0 = e.originalEvent,
            delta = e0.wheelDelta || -e0.detail;
        this.scrollTop += (delta < 0 ? 1 : -1) * 30;
        e.preventDefault();
    });
    // Scroll to top button appear
    $(document).scroll(function () {
        var scrollDistance = $(this).scrollTop();
        if (scrollDistance > 100) {
            $('.scroll-to-top').fadeIn();
        } else {
            $('.scroll-to-top').fadeOut();
        }
    });
    // Configure tooltips globally
    $('[data-toggle="tooltip"]').tooltip()
    // Smooth scrolling using jQuery easing
    $(document).on('click', 'a.scroll-to-top', function (event) {
        var $anchor = $(this);
        $('html, body').stop().animate({
            scrollTop: ($($anchor.attr('href')).offset().top)
        }, 1000, 'easeInOutExpo');
        event.preventDefault();
    });
})(jQuery); // End of use strict
