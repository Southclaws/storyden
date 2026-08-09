import { describe, expect, it } from "vitest";

import { isPrivatePage } from "./useAuthProvider";

describe("isPrivatePage", () => {
  it.each(["/settings", "/new", "/admin"])(
    "matches the private route %s",
    (pathname) => {
      expect(isPrivatePage(pathname)).toBe(true);
    },
  );

  it.each([
    "/settings/account",
    "/settings/account/security",
    "/new/thread",
    "/admin/settings",
  ])("matches the private descendant route %s", (pathname) => {
    expect(isPrivatePage(pathname)).toBe(true);
  });

  it.each(["/", "/settings-updates", "/newest", "/administrator"])(
    "does not match the public route %s",
    (pathname) => {
      expect(isPrivatePage(pathname)).toBe(false);
    },
  );
});
