import "./Monitor.scss"
import {Component} from "react";
import {Card, Col, Table} from "react-bootstrap";

class Monitor extends Component {
    render() {
        return (
            <Col xs={12} lg={6}  className="monitor">
               <Card>
                   <Card.Header>Monitoring</Card.Header>
                   <Table className="mb-0">
                       <tbody>
                       <tr>
                           <td>Slaves available</td>
                           <td>{this.props.state.slavesAvailable}</td>
                       </tr>
                       <tr>
                           <td>Slaves working</td>
                           <td>{this.props.state.slavesWorking}</td>
                       </tr>
                       <tr>
                           <td>Queue</td>
                           <td>{this.props.state.queue}</td>
                       </tr>
                       <tr>
                           <td>Searching</td>
                           <td>{this.props.state.searching}</td>
                       </tr>
                       </tbody>
                   </Table>
               </Card>
            </Col>
        )
    }
}

export default Monitor