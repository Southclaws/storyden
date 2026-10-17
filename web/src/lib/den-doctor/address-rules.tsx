import type { DoctorContext, DoctorRule } from "./diagnostics";

export const addressRules: readonly DoctorRule[] = [
  {
    code: "scheme_mismatch",
    category: "authentication",
    severity: "error",
    check({ api, web }) {
      if (api.protocol === web.protocol) {
        return;
      }

      return {
        title: "Your website and API use different connection types",
        explanation: (
          <>
            One address uses <code>http://</code> and the other uses{" "}
            <code>https://</code>. Browser security rules can prevent sign-in
            with this combination.
          </>
        ),
        remedy:
          "Set the website and API addresses to use the same connection type: HTTPS for a public community, or HTTP for local and trusted internal networks.",
      };
    },
  },
  {
    code: "localhost_api",
    category: "addresses",
    severity: "warning",
    check(context) {
      if (!hasRemoteLoopbackAPI(context)) {
        return;
      }

      return {
        title: "Your API address points to the visitor’s computer",
        explanation: (
          <>
            A <code>localhost</code> or loopback address connects to the
            computer running the browser, rather than your Storyden server.
          </>
        ),
        remedy: (
          <>
            Change <code>PUBLIC_API_ADDRESS</code> to a hostname or IP address
            that visitors can reach. Container-only addresses belong in{" "}
            <code>SSR_API_ADDRESS</code> instead.
          </>
        ),
      };
    },
  },
  {
    code: "cookie_hosts",
    category: "authentication",
    severity: "error",
    check(context) {
      const { api, web } = context;

      if (hasRemoteLoopbackAPI(context)) {
        return;
      }

      if (canShareCookie(api.hostname, web.hostname)) {
        return;
      }

      return {
        title: "Your website and API cannot share the sign-in cookie",
        explanation: (
          <>
            With these hostnames, Storyden’s session cookie cannot reach both
            the API and the website. This affects unrelated domains and separate
            sibling hosts, such as <code>www.community.example</code> and{" "}
            <code>api.community.example</code>.
          </>
        ),
        remedy: (
          <>
            Use one hostname for both addresses, or place the API below the
            website’s hostname: for example, <code>community.example</code> and{" "}
            <code>api.community.example</code>.
          </>
        ),
      };
    },
  },
  {
    code: "website_origin",
    category: "addresses",
    severity: "warning",
    check({ browser, web }) {
      if (!browser || browser.origin === web.origin) {
        return;
      }

      return {
        title: "You’re visiting a different website address",
        explanation:
          "This browser’s address differs from the website address in your configuration. Sign-in requests from this address may not be accepted.",
        remedy:
          "Open the configured website address, or change the website address in your deployment configuration to the address you want people to use.",
      };
    },
  },
  {
    code: "website_path",
    category: "routing",
    severity: "error",
    check({ web }) {
      if (web.pathname === "/") {
        return;
      }

      return {
        title: "Your website address contains a sub-path",
        explanation: (
          <>
            Storyden’s supplied web app runs at the root of a hostname. Setting
            an address such as <code>example.com/community</code> does not move
            its pages and assets under <code>/community</code>.
          </>
        ),
        remedy: (
          <>
            Host Storyden at a root address, such as{" "}
            <code>community.example.com</code>, and remove the path from the
            website address.
          </>
        ),
      };
    },
  },
  {
    code: "api_path",
    category: "routing",
    severity: "warning",
    check({ api }) {
      if (api.pathname === "/") {
        return;
      }

      return {
        title: "Your API address contains a path prefix",
        explanation: (
          <>
            Storyden adds <code>/api</code> to the API address, so browser
            requests will use <code>{`${api.pathname}/api`}</code>. The backend
            itself serves <code>/api</code>; a custom prefix requires a reverse
            proxy that removes it before forwarding requests.
          </>
        ),
        remedy: (
          <>
            Use the API server’s root address for the simplest setup, without a
            trailing <code>/api</code>. If the prefix is intentional, verify
            your proxy routing for API requests, assets and discovery endpoints.
          </>
        ),
      };
    },
  },
  {
    code: "address_suffix",
    category: "routing",
    severity: "error",
    check({ api, web }) {
      if (!api.search && !api.hash && !web.search && !web.hash) {
        return;
      }

      return {
        title: "Your public addresses contain a query or fragment",
        explanation: (
          <>
            Query strings (after <code>?</code>) and fragments (after{" "}
            <code>#</code>) interfere with the URLs Storyden builds for pages
            and API requests.
          </>
        ),
        remedy:
          "Remove queries and fragments from both public addresses. Use only the server address and, for a deliberately configured API proxy, its path prefix.",
      };
    },
  },
];

function isLoopback(hostname: string): boolean {
  return (
    hostname === "localhost" ||
    hostname.endsWith(".localhost") ||
    /^127(?:\.\d+){3}$/.test(hostname) ||
    hostname === "[::1]"
  );
}

function hasRemoteLoopbackAPI({ api, web, browser }: DoctorContext): boolean {
  return (
    isLoopback(api.hostname) && !isLoopback(browser?.hostname ?? web.hostname)
  );
}

function canShareCookie(apiHost: string, webHost: string): boolean {
  if (apiHost === webHost) {
    return true;
  }

  if ([apiHost, webHost].some((h) => h.startsWith("[") || /^[\d.]+$/.test(h))) {
    return false;
  }

  if ([apiHost, webHost].some((h) => !h.includes(".") && h !== "localhost")) {
    return false;
  }

  return apiHost.endsWith(`.${webHost}`) || webHost.endsWith(`.${apiHost}`);
}
