import "./Message.scss"
import {Component} from "react";

class Message extends Component {
    constructor(props) {
        super(props);
        let temp = this.props.message;
        this.state = {
            message: temp
        }
    }

    render() {
        return <div className="message">{this.state.message}</div>
    }
}

export default Message;