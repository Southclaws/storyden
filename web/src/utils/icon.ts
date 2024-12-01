import { API_ADDRESS } from "@/config";

export function getIconURL(size: "512x512") {
  return `${API_ADDRESS}/api/info/icon/${size}`;
}

export function getBannerURL() {
  return `${API_ADDRESS}/api/info/banner`;
}
