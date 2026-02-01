import type { GamePacket } from "../types/packets";

export function handlePacket(packet: GamePacket) {
  console.log("From server:", packet.type);
  switch (packet.type) {
    case "icons-available":
      console.log(packet.icons);
      break;

    case "icon-disabled":
      console.log(packet.icon);
      break;

    case "name-validation":
      console.log(packet.guid);
      break;

    case "microgame-time-up": 
      console.log("Microgame: time up");
      break;
    
    case "microgame-end": 
      console.log("Microgame: end");
      break;

    case "microgame-start":
      console.log(packet.mask, packet.time);
      break;

    case "display-text":
      console.log(packet.text);
      break;
  }
}
