import "./CommandPanel.scss"
import {Component} from "react";
import ResultHistory from "./ResultHistory";
import CheckRequest from "./CheckRequest";
import {Button, Row} from "react-bootstrap";
import Monitor from "./Monitor";
import Configuration from "./Configuration";
import Websocket from "./Websocket";

class CommandPanel extends Component {
    render() {
        return(
            <div className="commandPanel">
                <div className="d-flex justify-content-between align-items-center">
                    <h2>MD5 Cracker</h2>
                    <div>
                        <Button onClick={this.props.toggleConsole}>{this.props.state.consoleDeveloped ? "Reduce" : "Develop"}</Button>
                    </div>
                </div>

                <Row>
                    <Websocket socketConnector={this.props.socketConnector}/>
                    <CheckRequest socket={this.props.state.socket}/>
                    <Monitor state={this.props.state}/>
                    <Configuration state={this.props.state} />
                    <ResultHistory resultHistory={this.props.state.resultHistory}/>
                </Row>
            </div>
        )
    }
}

export default CommandPanel