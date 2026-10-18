import { API_ADDRESS } from "@/config";

export function getIconURL(size: "32x32" | "120x120" | "512x512") {
  return `${API_ADDRESS}/api/info/icon/${size}`;
}

export function getBannerURL(variant?: "s" | "m" | "l") {
  return `${API_ADDRESS}/api/info/banner${variant ? `?variant=${variant}` : ""}`;
}
