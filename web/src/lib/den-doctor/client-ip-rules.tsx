import type { DoctorRule } from "./diagnostics";
import { isPublicSingleHost, parseForwardedFor } from "./ip";

export const clientIPRules: readonly DoctorRule[] = [
  {
    code: "client_ip_proxy_address",
    category: "networking",
    severity: "error",
    check({ clientIP }) {
      if (clientIP?.mode !== "xff_trusted_proxies") {
        return;
      }

      const hops = parseForwardedFor(clientIP.forwardedFor);
      const proxy = hops.at(-1);

      if (hops.length < 2 || !proxy || clientIP.resolvedAddress !== proxy) {
        return;
      }

      return {
        title: "Visitors are identified by a proxy’s address",
        explanation: (
          <>
            Storyden identified your request as <code>{proxy}</code>, which is
            the last proxy in its <code>X-Forwarded-For</code> chain, not your
            address. Everyone routed through that proxy shares one rate limit,
            so a busy period can block all of them at once.
          </>
        ),
        remedy: (
          <>
            Add your proxy’s network to <strong>Trusted proxy CIDRs</strong>.
            Hosting providers rotate their edge proxies, so trust the provider’s
            published IP ranges rather than <code>{proxy}</code> alone.
          </>
        ),
      };
    },
  },
  {
    code: "client_ip_single_proxy_hosts",
    category: "networking",
    severity: "warning",
    check({ clientIP }) {
      if (clientIP?.mode !== "xff_trusted_proxies") {
        return;
      }

      const hosts = clientIP.trustedProxyCIDRs.filter(isPublicSingleHost);

      if (!hosts.length) {
        return;
      }

      return {
        title: "Your trusted proxies are single public addresses",
        explanation: (
          <>
            {hosts.map((host, index) => (
              <span key={host}>
                {index > 0 && ", "}
                <code>{host}</code>
              </span>
            ))}{" "}
            {hosts.length === 1 ? "is" : "are"} trusted as{" "}
            {hosts.length === 1 ? "a single address" : "single addresses"}. Edge
            networks such as Fly.io and Cloudflare route visitors through many
            proxies and change them over time. Visitors who arrive through any
            other proxy are identified by that proxy’s address and share its
            rate limit.
          </>
        ),
        remedy: (
          <>
            Replace single proxy addresses in{" "}
            <strong>Trusted proxy CIDRs</strong> with your provider’s published
            IP ranges.
          </>
        ),
      };
    },
  },
  {
    code: "client_ip_ignores_forwarding",
    category: "networking",
    severity: "warning",
    check({ clientIP }) {
      if (
        clientIP?.mode !== "remote_addr" ||
        !parseForwardedFor(clientIP.forwardedFor).length
      ) {
        return;
      }

      return {
        title: "Storyden is ignoring your proxy’s forwarded addresses",
        explanation: (
          <>
            Requests arrive with an <code>X-Forwarded-For</code> header, but the
            client IP mode is <strong>Raw IP</strong>. Every visitor is
            identified by your proxy’s address and shares its rate limit.
          </>
        ),
        remedy: (
          <>
            Set <strong>Client IP mode</strong> to{" "}
            <strong>X-Forwarded-For with trusted proxy CIDRs</strong> and add
            your proxy’s IP ranges.
          </>
        ),
      };
    },
  },
];
