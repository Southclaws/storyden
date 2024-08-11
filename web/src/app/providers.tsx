"use client";

import { PropsWithChildren } from "react";

import { AuthProvider } from "src/auth/AuthProvider";

import { NavigationProvider } from "@/components/site/Navigation/Right/context";

export function Providers({ children }: PropsWithChildren) {
  return (
    <AuthProvider>
      <NavigationProvider>
        {/* -- */}
        {children}
        {/* -- */}
      </NavigationProvider>
    </AuthProvider>
  );
}
