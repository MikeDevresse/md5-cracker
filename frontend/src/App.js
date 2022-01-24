import './App.scss';
import {Component} from "react";
import {connect} from "./api";
import CommandHistory from "./components/CommandHistory";
import CommandPanel from "./components/CommandPanel";

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            commandHistory: [],
            resultHistory: [],
            slaves: 0,
            queue: 0
        }
    }

    componentDidMount() {
        connect((msg) => {
            this.setState(_ => {
                let msgSplit = msg.data.split(" ")
                let state = {
                    commandHistory: [...this.state.commandHistory, {msg: msg.data, date: (new Date()).toLocaleString()}],
                    resultHistory: this.state.resultHistory,
                    slaves: this.state.slaves,
                    queue: this.state.queue
                }
                if(msgSplit[0] === "found" && msgSplit.length === 3) {
                    state.resultHistory = [...this.state.resultHistory, {hash: msgSplit[1], result: msgSplit[2]}]
                }
                if(msgSplit[0] === "slaves" && msgSplit.length === 2) {
                    state.slaves = msgSplit[1]
                }
                if(msgSplit[0] === "queue" && msgSplit.length === 2) {
                    state.queue = msgSplit[1]
                }

                return state
            })
        });
    }

    render() {
        return (
            <div className="App">
                <div className="content">
                    <CommandPanel state={this.state} />
                    <CommandHistory commandHistory={this.state.commandHistory} />
                </div>
            </div>
        );
    }
}

export default App;
