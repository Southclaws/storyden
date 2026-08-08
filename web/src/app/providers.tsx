"use client";

import { PropsWithChildren } from "react";
import { NuqsAdapter } from "nuqs/adapters/next/app";
import { Toaster } from "sonner";
import { SWRConfig } from "swr";

import { AuthProvider } from "@/auth/AuthProvider";
import { useCacheProvider } from "@/lib/cache/swr-cache";
import { DndProvider } from "@/lib/dragdrop/provider";

export function Providers({ children }: PropsWithChildren) {
  const provider = useCacheProvider();

  return (
    <NuqsAdapter>
      <AuthProvider>
        <SWRConfig
          value={{
            keepPreviousData: true,
            // provider: provider,
          }}
        >
          <DndProvider>
            <Toaster />

            {/* -- */}
            {children}
            {/* -- */}
          </DndProvider>
        </SWRConfig>
      </AuthProvider>
    </NuqsAdapter>
  );
}
