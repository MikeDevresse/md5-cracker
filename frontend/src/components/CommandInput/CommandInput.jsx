import "./CommandInput.scss"
import {Component} from "react";

class CommandInput extends Component {
    render() {
        return (
            <div className="commandInput">
                <input onKeyDown={this.props.send} />
            </div>
        )
    }
}

export default CommandInput