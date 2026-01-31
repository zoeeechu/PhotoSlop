<script lang="ts">
let socket: WebSocket;
let status = "Connecting...";
let statusColor = "black";
let buttonsDisabled = true;
let message = "";

function getWsUrl() {
    const params = new URLSearchParams(window.location.search);
    const wsFromQuery = params.get("ws");
    if (wsFromQuery) {
        console.log("Using WS from query:", wsFromQuery);
        return wsFromQuery;
    }
    return (location.protocol === "https:" ? "wss://" : "ws://") + location.host + "/game";
}

function connect() {
    const wsUrl = getWsUrl();
    socket = new WebSocket(wsUrl);
    console.log("Connecting to:", wsUrl);

    socket.onopen = () => {
        status = "Connected ✓";
        statusColor = "green";
        buttonsDisabled = false;
    };

    socket.onmessage = e => console.log("From Unity:", e.data);

    socket.onerror = () => {
        status = "Connection error";
        statusColor = "red";
    };

    socket.onclose = () => {
        status = "Disconnected – Reconnecting...";
        statusColor = "orange";
        buttonsDisabled = true;
        setTimeout(connect, 2000);
    };
}

function move(direction: string) {
    socket.send(direction);
}

function sendText() {
    socket.send(message);
}

connect();
</script>

<div style="font-family:sans-serif; text-align:center;">
    <h1>Controller</h1>
    <div style="color: {statusColor}">{status}</div>
    <button on:click={() => move('left')} disabled={buttonsDisabled}>LEFT</button>
    <button on:click={() => move('right')} disabled={buttonsDisabled}>RIGHT</button>
    <input bind:value={message}>
    <button on:click={sendText}>Send</button>
</div>
