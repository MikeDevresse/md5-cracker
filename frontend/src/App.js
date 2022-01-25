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
            slavesAvailable: 0,
            slavesWorking: 0,
            queue: 0,
            searching: 0
        }
    }

    componentDidMount() {
        connect((msg) => {
            this.setState(_ => {
                let msgSplit = msg.data.split(" ")
                let state = {
                    commandHistory: [...this.state.commandHistory, {msg: msg.data, date: (new Date()).toLocaleString()}],
                    resultHistory: this.state.resultHistory,
                    slavesAvailable: this.state.slavesAvailable,
                    slavesWorking: this.state.slavesWorking,
                    queue: this.state.queue,
                    searching: this.state.searching
                }
                if(msgSplit[0] === "found" && msgSplit.length === 3) {
                    state.resultHistory = [...this.state.resultHistory, {hash: msgSplit[1], result: msgSplit[2]}]
                }
                if(msgSplit[0] === "slaves" && msgSplit.length === 3) {
                    state.slavesAvailable = msgSplit[1]
                    state.slavesWorking = msgSplit[2]
                }
                if(msgSplit[0] === "queue" && msgSplit.length === 3) {
                    state.queue = msgSplit[1]
                    state.searching = msgSplit[2]
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
