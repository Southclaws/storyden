const BASE32_CHARS = "0123456789abcdefghijklmnopqrstuv";

function encodeBase32(bytes: Uint8Array): string {
  let result = "";
  let bits = 0;
  let value = 0;

  for (let i = 0; i < bytes.length; i++) {
    const byte = bytes[i];
    if (byte === undefined) continue;
    value = (value << 8) | byte;
    bits += 8;

    while (bits >= 5) {
      result += BASE32_CHARS[(value >>> (bits - 5)) & 31];
      bits -= 5;
    }
  }

  if (bits > 0) {
    result += BASE32_CHARS[(value << (5 - bits)) & 31];
  }

  return result;
}

export function generateXid(): string {
  const uuid = crypto.randomUUID();
  const bytes = new Uint8Array(
    uuid
      .replace(/-/g, "")
      .match(/.{2}/g)!
      .map((byte) => parseInt(byte, 16)),
  );

  const timestamp = Math.floor(Date.now() / 1000);
  const timestampBytes = new Uint8Array([
    (timestamp >>> 24) & 0xff,
    (timestamp >>> 16) & 0xff,
    (timestamp >>> 8) & 0xff,
    timestamp & 0xff,
  ]);

  const combined = new Uint8Array(12);
  combined.set(timestampBytes, 0);
  combined.set(bytes.slice(0, 8), 4);

  return encodeBase32(combined);
}
