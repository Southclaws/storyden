import { trimCharsEnd } from "lodash/fp";

const apiAddress = new URL(
  process.env["NEXT_PUBLIC_API_ADDRESS"] ?? "http://localhost:8000",
);

const webAddress = new URL(
  process.env["NEXT_PUBLIC_WEB_ADDRESS"] ?? "http://localhost:3000",
);

const cleanAddress = trimCharsEnd("/");

function getWebsocketAddress(u: URL): string {
  u.protocol = u.protocol.replace("http", "ws"); // https -> wss
  u.pathname = "/ws";

  return u.toString();
}

export const API_ADDRESS = cleanAddress(apiAddress.toString());

export const WEB_ADDRESS = cleanAddress(webAddress.toString());

export const WS_ADDRESS = getWebsocketAddress(apiAddress);
