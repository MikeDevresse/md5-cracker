import "./ResultHistory.scss"
import {Component} from "react";

class ResultHistory extends Component {
    keyCounter = 0;

    render() {
        const history = this.props.resultHistory.map(res => {
            console.log(res)
            return <tr key={this.keyCounter++}><td>{res.hash}</td><td>{res.result}</td></tr>
        });

        return (
            <div className="resultHistory card">
                <div className="card-header">Result History</div>
                <table>
                    <thead>
                    <tr>
                        <th>Hash</th>
                        <th>Value</th>
                    </tr>
                    </thead>
                    <tbody>
                        {history}
                    </tbody>
                </table>
            </div>
        )
    }
}

export default ResultHistory