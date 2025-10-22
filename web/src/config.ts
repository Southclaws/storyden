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
    const parsed = ConfigSchema.safeParse((window as any).__storyden__);

    // Don't bail out if the config is invalid, just log it and use the default.
    // This is a rare occurrance but it does happen when Next.js renders either
    // not-found.tsx or error.tsx as it does not load the <script> tag for some
    // weird unknown reason. This seems like a bug in Next.js and not Storyden.
    // Either way, the config isn't necessary on either of those pages as they
    // don't make client side API calls, they just render a basic error page.
    if (!parsed.success) {
      console.error(
        `Invalid config loaded from \`window.__storyden__\`, this indicates a problem running the root layout <script> tag on the server render which should inject the API_ADDRESS and WEB_ADDRESS environment variables.

A default configuration will be used, however this configuration will most likely not work correctly in most production environments.

If you see this, please open an issue at https://github.com/Southclaws/storyden/issues/new
`,
        parsed.error.issues,
      );

      return {
        API_ADDRESS: DEFAULT_API_ADDRESS,
        WEB_ADDRESS: DEFAULT_WEB_ADDRESS,
        source: "server",
      };
    }

    const config = parsed.data;
    console.log("loaded window config", config);
    return config;
  } else {
    const config = serverEnvironment();
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
