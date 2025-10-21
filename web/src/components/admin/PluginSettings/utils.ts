import {
  Plugin,
  PluginActiveState,
  PluginStatusActive,
  PluginStatusError,
  PluginStatusInactive,
} from "@/api/openapi-schema";

export function getPluginActiveState(plugin: Plugin): PluginActiveState {
  if (isPluginStatusActive(plugin.status)) {
    return PluginActiveState.active;
  } else if (isPluginStatusError(plugin.status)) {
    return PluginActiveState.error;
  }
  return PluginActiveState.inactive;
}

export function isPluginStatusActive(
  status: Plugin["status"],
): status is PluginStatusActive {
  return "activated_at" in status;
}

export function isPluginStatusInactive(
  status: Plugin["status"],
): status is PluginStatusInactive {
  return "deactivated_at" in status;
}

export function isPluginStatusError(
  status: Plugin["status"],
): status is PluginStatusError {
  return "message" in status;
}
