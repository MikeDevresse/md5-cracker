import "./ResultHistory.scss"
import {Component} from "react";
import {Card, Col, Table} from "react-bootstrap";
import CardHeader from "react-bootstrap/CardHeader";

class ResultHistory extends Component {
    keyCounter = 0;

    render() {
        const history = this.props.resultHistory.map(res => {
            return <tr key={this.keyCounter++}><td>{res.hash}</td><td>{res.result}</td></tr>
        });

        return (
            <Col xs={12} className="resultHistory">
                <Card>
                    <CardHeader className="card-header">Result History</CardHeader>
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