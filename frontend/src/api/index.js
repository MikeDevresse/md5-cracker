var socket = new WebSocket("ws://localhost:8080/ws");

let connect = callback => {
    console.log("connecting");

    socket.onopen = () => {
        console.log("Successfully Connected");
        sendMsg("client")
    };

    socket.onmessage = msg => {
        callback(msg);
    };

    socket.onclose = event => {
        console.log("Socket Closed Connection: ", event);
    };

    socket.onerror = error => {
        console.log("Socket Error: ", error);
    };
};

let sendMsg = msg => {
    socket.send(msg);
};

export { connect, sendMsg };