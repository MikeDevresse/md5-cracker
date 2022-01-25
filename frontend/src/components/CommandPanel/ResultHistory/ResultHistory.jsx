import "./ResultHistory.scss"
import {Component} from "react";
import {Card, Col, Table} from "react-bootstrap";

class ResultHistory extends Component {
    keyCounter = 0;


    pad(num, size) {
        num = num.toString();
        while (num.length < size) num = "0" + num;
        return num;
    }

    calculateTimeDiff(date1, date2) {
        let diff;
        if(date1 > date2) {
            diff = date1 - date2
        }
        else {
            diff = date2 - date1
        }
        let msec = diff;
        let mm = Math.floor(msec / 1000 / 60);
        msec -= mm * 1000 * 60;
        let ss = Math.floor(msec / 1000);
        msec -= ss * 1000;

        return this.pad(mm,2) + ":" + this.pad(ss,2) + "." + this.pad(msec,3)
    }

    render() {
        const history = this.props.resultHistory.map(res => {
            return <tr key={this.keyCounter++}>
                <td>{res.hash}</td>
                <td>{this.calculateTimeDiff(res.requestedAt,res.startedAt)}</td>
                <td>{this.calculateTimeDiff(res.startedAt,res.endedAt)}</td>
                <td>{res.result}</td>
            </tr>
        });

        return (
            <Col xs={12} className="resultHistory">
                <Card>
                    <Card.Header className="card-header">Result History</Card.Header>
                    <Table className="mb-0">
                        <thead>
                        <tr>
                            <th>Hash</th>
                            <th>Time to init</th>
                            <th>Time taken</th>
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