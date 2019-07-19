const { exec } = require('child_process');
const { spawn } = require('child_process');
const request = require('request');
const Cookies = require('js-cookie');
const panexe = '../panutil.exe';
const url = "http://localhost:8080";
const storage = require('electron-json-storage');
const Store = require('electron-store');
const Panutil = require('./js/api');
const InputList = require('./js/widgets');
const storebasic = new Store()

class Main {
    constructor(elementid) {
        this.elementid = elementid;

    }

    SetLoginStatus(status) {
        // Update the status of the device connection.
        var statusVal = status['Status'];
        if (statusVal === 0) {
            $("#dstatus").replaceWith(`<span id="dstatus" class="good-dot"></span>`);
            this.DisplayIndex();
            store("keyinfo", status)
        } else if (statusVal === 1){
            $("#dstatus").replaceWith(`<span id="dstatus" class="bad-dot"></span>`);
            var logs = $("#logs")
            logs.css('height', '30px');
            logs.css('padding-top', '20px');
            logs.css('padding-bottom', '20px');
            logs.on('transitionend webkitTransitionEnd oTransitionEnd otransitionend MSTransitionEnd',
                function() {
                    $(this).html("Connection failed: bad password.")
                });

        }
    }

    DisplayIndex() {
        $("#main").load('htmlsnippets/links.html')
    }

    DisplayRegister(panutil) {
        /*
        This function displays the "register" submenu
         */
        $("#main").load('htmlsnippets/register.html', function() {
            var inputList = new InputList("list-widget");
            inputList.Add(null);
            var html = inputList.Render()
            $("#main").on('click', '#submit-register', function () {
                // On click we retrieve all the values from the input list widget
                var values = inputList.GetValues();
                // We then push the tag as the final argument - this will be sent to the panutil API
                values.push($("#tag").val());
                var apikey = get('keyinfo');
                var jsonform = {
                    'ApiKey': apikey['ApiKey'],
                    'Command': 'register',
                    'args': values,
                }
                console.log(jsonform)
                panutil.Register(jsonform).then(function (val) {
                    $("#main").html(val['Message'])
                }, function (err) {
                    console.log(val)
                });
            });
            $("#register-list").html(html);
        })
    }
}

// In the below, the string is the command
// The (*)=>{*} syntax represents a function, a callback function.
function PrintHelp() {
    exec(panexe + ' --help', (err, stdout, stderr) => {
        if (err) {
            alert("Error running command!");
            return
        }
        console.log(`stdout: ${stdout}`);
        console.log(`stderr: ${stderr}`);
    });
}

function StartPanutil() {
    /*
    Attempt to start the panutil runtime.
    If it exits prematurely, we will restart after a short time to ensure it's fully closed.
    Stderr and stdout go the chromium console log.
     */
    let child = spawn(panexe, ['api']);
    child.on('exit', code => {
        console.log("Panutil API exited, restarting in 2...")
        setTimeout(function () {
            StartPanutil();
        }, 2000)
    });

    child.stdout.on('data', (data) => {
        console.log(`child stdout:\n${data}`);
    });

    child.stderr.on('data', (data) => {
        console.log(`child stderr:\n${data}`);
    });

    var panutil = new Panutil(url);
    return panutil
}

function SetStatus(status) {
    //Update the panutil status.
    // { 'Message': "This is a status message" }
    if (status['Status'] === 0) {
        $("#pstatus").replaceWith(`<span class="good-dot"></span>`)
    } else {
        console.log(status)
        $("#pstatus").replaceWith(`<span class="bad-dot"></span>`)
    }

}

function ConvertFormToJSON(form){
    // Given a form object (select it first using jq), converts it to Json and returns.
    var array = jQuery(form).serializeArray();
    var json = {};

    jQuery.each(array, function() {
        json[this.name] = this.value || '';
    });
    return json;
}

function store(key, jsonValue) {
    storebasic.set(key, jsonValue)
}

function get(key) {
    return storebasic.get(key)
}

function Join(error,data) {
    var panutil = new Panutil(url);
    panutil.Join(data).then(function (val) {
        console.log(val)
    }, function (err) {
        console.log(val)
    });
}

function poploginfields(data) {
    $("#username").val(data['Username'])
    $("#hostname").val(data['Hostname'])
}

$(document).ready(function(){
    // Register the page classes
    var main = new Main("#main");

    poploginfields(get("loginForm"));

    var panutil = StartPanutil()
    // Set the status of Panutil
    panutil.GetStatus().then(function (val) {
        SetStatus(val)
    }, function (err) {
        SetStatus(err)
    });

    // Convert the form login-form into JSON and attempt to login.
    // This will login to the firewall.
    $("#loginb").click(function () {
        var loginJson = ConvertFormToJSON($("#login-form"));

        panutil.Login(loginJson).then(function (val) {
            main.SetLoginStatus(val)
        }, function (err) {
            main.SetLoginStatus(val)
        });
        delete loginJson['Password']
        store("loginForm", loginJson)

    });

    $("#main").on('click', '#register-link', function () {
        main.DisplayRegister(panutil)
    });
});
