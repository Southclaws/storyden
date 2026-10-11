import type { DiscoveryServiceSettings } from "@/api/openapi-schema";

export const DefaultDiscoverySettings = {
  search_indexing: true,
  training_crawling: false,
  llms_txt: false,
  auth_md: false,
  ai_catalog: false,
  api_catalog: false,
  mcp_server_card: false,
  robots_txt: "",
} satisfies Required<DiscoveryServiceSettings>;

export function resolveDiscoverySettings(settings?: DiscoveryServiceSettings) {
  return { ...DefaultDiscoverySettings, ...settings };
}
