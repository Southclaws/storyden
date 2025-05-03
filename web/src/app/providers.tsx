"use client";

import { PropsWithChildren } from "react";
import { Toaster } from "sonner";
import { SWRConfig } from "swr";

import { AuthProvider } from "src/auth/AuthProvider";

import { useCacheProvider } from "@/lib/cache/swr-cache";
import { DndProvider } from "@/lib/dragdrop/provider";

export function Providers({ children }: PropsWithChildren) {
  const provider = useCacheProvider();

  return (
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
  );
}
