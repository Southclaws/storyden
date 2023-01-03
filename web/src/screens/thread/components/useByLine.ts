import Identicon from "identicon.js";
import { useEffect, useState } from "react";

export function useByLine(handle: string) {
  const [fallback, setFallback] = useState("");

  useEffect(() => {
    (async () => {
      const hash = await sha256(handle);
      const fallback = new Identicon(hash, 420).toString();

      setFallback(`data:image/png;base64,${fallback}`);
    })().catch(console.error);
  }, [handle]);

  return { fallback, src: `/api/v1/accounts/${handle}/avatar` };
}

async function sha256(source: string) {
  const sourceBytes = new TextEncoder().encode(source);
  const digest = await crypto.subtle.digest("SHA-256", sourceBytes);
  const resultBytes = [...new Uint8Array(digest)];
  return resultBytes.map((x) => x.toString(16).padStart(2, "0")).join("");
}
