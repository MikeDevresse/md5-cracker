import "./CommandPanel.scss"
import {Component} from "react";
import ResultHistory from "./ResultHistory";

class CommandPanel extends Component {
    render() {
        return(
            <div className="commandPanel">
                <h2>MD5 Cracker</h2>
                <div className="panels">
                    <ResultHistory resultHistory={this.props.state.resultHistory}/>
                </div>
            </div>
        )
    }
}

export default CommandPanel