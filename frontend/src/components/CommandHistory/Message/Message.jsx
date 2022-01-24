import "./Message.scss"
import {Component} from "react";

class Message extends Component {
    constructor(props) {
        super(props);
        console.log(this.props)
        this.state = this.props.message
    }

    render() {
        return <div className="message">
            <span className="date">{this.state.date}</span>
            <span className="messageText">{this.state.msg}</span>
        </div>
    }
}

export default Message;