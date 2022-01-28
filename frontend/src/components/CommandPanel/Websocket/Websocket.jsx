import "./Websocket.scss"
import {Component} from "react";
import {Button, Col, FormControl, InputGroup} from "react-bootstrap";

class Websocket extends Component {
    constructor(props) {
        super(props);

        this.state = {
            link: 'ws://127.0.0.1:8080/ws',
            isLoading: false
        }

        this.connect = this.connect.bind(this);
        this.handleChange = this.handleChange.bind(this);
    }

    connect(event) {
        event.preventDefault()
        this.props.socketConnector(this.state.link)
    }

    handleChange(event) {
        this.setState({link: event.target.value});
    }

    render() {
        return (
            <Col xs={12}>
                <InputGroup className="mb-3">
                    <FormControl
                        onChange={this.handleChange}
                        value={this.state.link}
                        placeholder="Websocket"
                    />
                    <Button variant="primary" onClick={this.state.isLoading ? null : this.connect} disabled={this.state.isLoading}>
                        { this.state.isLoading ? 'Connecting ... ' : 'Connect' }
                    </Button>
                </InputGroup>
            </Col>
        )
    }
}

export default Websocket