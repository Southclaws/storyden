"use client";

import { PropsWithChildren, Suspense } from "react";
import { Toaster } from "sonner";
import { SWRConfig } from "swr";

import { useCacheProvider } from "@/lib/cache/swr-cache";
import { DndProvider } from "@/lib/dragdrop/provider";

export function Providers({ children }: PropsWithChildren) {
  const provider = useCacheProvider();

  return (
    <Suspense fallback={children}>
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
    </Suspense>
  );
}
