import { z } from "zod";

export const DEFAULT_API_ADDRESS = "http://localhost:8000";
export const DEFAULT_WEB_ADDRESS = "http://localhost:3000";

export const ConfigSchema = z.object({
  API_ADDRESS: z.string(),
  WEB_ADDRESS: z.string(),
  source: z.union([z.literal("server"), z.literal("script")]),
});
export type Config = z.infer<typeof ConfigSchema>;

export function serverEnvironment() {
  return {
    API_ADDRESS:
      global.process.env["NEXT_PUBLIC_API_ADDRESS"] ??
      global.process.env["PUBLIC_API_ADDRESS"] ??
      DEFAULT_API_ADDRESS,
    WEB_ADDRESS:
      global.process.env["NEXT_PUBLIC_WEB_ADDRESS"] ??
      global.process.env["PUBLIC_WEB_ADDRESS"] ??
      DEFAULT_WEB_ADDRESS,
    source: "server" as const,
  };
}

function isomorphicEnvironment(): Config {
  if (typeof window !== "undefined") {
    const config = ConfigSchema.parse((window as any).__storyden__);
    console.log("loaded window config", config);
    return config;
  } else {
    const config = serverEnvironment();
    console.log("loaded server config", config);
    return config;
  }
}

const env = isomorphicEnvironment();

export const API_ADDRESS = env.API_ADDRESS;

export const WEB_ADDRESS = env.WEB_ADDRESS;

export function getAPIAddress() {
  if (typeof window !== "undefined") {
    // When called on the client, return the public API address, such as:
    // https://api.mystorydencommunity.com
    return env.API_ADDRESS;
  } else {
    // When called on the server side, we may be running in a container beside
    // the API (the "fullstack" container image) so return the internal address
    // if it's set, otherwise the regular API_ADDRESS configuration value.
    return (
      // The default fullstack image will set this value automatically to :8000.
      global.process.env["SSR_API_ADDRESS"] ??
      // There's no SSR address set so just use the public API address.
      env.API_ADDRESS
    );
  }
}
