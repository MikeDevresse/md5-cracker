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
                       <tr>
                           <td>Slaves</td>
                           <td>{this.props.state.slaves}</td>
                       </tr>
                       <tr>
                           <td>Queue</td>
                           <td>{this.props.state.queue}</td>
                       </tr>
                   </Table>
               </Card>
            </Col>
        )
    }
}

export default Monitor