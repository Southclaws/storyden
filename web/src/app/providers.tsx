"use client";

import { CacheProvider } from "@chakra-ui/next-js";
import { ChakraProvider } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import { AuthProvider } from "src/auth/AuthProvider";
import { extended } from "src/theme";

export function Providers({ children }: PropsWithChildren) {
  return (
    <CacheProvider>
      <ChakraProvider theme={extended}>
        <AuthProvider>{children}</AuthProvider>
      </ChakraProvider>
    </CacheProvider>
  );
}
