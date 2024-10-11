"use client";

import { PropsWithChildren } from "react";
import { Toaster } from "sonner";
import { SWRConfig } from "swr";

import { AuthProvider } from "src/auth/AuthProvider";

export function Providers({ children }: PropsWithChildren) {
  return (
    <AuthProvider>
      <SWRConfig
        value={{
          keepPreviousData: true,
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
