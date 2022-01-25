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
            slavesAvailable: 0,
            slavesWorking: 0,
            maxSearch: "",
            maxSlavesPerRequest: 0,
            autoScale: false,
            queue: 0,
            searching: 0,
            consoleDeveloped: true
        }
    }

    toggleConsole = () => {
        this.setState({consoleDeveloped: !this.state.consoleDeveloped})
    }

    componentDidMount() {
        connect((msg) => {
            this.setState(_ => {
                let msgSplit = msg.data.split(" ")
                let state = {
                    commandHistory: [...this.state.commandHistory, {msg: msg.data, date: (new Date()).toLocaleString()}],
                    resultHistory: this.state.resultHistory,
                    slaves: this.state.slaves,
                    slavesAvailable: this.state.slavesAvailable,
                    slavesWorking: this.state.slavesWorking,
                    maxSearch: this.state.maxSearch,
                    maxSlavesPerRequest: this.state.maxSlavesPerRequest,
                    autoScale: this.state.autoScale,
                    queue: this.state.queue,
                    searching: this.state.searching,
                    consoleDeveloped: this.state.consoleDeveloped
                }
                if(msgSplit[0] === "found" && msgSplit.length === 6) {
                    state.resultHistory = [{
                        hash: msgSplit[1],
                        result: msgSplit[2],
                        requestedAt: new Date(parseInt(msgSplit[3])),
                        startedAt: new Date(parseInt(msgSplit[4])),
                        endedAt: new Date(parseInt(msgSplit[5]))
                    }, ...this.state.resultHistory]
                }
                else if(msgSplit[0] === "slaves" && msgSplit.length === 4) {
                    state.slaves = msgSplit[1]
                    state.slavesAvailable = msgSplit[2]
                    state.slavesWorking = msgSplit[3]
                }
                else if(msgSplit[0] === "queue" && msgSplit.length === 3) {
                    state.queue = msgSplit[1]
                    state.searching = msgSplit[2]
                }
                else if(msgSplit[0] === "max-search" && msgSplit.length === 2) {
                    state.maxSearch = msgSplit[1]
                }
                else if(msgSplit[0] === "max-slaves-per-request" && msgSplit.length === 2) {
                    state.maxSlavesPerRequest = msgSplit[1]
                }
                else if(msgSplit[0] === "auto-scale" && msgSplit.length === 2) {
                    state.autoScale = msgSplit[1] === "true"
                }
                return state
            })
        });
    }

    render() {
        return (
            <div className="App">
                <div className="content">
                    <CommandPanel state={this.state} toggleConsole={this.toggleConsole} />
                    <CommandHistory isDeveloped={this.state.consoleDeveloped} commandHistory={this.state.commandHistory} />
                </div>
            </div>
        );
    }
}

export default App;
