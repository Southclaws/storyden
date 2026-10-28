import { describe, expect, it } from "vitest";

import { getSafeReturnPath } from "./return-path";

describe("getSafeReturnPath", () => {
  it.each([
    "/",
    "/settings/account",
    "/t/welcome?sort=new#reply",
    "/oauth/authorise?redirect_uri=https%3A%2F%2Fclient.example%2Fcallback",
    "/search?q=https://example.com",
  ])("preserves the local path %s", (value) => {
    expect(getSafeReturnPath(value)).toBe(value);
  });

  it.each([
    undefined,
    null,
    "",
    "https://example.com",
    "//example.com",
    "///example.com",
    "/\\example.com",
    "\\\\example.com",
    "javascript:void(0)",
    "data:text/html,example",
    " /settings",
    "/\t/example.com",
    "/\n/example.com",
    "/\r\n/example.com",
    "/\u0000/example.com",
    "/\u007f/example.com",
  ])("falls back to the homepage for %j", (value) => {
    expect(getSafeReturnPath(value)).toBe("/");
  });
});
