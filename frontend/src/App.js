import './App.scss';
import {Component} from "react";
import {connect} from "./api";
import CommandHistory from "./components/CommandHistory";
import CommandPanel from "./components/CommandPanel";

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            commandHistory: []
        }
    }

    componentDidMount() {
        connect((msg) => {
            this.setState(_ => ({
                commandHistory: [...this.state.commandHistory, {msg: msg.data, date: (new Date()).toLocaleString()}]
            }))
        });
    }

    render() {
        return (
            <div className="App">
                <div className="content">
                    <CommandPanel />
                    <CommandHistory commandHistory={this.state.commandHistory} />
                </div>
            </div>
        );
    }
}

export default App;
