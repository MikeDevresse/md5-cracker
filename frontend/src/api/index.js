class Socket {
    constructor(url, callback) {
        this.socket = new WebSocket(url)

        this.socket.onopen = () => {
            console.log("Successfully Connected");
            this.sendMsg("client")
        };

        this.socket.onmessage = msg => {
            callback(msg);
        };

        this.socket.onclose = event => {
            console.log("Socket Closed Connection: ", event);
        };

        this.socket.onerror = error => {
            console.log("Socket Error: ", error);
        };
    }

    sendMsg = msg => {
        this.socket.send(msg);
    };
}

export default Socket;