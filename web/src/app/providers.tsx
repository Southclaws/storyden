"use client";

import { NuqsAdapter } from "nuqs/adapters/next/app";
import { PropsWithChildren } from "react";
import { Toaster } from "sonner";
import { SWRConfig } from "swr";

import { AuthProvider } from "@/auth/AuthProvider";
import { ThemeEditorHost } from "@/components/site/ThemeEditor";
import { useCacheProvider } from "@/lib/cache/swr-cache";
import { DndProvider } from "@/lib/dragdrop/provider";
import { ThemeRuntime } from "@/lib/theme/ThemeRuntime";

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
            <ThemeRuntime />
            <Toaster />

            {/* -- */}
            {children}
            <ThemeEditorHost />
            {/* -- */}
          </DndProvider>
        </SWRConfig>
      </AuthProvider>
    </NuqsAdapter>
  );
}
