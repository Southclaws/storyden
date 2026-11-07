import { build } from "esbuild";

await build({
  entryPoints: ["src/workers/notifications.ts"],
  outfile: "public/sw.js",
  bundle: true,
  format: "iife",
  platform: "browser",
  target: "es2022",
  tsconfig: "src/workers/tsconfig.json",
});
