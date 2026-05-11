"use client";

import { PropsWithChildren } from "react";
import { Toaster } from "sonner";
import { SWRConfig } from "swr";

import { AuthProvider } from "src/auth/AuthProvider";

import { NextDevtoolsI18n } from "@/components/site/NextDevtoolsI18n";
import { Locale } from "@/i18n/config";
import { I18nProvider } from "@/i18n/provider";
import { useCacheProvider } from "@/lib/cache/swr-cache";
import { DndProvider } from "@/lib/dragdrop/provider";

type Props = PropsWithChildren<{
  initialLocale: Locale;
}>;

export function Providers({ children, initialLocale }: Props) {
  const provider = useCacheProvider();

  return (
    <I18nProvider initialLocale={initialLocale}>
      <AuthProvider>
        <SWRConfig
          value={{
            keepPreviousData: true,
            // provider: provider,
          }}
        >
          <DndProvider>
            <Toaster />
            <NextDevtoolsI18n />

            {/* -- */}
            {children}
            {/* -- */}
          </DndProvider>
        </SWRConfig>
      </AuthProvider>
    </I18nProvider>
  );
}
