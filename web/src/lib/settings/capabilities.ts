import {
  InstanceCapability,
  InstanceCapabilityList,
} from "@/api/openapi-schema";

import { useSettings } from "./settings-client";

export function useCapabilities() {
  const { settings } = useSettings();

  return settings?.capabilities ?? [];
}

export function useCapability(cap: InstanceCapability) {
  const { settings } = useSettings();

  return hasCapability(cap, settings?.capabilities);
}

export function hasCapability(
  cap: InstanceCapability,
  cs?: InstanceCapabilityList,
) {
  const capabilities = cs ?? [];
  return capabilities.includes(cap);
}
