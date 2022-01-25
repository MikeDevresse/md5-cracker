import "./Configuration.scss"
import {Component} from "react";
import {Button, Card, Col, FormControl, Table} from "react-bootstrap";
import {sendMsg} from "../../../api";

class Configuration extends Component {
    constructor(props) {
        super(props);
        this.state = {
            maxSearch: props.maxSearch,
            slaves: props.slaves,
            maxSlavesPerRequest: props.maxSlavesPerRequest,
        }
        this.handleInputChange = this.handleInputChange.bind(this);
        this.setMaxSearch = this.setMaxSearch.bind(this);
        this.setSlaves = this.setSlaves.bind(this);
        this.setMaxSlavesPerRequest = this.setMaxSlavesPerRequest.bind(this);
    }

    componentDidUpdate(prevProps, prevState) {
        if (prevProps !== this.props) {
            this.setState({
                maxSearch: this.props.maxSearch,
                slaves: this.props.slaves,
                maxSlavesPerRequest: this.props.maxSlavesPerRequest,
            })
        }
    }


    handleInputChange(event) {
        const target = event.target;
        const value = target.type === 'checkbox' ? target.checked : target.value;
        const name = target.name;

        this.setState({
            [name]: value
        });
    }

    setMaxSearch(event) {
        event.preventDefault()
        sendMsg("max-search " + this.state.maxSearch)
    }

    setSlaves(event) {
        event.preventDefault()
        sendMsg("slaves " + this.state.slaves)
    }

    setMaxSlavesPerRequest(event) {
        event.preventDefault()
        sendMsg("max-slaves-per-request " + this.state.maxSlavesPerRequest)
    }

    stopAll(event) {
        event.preventDefault()
        sendMsg("stop-all")
    }

    render() {
        return (
            <Col xs={12} lg={6} className="configuration">
                <Card>
                    <Card.Header>Configuration</Card.Header>
                    <Table className="mb-0">
                        <tbody>
                        <tr>
                            <td>Max search</td>
                            <td>
                                <FormControl
                                    name="maxSearch"
                                    type="text"
                                    value={this.state.maxSearch}
                                    onChange={this.handleInputChange}
                                />
                            </td>
                            <td>
                                <Button onClick={this.setMaxSearch}>
                                    Update
                                </Button>
                            </td>
                        </tr>
                        <tr>
                            <td>Slaves</td>
                            <td>
                                <FormControl
                                    name="slaves"
                                    type="number"
                                    value={this.state.slaves}
                                    onChange={this.handleInputChange}
                                />
                            </td>
                            <td>
                                <Button onClick={this.setSlaves}>
                                    Update
                                </Button>
                            </td>
                        </tr>
                        <tr>
                            <td>Max slave per request</td>
                            <td>
                                <FormControl
                                    name="maxSlavesPerRequest"
                                    type="number"
                                    value={this.state.maxSlavesPerRequest}
                                    onChange={this.handleInputChange}
                                />
                            </td>
                            <td>
                                <Button onClick={this.setMaxSlavesPerRequest}>
                                    Update
                                </Button>
                            </td>
                        </tr>
                        <tr>
                            <td colSpan="3"><Button variant="danger" className="w-100" onClick={this.stopAll}>STOP ALL</Button></td>
                        </tr>
                        </tbody>
                    </Table>
                </Card>
            </Col>
        )
    }
}

export default Configuration