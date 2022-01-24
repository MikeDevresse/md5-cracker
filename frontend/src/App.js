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
            resultHistory: []
        }
    }

    componentDidMount() {
        connect((msg) => {
            this.setState(_ => {
                let msgSplit = msg.data.split(" ")
                let state = {
                    commandHistory: [...this.state.commandHistory, {msg: msg.data, date: (new Date()).toLocaleString()}],
                    resultHistory: this.state.resultHistory
                }
                console.log(msgSplit[0], msgSplit.length)
                if(msgSplit[0] === "found" && msgSplit.length === 3) {
                    state.resultHistory = [...this.state.resultHistory, {hash: msgSplit[1], result: msgSplit[2]}]
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
