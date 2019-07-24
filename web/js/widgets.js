class LogBox {
    constructor(id) {
        this.id = id;
        this.logs = $(this.id);
    }

    OpenWithFill(contents) {
        console.log("We got called to open")
        var logs = this.logs;
        logs.css('height', '30px');
        logs.css('padding-top', '20px');
        logs.css('padding-bottom', '20px');
        logs.on('transitionend webkitTransitionEnd oTransitionEnd otransitionend MSTransitionEnd',
            function() {
                $(this).html(contents)
            });

    }
    Close() {
        var logs = this.logs;
        logs.html("");
        logs.css('height', '0px');
        logs.css('padding-top', '0px');
        logs.css('padding-bottom', '0px');
        // This clears the previously registered thing
        logs.on('transitionend webkitTransitionEnd oTransitionEnd otransitionend MSTransitionEnd',
            function() {
                $(this).html("");
            });
    }
}


class InputList {
    /*
    Widget for bulding dynamic input lists
    List output is returned as array but not nessecarily ordered by which it appears on the page.
     */
    constructor(className) {
        this.className = className;
        this.items = []
    }

    Add(value) {
        var html;
        if (value != null) {
            html = `<div><input class="${this.className}" type="text" value="${value}"></div>`;
        } else {
            html = `<div><input class="${this.className}" type="text"></div>`;
        }

        this.items.push(html)
    }

    Render() {
        var html = this.items.join("")
        html = html + "<button class='add-item'>+</bbutton>"
        return html;
    }

    GetValues() {
        var values = []
        $('.'+this.className).each(function() {
            values.push(this.value)
        });
        return values;
    }
}

class ExpandTable {
    constructor(id) {
        this.id = id;
    }

    DrawFromList(list) {
        var html = "";
        $.each(list, function (index, element) {
            console.log(element)
            html = html + `<div class='tag'>${index}</div>`
        });
        return html;
    }

}

module.exports = {
    InputList,
    LogBox,
    ExpandTable
};