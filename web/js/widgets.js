
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
            html = `<input class="${this.className}" type="text" value="${value}">`;
        } else {
            html = `<input class="${this.className}" type="text">`;
        }

        this.items.push(html)
    }

    Render() {
        var html = this.items.join("")
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