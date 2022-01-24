import './App.css';
import {Component} from "react";
import {connect, sendMsg} from "./api";
import Header from "./components/Header";
import CommandHistory from "./components/CommandHistory";
import CommandInput from "./components/CommandInput";

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            commandHistory: []
        }
    }

    componentDidMount() {
        connect((msg) => {
            console.log("New Message")
            this.setState(_ => ({
                commandHistory: [...this.state.commandHistory, msg]
            }))
            console.log(this.state);
        });
    }

    send(event) {
        if(event.keyCode === 13) {
            sendMsg(event.target.value);
            event.target.value = "";
        }
    }

    render() {
        return (
            <div className="App">
                <Header/>
                <CommandHistory commandHistory={this.state.commandHistory} />
                <CommandInput send={this.send} />
            </div>
        );
    }
}

export default App;
