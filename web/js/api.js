
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

    GetRegistered(formJson) {
        var url = this.url + '/showregistered';
        console.log(url)
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
                    var ri = new RegisteredIPs(body);
                    resolve(ri)
                });
        });
    }
}

class RegisteredIPs {
    constructor(json) {
        this.json = json;
        this.byTag = {};
        this.byAddr = {};
        this.parseJson(this.json)
    }

    parseJson(json) {
        var self = this;
        $.each(json.Entries, function(index, element) {
            $.each(element.Tags, function (_, tag) {
                if (tag in self.byTag) {
                    self.byTag[tag].push(element.Ip);
                } else {
                    self.byTag[tag] = [element.Ip]
                }
            })
        });
    }

    GetTags() {
        return this.byTag
    }
}

module.exports = Panutil;