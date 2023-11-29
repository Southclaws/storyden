"use client";

import { CacheProvider } from "@chakra-ui/next-js";
import { ChakraProvider, extendTheme } from "@chakra-ui/react";
import { PropsWithChildren } from "react";

import { AuthProvider } from "src/auth/AuthProvider";

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
        theme={extendTheme({
          styles: {
            // Remove the Chakra defaults - we don't need them with Panda.
            // https://chakra-ui.com/docs/styled-system/global-styles#default-styles
            global: {
              body: {
                fontFamily: "unset",
                color: "unset",
                bg: "unset",
                lineHeight: "unset",
              },

              "*, *::before, &::after": {
                borderColor: "unset",
                wordWrap: "unset",
              },
            },
          },
        })}
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
