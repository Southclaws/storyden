import { includes } from "lodash/fp";

import {
  InstanceCapability,
  InstanceCapabilityList,
} from "@/api/openapi-schema";

import { useSettings } from "./settings-client";

const findCapability = (cap: InstanceCapability) => includes(cap);

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
  const found = findCapability(cap)(capabilities);

  return Boolean(found);
}
