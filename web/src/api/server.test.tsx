import { describe, expect, it } from "vitest";

import { forwardSSRNetworkHeaders } from "./server";

describe("forwardSSRNetworkHeaders", () => {
  it("forwards x-forwarded-for unchanged", () => {
    const incoming = new Headers({
      "x-forwarded-for": "203.0.113.1, 203.0.113.2",
    });
    const out = new Headers();

    forwardSSRNetworkHeaders(out, incoming);

    expect(out.get("x-forwarded-for")).toBe("203.0.113.1, 203.0.113.2");
  });

  it("replaces existing x-forwarded-for instead of appending", () => {
    const incoming = new Headers({
      "x-forwarded-for": "203.0.113.1",
    });
    const out = new Headers({
      "x-forwarded-for": "198.51.100.1",
    });

    forwardSSRNetworkHeaders(out, incoming);

    expect(out.get("x-forwarded-for")).toBe("203.0.113.1");
  });

  it("does not forward non-whitelisted headers", () => {
    const incoming = new Headers({
      cookie: "session=abc",
      authorization: "Bearer token",
      "x-custom-ip": "198.51.100.10",
      "x-another-client-ip": "198.51.100.11",
    });
    const out = new Headers();

    forwardSSRNetworkHeaders(out, incoming);

    expect(out.get("cookie")).toBeNull();
    expect(out.get("authorization")).toBeNull();
    expect(out.get("x-custom-ip")).toBeNull();
    expect(out.get("x-another-client-ip")).toBeNull();
  });
});
