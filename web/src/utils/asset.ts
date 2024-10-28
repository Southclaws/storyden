import { API_ADDRESS } from "@/config";

export function getAssetURL(filename?: string) {
  if (!filename) return;

  return `${API_ADDRESS}${filename}`;
}
