
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

module.exports = InputList;