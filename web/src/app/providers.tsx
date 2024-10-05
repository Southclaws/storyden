"use client";

import { PropsWithChildren } from "react";
import { Toaster, toast } from "sonner";

import { AuthProvider } from "src/auth/AuthProvider";

export function Providers({ children }: PropsWithChildren) {
  return (
    <AuthProvider>
      <Toaster />

      {/* -- */}
      {children}
      {/* -- */}
    </AuthProvider>
  );
}
