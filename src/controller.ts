let socket: WebSocket;

function getWsUrl() {
    const params = new URLSearchParams(window.location.search);
    const wsFromQuery = params.get("ws");

    if (wsFromQuery) {
        console.log("Using WS from query:", wsFromQuery);
        return wsFromQuery;
    }

    return (location.protocol === "https:" ? "wss://" : "ws://") +
           location.host +
           "/game";
}

function connect() {
    const wsUrl = getWsUrl();
    socket = new WebSocket(wsUrl);

    console.log("Connecting to:", wsUrl);

    socket.onopen = () => {
        document.getElementById("status")!.textContent = "Connected ✓";
        (document.getElementById("status") as HTMLElement).style.color = "green";
        (document.getElementById("leftBtn") as HTMLButtonElement).disabled = false;
        (document.getElementById("rightBtn") as HTMLButtonElement).disabled = false;
    };

    socket.onmessage = e => console.log("From Unity:", e.data);

    socket.onerror = () => {
        document.getElementById("status")!.textContent = "Connection error";
        (document.getElementById("status") as HTMLElement).style.color = "red";
    };

    socket.onclose = () => {
        document.getElementById("status")!.textContent = "Disconnected – Reconnecting...";
        (document.getElementById("status") as HTMLElement).style.color = "orange";
        (document.getElementById("leftBtn") as HTMLButtonElement).disabled = true;
        (document.getElementById("rightBtn") as HTMLButtonElement).disabled = true;
        setTimeout(connect, 2000);
    };
}

function sendText() {
    const text = (document.getElementById("msg") as HTMLInputElement).value;
    socket.send(text);
}

(window as any).move = (direction: string) => {
    socket.send(direction);
};

(window as any).sendText = sendText;

connect();
