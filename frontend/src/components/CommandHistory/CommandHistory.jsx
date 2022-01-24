import "./CommandHistory.scss"
import {Component} from "react";
import Message from "./Message"
import CommandInput from "./CommandInput";
import {sendMsg} from "../../api";

class CommandHistory extends Component {
    keyCounter = 0;

    send(event) {
        if(event.keyCode === 13) {
            sendMsg(event.target.value);
            event.target.value = "";
        }
    }

    render() {
        const messages = this.props.commandHistory.map(msg => <Message key={this.keyCounter++} message={msg} />);

        return (
            <div className="commandHistory">
                <h2>Command History</h2>
                <div className="commandHistoryList">
                    {messages}
                </div>
                <CommandInput send={this.send} />
            </div>
        );
    }
}

export default CommandHistory