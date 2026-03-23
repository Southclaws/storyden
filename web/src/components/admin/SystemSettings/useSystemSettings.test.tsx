import { describe, expect, it } from "vitest";

import {
  buildClientIPSettingsPayload,
  parseTrustedProxyCidrs,
} from "./useSystemSettings";

describe("useSystemSettings helpers", () => {
  it("parses trusted proxy CIDRs from comma/newline input", () => {
    expect(parseTrustedProxyCidrs("10.0.0.0/8,\n 172.16.0.0/12 ,, ")).toEqual([
      "10.0.0.0/8",
      "172.16.0.0/12",
    ]);
  });

  it("sends only mode for remote_addr", () => {
    expect(
      buildClientIPSettingsPayload(
        { client_ip_mode: "remote_addr", client_ip_header: "X-Real-IP" },
        ["10.0.0.0/8"],
      ),
    ).toEqual({ client_ip_mode: "remote_addr" });
  });

  it("sends only header for single_header mode", () => {
    expect(
      buildClientIPSettingsPayload(
        {
          client_ip_mode: "single_header",
          client_ip_header: "CF-Connecting-IP",
        },
        ["10.0.0.0/8"],
      ),
    ).toEqual({
      client_ip_mode: "single_header",
      client_ip_header: "CF-Connecting-IP",
    });
  });

  it("sends only CIDRs for xff_trusted_proxies mode", () => {
    expect(
      buildClientIPSettingsPayload(
        {
          client_ip_mode: "xff_trusted_proxies",
          client_ip_header: "X-Real-IP",
        },
        ["10.0.0.0/8", "172.16.0.0/12"],
      ),
    ).toEqual({
      client_ip_mode: "xff_trusted_proxies",
      trusted_proxy_cidrs: ["10.0.0.0/8", "172.16.0.0/12"],
    });
  });
});
