// src/lib/types/packets.ts

export type PacketBase = {
  type: string;
  from: "server" | "client";
};

export type IconsAvailablePacket = {
  type: "icons-available";
  from: "server";
  icons: string[];
};

export type IconSelectionPacket = {
  type: "icon-selection";
  from: "client";
  icon: string;
};

export type IconDisabledPacket = {
  type: "icon-disabled";
  from: "server";
  icon: string;
};

export type NameEntryPacket = {
  type: "name-entry";
  from: "client";
  name: string;
};

export type NameValidationPacket = {
  type: "name-validation";
  from: "server";
  guid: string;
};

export type MicrogameStartPacket = {
  type: "microgame-start";
  from: "server";
  time: number;
  mask: { x: number; y: number }[];
  direction: string;
};

export type MicrogameTimeUpPacket = {
  type: "microgame-time-up";
  from: "server";
};

export type MicrogameResultPacket = {
  type: "microgame-result";
  from: "client";
  id: string;
  result: { x: number; y: number }[];
};

export type MicrogameEndPacket = {
  type: "microgame-end";
  from: "server";
};

export type DisplayTextPacket = {
  type: "display-text";
  from: "server";
  text: string;
};

export type GamePacket =
  | IconsAvailablePacket
  | IconSelectionPacket
  | IconDisabledPacket
  | NameEntryPacket
  | NameValidationPacket
  | MicrogameStartPacket
  | MicrogameTimeUpPacket
  | MicrogameResultPacket
  | MicrogameEndPacket
  | DisplayTextPacket;
