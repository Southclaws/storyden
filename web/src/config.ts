import { z } from "zod";

const DEFAULT_API_ADDRESS = "http://localhost:8000";
const DEFAULT_WEB_ADDRESS = "http://localhost:3000";

export const ConfigSchema = z.object({
  API_ADDRESS: z.string(),
  WEB_ADDRESS: z.string(),
  source: z.union([z.literal("server"), z.literal("script")]),
});
export type Config = z.infer<typeof ConfigSchema>;

function isomorphicEnvironment(): Config {
  if (typeof window !== "undefined") {
    const config = ConfigSchema.parse((window as any).__storyden__);
    console.log("loaded window config", config);
    return config;
  } else {
    const config = {
      API_ADDRESS:
        global.process.env["NEXT_PUBLIC_API_ADDRESS"] ?? DEFAULT_API_ADDRESS,
      WEB_ADDRESS:
        global.process.env["NEXT_PUBLIC_WEB_ADDRESS"] ?? DEFAULT_WEB_ADDRESS,
      source: "server" as const,
    };
    console.log("loaded server config", config);
    return config;
  }
}

const env = isomorphicEnvironment();

export const API_ADDRESS = env.API_ADDRESS;

export const WEB_ADDRESS = env.WEB_ADDRESS;
