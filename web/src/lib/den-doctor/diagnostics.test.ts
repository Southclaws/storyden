import { renderToStaticMarkup } from "react-dom/server";
import { describe, expect, it } from "vitest";

import { diagnoseInstallation } from "./diagnostics";

describe("Den Doctor address checks", () => {
  it.each([
    ["same host", "https://example.com", "https://example.com"],
    ["API subdomain", "https://api.example.com", "https://example.com"],
    [
      "community API",
      "https://api.community.example",
      "https://community.example",
    ],
    [
      "website subdomain",
      "https://example.com",
      "https://community.example.com",
    ],
    ["local ports", "http://localhost:8000", "http://localhost:3000"],
    ["internal hostname", "http://storyden:8000", "http://storyden:3000"],
    ["internal IP", "http://192.168.1.5:8000", "http://192.168.1.5:3000"],
    ["IPv6", "http://[::1]:8000", "http://[::1]:3000"],
    ["private IPv6", "http://[fd00::1]:8000", "http://[fd00::1]:3000"],
    ["loopback IP", "http://127.0.0.1:8000", "http://127.0.0.1:3000"],
    ["default ports", "https://EXAMPLE.com:443/", "https://example.com/"],
  ])("accepts %s", (_name, api, web) => {
    expect(diagnoseInstallation(api, web)).toEqual([]);
  });

  it.each([
    ["https://api.example.com", "http://example.com"],
    ["http://api.example.com", "https://example.com"],
  ])("warns about mixed schemes: %s and %s", (api, web) => {
    expect(diagnoseInstallation(api, web).map((i) => i.code)).toEqual([
      "scheme_mismatch",
    ]);
  });

  it.each([
    ["https://api.example.com", "https://www.example.com"],
    ["https://example.com", "https://example.org"],
    ["https://evil-example.com", "https://example.com"],
    ["https://one.example.co.uk", "https://two.example.co.uk"],
    ["https://one.github.io", "https://two.github.io"],
    ["http://192.168.1.5", "http://192.168.1.6"],
    ["http://[fd00::1]", "http://[fd00::2]"],
    ["http://api.storyden", "http://storyden"],
  ])("rejects cookie hosts: %s and %s", (api, web) => {
    expect(diagnoseInstallation(api, web).map((i) => i.code)).toEqual([
      "cookie_hosts",
    ]);
  });

  it.each(["localhost", "127.0.0.1", "127.1.2.3", "[::1]", "api.localhost"])(
    "explains a remote website pointing at %s",
    (host) => {
      expect(
        diagnoseInstallation(`http://${host}:8000`, "http://example.com").map(
          (i) => i.code,
        ),
      ).toEqual(["localhost_api"]);
    },
  );

  it("does not treat a hostname beginning with 127 as a loopback IP", () => {
    expect(
      diagnoseInstallation("https://127.example.com", "https://example.com"),
    ).toEqual([]);
  });

  it("uses the browser address to detect a remote visitor with localhost defaults", () => {
    expect(
      diagnoseInstallation(
        "http://localhost:8000",
        "http://localhost:3000",
        "http://192.168.1.5:3000",
      ).map((i) => i.code),
    ).toEqual(["localhost_api", "website_origin"]);
  });

  it("compares complete origins, including non-default ports", () => {
    expect(
      diagnoseInstallation(
        "http://localhost:8000",
        "http://localhost:3000",
        "http://localhost:3001",
      ).map((i) => i.code),
    ).toEqual(["website_origin"]);
    expect(
      diagnoseInstallation(
        "https://example.com",
        "https://EXAMPLE.com:443/",
        "https://example.com",
      ),
    ).toEqual([]);
  });

  it("identifies an unsupported website sub-path", () => {
    expect(
      diagnoseInstallation(
        "https://example.com",
        "https://example.com/community",
      ),
    ).toEqual([
      expect.objectContaining({
        code: "website_path",
        category: "routing",
        severity: "error",
      }),
    ]);
  });

  it.each(["/backend", "/api", "/prefix/"])(
    "treats API prefix %s as requiring proxy configuration rather than impossible",
    (path) => {
      const issues = diagnoseInstallation(
        `https://example.com${path}`,
        "https://example.com",
      );

      expect(issues).toEqual([
        expect.objectContaining({
          code: "api_path",
          category: "routing",
          severity: "warning",
        }),
      ]);
      expect(renderToStaticMarkup(issues[0]?.explanation)).toContain(
        `<code>${path}/api</code>`,
      );
    },
  );

  it.each([
    ["https://example.com?mode=api", "https://example.com"],
    ["https://example.com", "https://example.com#community"],
  ])("rejects query strings and fragments: %s and %s", (api, web) => {
    expect(diagnoseInstallation(api, web)).toEqual([
      expect.objectContaining({
        code: "address_suffix",
        category: "routing",
        severity: "error",
      }),
    ]);
  });
});
