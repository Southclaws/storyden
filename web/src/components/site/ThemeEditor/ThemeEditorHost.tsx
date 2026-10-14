"use client";

import dynamic from "next/dynamic";

import { useAccountGet } from "@/api/openapi-client/accounts";
import { Permission } from "@/api/openapi-schema";
import { useThemeEditingEnabled } from "@/lib/theme/theme-editing";
import { hasPermission } from "@/utils/permissions";

const ThemeEditorController = dynamic(
  () =>
    import("./ThemeEditorController").then(
      (module) => module.ThemeEditorController,
    ),
  { ssr: false },
);

export function ThemeEditorHost() {
  const editing = useThemeEditingEnabled();
  const account = useAccountGet();

  if (!editing || !hasPermission(account.data, Permission.ADMINISTRATOR)) {
    return null;
  }

  return <ThemeEditorController />;
}
