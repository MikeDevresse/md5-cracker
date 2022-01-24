import "./CommandPanel.scss"
import {Component} from "react";
import ResultHistory from "./ResultHistory";
import CheckRequest from "./CheckRequest";
import {Container, Row} from "react-bootstrap";

class CommandPanel extends Component {
    render() {
        return(
            <Container className="commandPanel">
                <h2>MD5 Cracker</h2>
                <Row>
                    <CheckRequest/>
                    <ResultHistory resultHistory={this.props.state.resultHistory}/>
                </Row>
            </Container>
        )
    }
}

export default CommandPanel