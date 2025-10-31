import { useCallback, useState } from "react";

export function useCopyToClipboard(): [
  string | null,
  (text: string) => Promise<void>,
  boolean,
] {
  const [copiedText, setCopiedText] = useState<string | null>(null);
  const isClipboardAvailable = !!navigator?.clipboard;

  const copyToClipboard = useCallback(async (text: string) => {
    if (!navigator?.clipboard) {
      console.warn("Clipboard not supported");
      return;
    }

    try {
      await navigator.clipboard.writeText(text);
      setCopiedText(text);
    } catch (error) {
      console.error("Failed to copy text: ", error);
    }
  }, []);

  return [copiedText, copyToClipboard, isClipboardAvailable];
}
