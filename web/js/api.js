
class Panutil {
    constructor(url) {
        this.url = url;
        this.options = {
            json:true
        }
        this.key = ""

    }

    GetStatus() {
        var options = this.options;
        var url = this.url + '/status';
        return new Promise(function(resolve, reject) {
            options['url'] = url;
            request(options, function (error, response, body) {
                if (error) {
                    console.log(error);
                    reject({
                        "Status": 255,
                        "Message": `Error in request! ${error}`
                    })
                }
                resolve(body)
            });
        });
    }

    Login(loginJson) {
        var url = this.url + '/login';
        var options = {
            uri: url,
            method: 'POST',
            json: loginJson
        };
        return new Promise(function(resolve, reject) {
            request(options,
                function (error, response, body) {
                if (error) {
                    console.log(error);
                    reject({
                        "Status": 255,
                        "Message": `Error in request! ${error}`
                    })
                }
                this.key = body["ApiKey"];
                resolve(body)
            });
        });
    }

    Join(formJson) {
        var url = this.url + '/join';
        var options = {
            uri: url,
            method: 'POST',
            json: formJson
        };
        return new Promise(function(resolve, reject) {
            request(options,
                function (error, response, body) {
                    if (error) {
                        console.log(error);
                        reject({
                            "Status": 255,
                            "Message": `Error in request! ${error}`
                        })
                    }
                    resolve(body)
                });
        });
    }

    Register(formJson) {
        var url = this.url + '/register';
        var options = {
            uri: url,
            method: 'POST',
            json: formJson
        };
        return new Promise(function(resolve, reject) {
            request(options,
                function (error, response, body) {
                    if (error) {
                        console.log(error);
                        reject({
                            "Status": 255,
                            "Message": `Error in request! ${error}`
                        })
                    }
                    resolve(body)
                });
        });
    }
}

module.exports = Panutil;