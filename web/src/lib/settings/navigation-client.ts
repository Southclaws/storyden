"use client";

import { handle } from "@/api/client";

import { useSettingsMutation } from "./mutation";
import { DefaultNavigationConfig, type NavigationConfig } from "./navigation";
import { type Settings } from "./settings";
import { useSettings } from "./settings-client";

export function useNavigationConfig(
  initialSettings?: Settings,
  revalidateOnMount = false,
): NavigationConfig {
  const { settings } = useSettings(initialSettings, revalidateOnMount);

  return (
    settings?.metadata.navigation ??
    initialSettings?.metadata.navigation ??
    DefaultNavigationConfig
  );
}

export function useNavigationMutation() {
  const { updateSettings } = useSettingsMutation();

  async function updateNavigation(navigation: NavigationConfig) {
    await handle(
      async () => {
        await updateSettings({
          metadata: {
            navigation,
          },
        });
      },
      {
        promiseToast: {
          loading: "Updating navigation...",
          success: "Navigation updated",
        },
      },
    );
  }

  return { updateNavigation };
}
