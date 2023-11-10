"use client";

import { CacheProvider } from "@chakra-ui/next-js";
import { ChakraProvider } from "@chakra-ui/react";
import { PropsWithChildren } from "react";

import { AuthProvider } from "src/auth/AuthProvider";
import { extended } from "src/theme";

// Force chakra to always be light mode - because we're removing it eventually.
type ColourMode = "light" | "dark";
const noopColourModeManager = {
  type: "localStorage" as const,
  ssr: false,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  get(_init?: ColourMode | undefined) {
    return "light" as const;
  },
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  set(_value: "system") {},
};

export function Providers({ children }: PropsWithChildren) {
  return (
    <CacheProvider>
      <ChakraProvider
        theme={extended}
        // We're not using Chakra's reset, instead we're using Panda CSS.
        resetCSS={false}
        colorModeManager={noopColourModeManager}
      >
        <AuthProvider>
          {/* -- */}
          {children}
          {/* -- */}
        </AuthProvider>
      </ChakraProvider>
    </CacheProvider>
  );
}
