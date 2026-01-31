<script lang="ts">
  import { onMount } from "svelte";
  import StatusBanner from "./StatusBanner.svelte";
  import CharacterSelector from "./CharacterSelector.svelte";
  import type { GamePacket } from "./lib/types/packets";
  import { handlePacket } from "./lib/network/handlePacket";
  let socket: WebSocket;
  let status = "Connecting...";
  let statusVisible = false;
  let statusTimeout: ReturnType<typeof setTimeout> | undefined;
  let statusColor = "black";
  let buttonsDisabled = true;
  let message = {};

  let page = "enter";
  let name = "";

  function getWsUrl() {
    const params = new URLSearchParams(window.location.search);
    const wsFromQuery = params.get("ws");
    if (wsFromQuery) {
      console.log("Using WS from query:", wsFromQuery);
      return wsFromQuery;
    }
    return (
      (location.protocol === "https:" ? "wss://" : "ws://") +
      location.host +
      "/game"
    );
  }
  function connect() {
    const wsUrl = getWsUrl();
    socket = new WebSocket(wsUrl);
    console.log("Connecting to:", wsUrl);

    socket.onopen = () => {
      status = "Connected ✓";
      statusColor = "green";
      statusVisible = true;
      buttonsDisabled = false;
    };

    //socket.onmessage = (e) => console.log("From Unity:", e.data);

    socket.onmessage = (event) => {
      const packet: GamePacket = JSON.parse(event.data);
      handlePacket(packet);
    };

    socket.onerror = () => {
      status = "Connection error";
      statusColor = "red";
      statusVisible = true;
    };

    socket.onclose = () => {
      if (status !== "Disconnected – Reconnecting...") {
        status = "Disconnected – Reconnecting...";
        statusColor = "orange";
      }
      statusVisible = true;
      buttonsDisabled = true;
      setTimeout(connect, 2000);
    };
  }

  function move(direction: string) {
    socket.send(direction);
  }

  function sendText() {
    message = {
      type: "entry",
      name: name,
    };
    page = "game";
    socket.send(JSON.stringify(message));
  }

  onMount(() => {
    connect();
    return () => socket?.close();
  });
</script>

<StatusBanner {status} {statusColor} bind:statusVisible />

{#if page === "enter"}
  <div class="gap-2 p-5 mt-2 flex flex-col items-center justify-center">
    <div class="flex items-center justify-center">
      <img src="/images/logo.svg" alt="My App Logo" class="logo w-30 mt-5" />
    </div>
    <div class="flex flex-row gap-3">
      <h1 class="text-3xl">pick your character:</h1>
    </div>
    <div>
    <CharacterSelector />

    </div>
    <div class="flex flex-row gap-3">
      <h1 class="text-3xl">Enter your name:</h1>
    </div>
    <div class="flex flex-row gap-2">
      <!-- <button
      class="border p-2 disabled:bg-gray-300 disabled:border-gray-500"
      on:click={() => move("left")}
      disabled={buttonsDisabled}>LEFT</button
    >
    <button
      class="border p-2 disabled:bg-gray-300 disabled:border-gray-500"
      on:click={() => move("right")}
      disabled={buttonsDisabled}>RIGHT</button
    > -->
      <form class="flex gap-2" on:submit|preventDefault={sendText}>
        <input class="border p-2 rounded-lg" bind:value={name} />
        <button class="border p-2 rounded-lg" type="submit">Start</button>
      </form>
    </div>
  </div>
{:else}
  <div class="gap-2 p-5 mt-2 flex flex-col items-center justify-center">
    <div class="flex flex-row items-center gap-4">
      <div class="flex items-center justify-center">
        <img src="/images/logo.svg" alt="My App Logo" class="logo w-20 mt-5" />
      </div>
      <h1 class="text-5xl">Hey, {name}!</h1>
    </div>

    <p>Look at the main screen...</p>
  </div>
{/if}
