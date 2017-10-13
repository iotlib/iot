(function ($) {


    "use strict" // Start of use strict

    var Profile = {};


    function setup() {
        selectPageSection("devices")

        // load profile
        $.getJSON("/api/profile", function (data) {
            console.log("Json resp:", data)
            Profile = data
            console.log(JSON.stringify(data, " ", " "))

            onProfileLoaded()
        })


    }

    setup()


    function api(method, path, data, success, error) {

        $.ajax({
            url: '/api' + path,
            method: method,
            success: success,
            error: error,
            data: data,
            dataType: "json",
        })
    }

    window.api = api


    function executeFunction(cmd) {
        $.post("/api/exec", JSON.stringify(cmd), function (resp) {
            console.log("api cmd response:", resp)
        }, 'json')
    }


    function onProfileLoaded() {
        setupDeviceSection()
        setupAccountSection()
    }

    function setupAccountSection() {
        var tmpl = $('.template-collapsecard').html()
        $(tmpl.formatUnicorn({
            id: "account-card",
            title: Profile.user.name,
            subtitle: Profile.user.email,
        })).appendTo('.page-section-account')
    }

    function setupDeviceSection() {
        var section = $('.page-section-devices')
        section.empty()

        var template = $('.template-devicecard').html()
        Profile.devices.forEach(function (device) {
            console.log("device:", device)
            var online = device.lastseen - new Date().getTime() < 60
            var deviceElement = $(template.formatUnicorn({
                title: device.name,
                deviceid: device.id.substr(0, 7),
                online: online ? "online" : offline,
                id: 'device-' + device.id,
            }))
            section.append(deviceElement)


            deviceElement.find('.btn-add-function').on('click', function (e) {
                addNewFunctionLineToDevice(device.id)

            })
        })
        if (Profile.functions !== null) {
            Profile.functions.forEach(function (f) {
                // add each existing function to the list
                addFunctionToDevice(f)
            })
        }
    }

    function addFunctionToDevice(f) {
        var id = 'device-' + f.deviceid
        var deviceElement = $('#' + id)
        try {
            newSwitchFunction(f).appendTo(deviceElement.find('ul'))
        } catch (e) {
            console.log("No such device:", id, e)
        }
    }

    function addNewFunctionLineToDevice(deviceid, savecallback, f) {
        var template = $('.template-functionline').html()
        var elm = $(template)

        function dismiss() {
            elm.slideUp(elm.remove)
        }

        elm.find('.btn-cancel').click(function () {
            dismiss()
        })

        if (f !== undefined) {
            console.log("editing")
            elm.find('.function-name').val(f.name || '')
            elm.find('.function-pin').val(f.pin || '')
            elm.find('.function-dw-invert').prop('checked', f.data.invert || false)

            elm.find('.btn-cancel').click(function () {
                // add it back
                addFunctionToDevice(f)
            })
        }


        elm.find('.btn-save').click(function () {

            console.log("saving...")
            var func = {
                "name": elm.find('.function-name').val(),
                "pin": parseInt(elm.find('.function-pin').val()),
                "cmd": "DW",
                "deviceid": deviceid,
                "data": {
                    "uielement": elm.find('.function-dw-type').val(),
                    "invert": elm.find('.function-dw-invert').is(':checked'),
                },
            }
            console.log("func:", func)
            api("POST", "/function", JSON.stringify(func), function (data) {
                func.id = data.id
            }, function (e) {
            })

            if (savecallback !== undefined) savecallback()

            addFunctionToDevice(func)

            dismiss()
        })

        var deviceElement = $('#device-' + deviceid)

        elm.hide()
            .appendTo(deviceElement.find('ul'))
            .slideDown();


        return elm
    }

    function newSwitchFunction(func) {
        var tmpl = $('.template-switchfunction').html()
        var elm = $(tmpl.formatUnicorn({
            name: func.name,
        }))

        function dismiss() {
            elm.slideUp(elm.remove)
        }

        var value = false
        elm.find('input').change(function (e) {
            console.log("Clicked", this.checked)
            executeFunction({
                "id": func.deviceid,
                "cmd": "DW {pin} {val}".formatUnicorn({
                    pin: func.pin,
                    val: (this.checked ^ func.data.invert) ? "HIGH" : "LOW"
                })
            })
        })
        elm.find('.btn-remove').click(function (e) {
            dismiss()
            console.log(func)
            api("DELETE", "/function/" + func.id, null, function () {
                console.log("deleted function", func.id)
            })
        })
        elm.find('.btn-edit').click(function (e) {
            dismiss()
            console.log(func)
            addNewFunctionLineToDevice(func.deviceid, function () {
                api("DELETE", "/function/" + func.id, null, function () {
                    console.log("deleted function", func.id)
                })
            }, func)
        })
        return elm
    }


    function collapseNavbar() {
        $('#navbarResponsive').collapse('hide')
        $('.tooltip.navbar-sidenav-tooltip.fade.bs-tooltip-right.show').hide()
    }

    function selectPageSection(section) {
        $('.page-section').hide()
        var sec = $('.page-section-' + section).show()
        var tit = $('li[data-nav-target=' + section + ']').text()
        $('.navbar-brand').text(tit)

        collapseNavbar()
    }

    $('#nav-signout').click(function (e) {
        collapseNavbar()
    })
    $('li[data-nav-target]').click(function (e) {
        e.preventDefault()
        selectPageSection($(this).attr('data-nav-target'))
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
