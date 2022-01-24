import "./ResultHistory.scss"
import {Component} from "react";
import {Card, Col, Table} from "react-bootstrap";

class ResultHistory extends Component {
    keyCounter = 0;

    render() {
        const history = this.props.resultHistory.map(res => {
            return <tr key={this.keyCounter++}><td>{res.hash}</td><td>{res.result}</td></tr>
        });

        return (
            <Col xs={12} className="resultHistory">
                <Card>
                    <Card.Header className="card-header">Result History</Card.Header>
                    <Table className="mb-0">
                        <thead>
                        <tr>
                            <th>Hash</th>
                            <th>Value</th>
                        </tr>
                        </thead>
                        <tbody>
                            {history}
                        </tbody>
                    </Table>
                </Card>
            </Col>
        )
    }
}

export default ResultHistory