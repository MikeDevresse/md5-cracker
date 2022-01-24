import "./CheckRequest.scss"
import {Component} from "react";
import {sendMsg} from "../../../api";
import {Button, Card, Col, Form} from "react-bootstrap";
import CardHeader from "react-bootstrap/CardHeader";

class CheckRequest extends Component {
    constructor(props) {
        super(props);

        this.state = {
            hash: ''
        }

        this.send = this.send.bind(this);
        this.handleChange = this.handleChange.bind(this);
    }

    handleChange(event) {
        this.setState({hash: event.target.value});
    }

    send(event) {
        event.preventDefault()
        sendMsg("search "+this.state.hash)
        this.setState({hash: ''})
    }

    render() {
        return (
            <Col xs={12} lg={6} className="checkRequest">
                <Card>
                    <Card.Header>Send a crack request</Card.Header>
                    <Card.Body>
                        <Form onSubmit={this.send}>
                            <Form.Group className="mb-2">
                                <Form.Control name="checkRequest" onChange={this.handleChange} value={this.state.hash} placeholder="Hash"/>
                            </Form.Group>
                            <Button type="submit" className="w-100" variant="success">Crack !</Button>
                        </Form>
                    </Card.Body>
                </Card>
            </Col>
        )
    }
}

export default CheckRequest