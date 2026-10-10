import { proxyDiscovery } from "@/lib/discovery/proxy";

export const dynamic = "force-dynamic";

export function GET(request: Request) {
  return proxyDiscovery(request, "/.well-known/api-catalog");
}

export function HEAD(request: Request) {
  return proxyDiscovery(request, "/.well-known/api-catalog");
}

export function OPTIONS(request: Request) {
  return proxyDiscovery(request, "/.well-known/api-catalog");
}
