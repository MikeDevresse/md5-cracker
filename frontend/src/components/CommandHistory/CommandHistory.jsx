import "./CommandHistory.scss"
import {Component} from "react";
import Message from "./Message"

class CommandHistory extends Component {
    render() {
        const messages = this.props.commandHistory.map(msg => <Message message={msg.data} />);

        return (
            <div className="commandHistory">
                <h2>Command History</h2>
                {messages}
            </div>
        );
    }
}

export default CommandHistory