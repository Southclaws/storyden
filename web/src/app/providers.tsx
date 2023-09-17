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
