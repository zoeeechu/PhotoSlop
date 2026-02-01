import type { GamePacket, MicrogameStartPacket } from "../types/packets";
import { writable } from "svelte/store";

export const availableIcons = writable<string[]>([]);
export const nameValidated = writable<string | null>(null);
export const microgameData = writable<MicrogameStartPacket | null>(null);
export const playerGuid = writable<string | null>(null);
export const displayText = writable<string | null>(null);

export function handlePacket(packet: GamePacket) {
  console.log("From server:", packet.type);
  switch (packet.type) {
    case "icons-available":
      console.log(packet.icons);
      availableIcons.set(packet.icons);
      break;

    case "icon-disabled":
      console.log(packet.icon);
      availableIcons.update(icons => icons.filter(i => i !== packet.icon));
      break;

    case "name-validation":
      console.log(packet.guid);
      nameValidated.set(packet.guid);
      playerGuid.set(packet.guid);
      break;

    case "microgame-time-up": 
      console.log("Microgame: time up");
      break;
    
    case "microgame-end": 
      console.log("Microgame: end");
      microgameData.set(null);
      break;

    case "microgame-start":
      console.log(packet.mask, packet.time);
      microgameData.set(packet);
      break;

    case "display-text":
      console.log(packet.text);
      displayText.set(packet.text);
      break;
  }
}
