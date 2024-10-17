"use client";

import { useIsClient } from "@uidotdev/usehooks";
import { PropsWithChildren } from "react";
import { Toaster } from "sonner";
import { SWRConfig } from "swr";

import { AuthProvider } from "src/auth/AuthProvider";

import { cacheProvider } from "@/lib/cache/swr-cache";

export function Providers({ children }: PropsWithChildren) {
  const isClient = useIsClient();

  const provider = isClient ? cacheProvider() : undefined;

  return (
    <AuthProvider>
      <SWRConfig
        value={{
          keepPreviousData: true,
          provider: provider,
        }}
      >
        <Toaster />

        {/* -- */}
        {children}
        {/* -- */}
      </SWRConfig>
    </AuthProvider>
  );
}
