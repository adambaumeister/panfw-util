class LogBox {
    /*
    A white box that will expand, contract, and fill with data.
    Useful for transitory information.
     */
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
    /*
    Simple, one row "table" where each option can be clicked to expand to additional values.
     */
    constructor(id) {
        this.id = id;
        this.limit = 10;
    }

    DrawFromList(list) {
        var html = "";
        var mainobj = $("#main");
        var totalips = 0;
        $.each(list, function (index, element) {
            var ipcount = element.length;
            totalips = totalips + ipcount;
        });

        $.each(list, function (index, element) {
            var ipcount = element.length;
            var divwidth = ipcount/totalips *100;

            var ipHtml = element.join("<br>")
            html = html + `<div id="${index}-tagname" class='tag' style="width:${divwidth}%">${index}</div>`
        });
        $("#"+this.id).html(html);
        return html;
    }

}

module.exports = {
    InputList,
    LogBox,
    ExpandTable
};