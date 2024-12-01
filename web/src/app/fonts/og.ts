import { WEB_ADDRESS } from "@/config";

export const interRegular = async () => {
  const url = new URL(`${WEB_ADDRESS}/Inter-Regular.ttf`);

  const response = await fetch(url, { cache: "no-cache" });

  const buf = await response.arrayBuffer();

  return buf;
};

export const interBold = async () => {
  const url = new URL(`${WEB_ADDRESS}/Inter-Bold.ttf`);

  const response = await fetch(url, { cache: "no-cache" });

  const buf = await response.arrayBuffer();

  return buf;
};
