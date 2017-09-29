(function ($) {
    // see https://stackoverflow.com/a/18234317
    String.prototype.formatUnicorn = function () {

        var str = this.toString();
        if (arguments.length) {
            var t = typeof arguments[0];
            var key;
            var args = ("string" === t || "number" === t) ?
                Array.prototype.slice.call(arguments)
                : arguments[0];

            for (key in args) {
                str = str.replace(new RegExp("\\{" + key + "\\}", "gi"), args[key]);
            }
        }

        return str;
    };

    "use strict" // Start of use strict


    function setup() {
        selectPageSection("devices")


    }

    setup()

    $.getJSON("/api/profile", function (data) {
        console.log("Json resp:", data)
        window.Profile = data
        console.log(JSON.stringify(data, " ", " "))

        onProfileLoaded()
    })

    function onLightButtonClick() {
        console.log("ok")
    }

    function executeFunction(cmd) {
        console.log("posting:", cmd)
        $.post("/api/exec", JSON.stringify(cmd), function (resp) {
            console.log("api cmd response:", resp)
        }, 'json')
    }


    function onProfileLoaded() {
        setupDeviceCards()
    }

    function setupDeviceCards() {
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
                console.log("hi")
                newFunctionLine(device)
                    .hide()
                    .appendTo(deviceElement.find('ul'))
                    .slideDown();
            })
        })
        if (Profile.functions !== null) {
            Profile.functions.forEach(function (f) {
                // add each existing function to the list
                var id = 'device-' + f.deviceid
                var deviceElement = $('#' + id)
                try {
                    getWidget(f).appendTo(deviceElement.find('ul'))
                } catch (e) {
                    console.log("No such device:", id, e)
                }
            })
        }
    }

    function newFunctionLine(device) {
        var template = $('.template-functionline').html()
        var elm = $(template)
        var functionline = $(elm.closest('li.list-group-item'))
        console.log("fl:", functionline)

        function dismiss() {
            functionline.slideUp(functionline.remove)
        }

        elm.find('.btn-cancel').click(function () {
            dismiss()
        })

        elm.find('.btn-save').click(function () {
            console.log("saving...")
            var func = {
                "name": functionline.find('.function-name').val(),
                "pin": parseInt(functionline.find('.function-pin').val()),
                "cmd": "DW",
                "deviceid": device.id,
                "data": {
                    "uielement": functionline.find('.function-dw-type').val(),
                    "invert": functionline.find('.function-dw-invert').is(':checked'),
                },
            }

            console.log("func:", func)
            $.post("/api/newfunction", JSON.stringify(func), function (resp) {
                console.log("api newfunction response:", resp)
            }, 'json')

            dismiss()
        })

        return elm
    }

    function getWidget(func) {
        console.log("newfunc!")
        var tmpl = '<li class="list-group-item"><span class="text-muted align-middle">{name}</span>\n' +
            '  <div class="switch">\n' +
            '    <label>\n' +
            '      <input id="lightswitch" type="checkbox" checked="checked"/><span class="slider round"></span>\n' +
            '    </label>\n' +
            '  </div>\n' +
            '</li>'

        var elm = $(tmpl.formatUnicorn({
            name: func.name,
        }))
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
