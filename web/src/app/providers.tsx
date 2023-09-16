"use client";

import { CacheProvider } from "@chakra-ui/next-js";
import { ChakraProvider } from "@chakra-ui/react";
import { PropsWithChildren } from "react";

import { InfoProvider } from "src/api/InfoProvider/InfoProvider";
import { AuthProvider } from "src/auth/AuthProvider";
import { extended } from "src/theme";

export function Providers({ children }: PropsWithChildren) {
  return (
    <CacheProvider>
      <ChakraProvider
        theme={extended}
        // We're not using Chakra's reset, instead we're using Panda CSS.
        resetCSS={false}
        // Similarly to above, we're using our own theming via Panda CSS.
        disableGlobalStyle={true}
      >
        <AuthProvider>
          <InfoProvider>
            {/* -- */}
            {children}
            {/* -- */}
          </InfoProvider>
        </AuthProvider>
      </ChakraProvider>
    </CacheProvider>
  );
}
