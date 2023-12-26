"use client";

import { PropsWithChildren } from "react";

import { AuthProvider } from "src/auth/AuthProvider";

export function Providers({ children }: PropsWithChildren) {
  return (
    <AuthProvider>
      {/* -- */}
      {children}
      {/* -- */}
    </AuthProvider>
  );
}
