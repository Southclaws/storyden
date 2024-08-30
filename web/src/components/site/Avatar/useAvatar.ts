import { API_ADDRESS } from "@/config";

export function useAvatar(handle: string) {
  return { src: `${API_ADDRESS}/api/accounts/${handle}/avatar` };
}
