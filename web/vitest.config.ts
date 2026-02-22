import path from "node:path";

import { defineConfig } from "vitest/config";

export default defineConfig({
  resolve: {
    alias: [
      {
        find: /^@\/styled-system\/(.*)$/,
        replacement: path.resolve(__dirname, "./styled-system/$1"),
      },
      {
        find: /^styled-system\/(.*)$/,
        replacement: path.resolve(__dirname, "./styled-system/$1"),
      },
      {
        find: "@",
        replacement: path.resolve(__dirname, "./src"),
      },
      {
        find: "src",
        replacement: path.resolve(__dirname, "./src"),
      },
    ],
  },
  test: {
    environment: "jsdom",
    setupFiles: ["./vitest.setup.ts"],
    include: ["src/**/*.test.tsx"],
    clearMocks: true,
    restoreMocks: true,
    mockReset: true,
  },
});
