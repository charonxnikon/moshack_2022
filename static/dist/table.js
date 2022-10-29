var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _possibleConstructorReturn(self, call) { if (!self) { throw new ReferenceError("this hasn't been initialised - super() hasn't been called"); } return call && (typeof call === "object" || typeof call === "function") ? call : self; }

function _inherits(subClass, superClass) { if (typeof superClass !== "function" && superClass !== null) { throw new TypeError("Super expression must either be null or a function, not " + typeof superClass); } subClass.prototype = Object.create(superClass && superClass.prototype, { constructor: { value: subClass, enumerable: false, writable: true, configurable: true } }); if (superClass) Object.setPrototypeOf ? Object.setPrototypeOf(subClass, superClass) : subClass.__proto__ = superClass; }

var Line = function (_React$Component) {
    _inherits(Line, _React$Component);

    function Line() {
        _classCallCheck(this, Line);

        return _possibleConstructorReturn(this, (Line.__proto__ || Object.getPrototypeOf(Line)).apply(this, arguments));
    }

    _createClass(Line, [{
        key: "render",
        value: function render() {
            return React.createElement(
                "div",
                { "class": "table__line" },
                React.createElement(
                    "div",
                    { "class": "table__line-fields" },
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        " ",
                        this.props.data.address,
                        " "
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        " ",
                        this.props.data.rooms,
                        " "
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        " ",
                        this.props.data.type,
                        " "
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        " ",
                        this.props.data.height,
                        " "
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        " ",
                        this.props.data.material,
                        " "
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        " ",
                        this.props.data.floor,
                        " "
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        " ",
                        this.props.data.area,
                        " "
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        " ",
                        this.props.data.kitchen,
                        " "
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        " ",
                        this.props.data.balcony,
                        " "
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        " ",
                        this.props.data.metro,
                        " "
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        " ",
                        this.props.data.condition,
                        " "
                    )
                ),
                React.createElement(
                    "div",
                    { "class": "table__line-navigation" },
                    React.createElement("button", { "class": "table__line-navbutton table__line-navbutton-up", type: "button" }),
                    React.createElement("button", { "class": "table__line-navbutton table__line-navbutton-down", type: "button" }),
                    "// ",
                    React.createElement("button", { "class": "table__line-navbutton table__line-navbutton-down", type: "button", onClick: this.moveDown })
                )
            );
        }
    }]);

    return Line;
}(React.Component);

var Table = function (_React$Component2) {
    _inherits(Table, _React$Component2);

    function Table(props) {
        _classCallCheck(this, Table);

        var _this2 = _possibleConstructorReturn(this, (Table.__proto__ || Object.getPrototypeOf(Table)).call(this, props));

        _this2.state = {
            data: []
        };

        _this2.state = { data: init_data };
        return _this2;
    }

    _createClass(Table, [{
        key: "render",
        value: function render() {
            return React.createElement(
                "div",
                { "class": "table__container" },
                React.createElement(
                    "div",
                    { "class": "table__header" },
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        "Address"
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        "Rooms"
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        "Type"
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        "Height"
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        "Material"
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        "Floor"
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        "Area"
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        "Kitchen"
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        "Balcony"
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        "Subway"
                    ),
                    React.createElement(
                        "div",
                        { "class": "table__field" },
                        "Condition"
                    )
                ),
                React.createElement(
                    "div",
                    { "class": "table__table" },
                    this.state.data.map(function (item, index) {
                        return React.createElement(Line, { data: item, index: index });
                    })
                )
            );
        }
    }]);

    return Table;
}(React.Component);

var init_data = [{
    address: "улица Пупкина д. 1",
    rooms: "3",
    type: "Новостройка",
    height: "10",
    material: "кирпич",
    floor: "7",
    area: "100",
    kitchen: "20",
    balcony: "есть",
    metro: "Авиамоторная",
    condition: "под чистовую отделку"
}, {
    address: "улица Пупкина д. 1",
    rooms: "3",
    type: "Новостройка",
    height: "10",
    material: "кирпич",
    floor: "7",
    area: "100",
    kitchen: "20",
    balcony: "есть",
    metro: "Авиамоторная",
    condition: "под чистовую отделку"
}, {
    address: "улица Пупкина д. 1",
    rooms: "3",
    type: "Новостройка",
    height: "10",
    material: "кирпич",
    floor: "7",
    area: "100",
    kitchen: "20",
    balcony: "есть",
    metro: "Авиамоторная",
    condition: "под чистовую отделку"
}, {
    address: "улица Пупкина д. 1",
    rooms: "3",
    type: "Новостройка",
    height: "10",
    material: "кирпич",
    floor: "7",
    area: "100",
    kitchen: "20",
    balcony: "есть",
    metro: "Авиамоторная",
    condition: "под чистовую отделку"
}];

var el = document.querySelector('#table');
ReactDOM.render(React.createElement(Table, null), el);