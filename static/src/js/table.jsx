class Line extends React.Component {
    render() {
        return (
            <div class="table__line">
                <div class="table__line-fields">
                    <div class="table__field"> {this.props.data.address} </div>
                    <div class="table__field"> {this.props.data.rooms} </div>
                    <div class="table__field"> {this.props.data.type} </div>
                    <div class="table__field"> {this.props.data.height} </div>
                    <div class="table__field"> {this.props.data.material} </div>
                    <div class="table__field"> {this.props.data.floor} </div>
                    <div class="table__field"> {this.props.data.area} </div>
                    <div class="table__field"> {this.props.data.kitchen} </div>
                    <div class="table__field"> {this.props.data.balcony} </div>
                    <div class="table__field"> {this.props.data.metro} </div>
                    <div class="table__field"> {this.props.data.condition} </div>
                </div>

                <div class="table__line-navigation">
                    <button class="table__line-navbutton table__line-navbutton-up" type="button"></button>
                    <button class="table__line-navbutton table__line-navbutton-down" type="button"></button>
                    // <button class="table__line-navbutton table__line-navbutton-down" type="button" onClick={this.moveDown}></button>
                </div>
            </div>
        );
    }
}

class Table extends React.Component {

    constructor(props) {
        super(props);
        this.state = { data: init_data };
    }

    state = {
        data: [],
    };

    render() {
        return (
            <div class="table__container">
                <div class="table__header">
                    <div class="table__field">Address</div>
                    <div class="table__field">Rooms</div>
                    <div class="table__field">Type</div>
                    <div class="table__field">Height</div>
                    <div class="table__field">Material</div>
                    <div class="table__field">Floor</div>
                    <div class="table__field">Area</div>
                    <div class="table__field">Kitchen</div>
                    <div class="table__field">Balcony</div>
                    <div class="table__field">Subway</div>
                    <div class="table__field">Condition</div>
                </div>
                <div class="table__table">
                    {this.state.data.map((item, index) => (
                        <Line data={item} index={index} />
                    ))}
                </div>
            </div>
        );
    }
}

let init_data = [
    {
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
        condition: "под чистовую отделку",
    },
    {
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
        condition: "под чистовую отделку",
    },{
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
        condition: "под чистовую отделку",
    },{
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
        condition: "под чистовую отделку",
    },
];

let el = document.querySelector('#table');
ReactDOM.render(<Table />, el);
