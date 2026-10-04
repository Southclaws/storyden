import { notFound, redirect } from "next/navigation";
import { beforeEach, describe, expect, it, vi } from "vitest";

import Page from "./page";

vi.mock("next/navigation", () => ({
  notFound: vi.fn(),
  redirect: vi.fn(),
}));

vi.mock("@/config", () => ({ WEB_ADDRESS: "https://community.example" }));

beforeEach(() => {
  vi.mocked(notFound).mockImplementation(() => {
    throw new Error("Not found");
  });
});

describe("frontend resolver", () => {
  it("opens registration reviews without requiring a verification code", async () => {
    await Page({
      params: Promise.resolve({ kind: "admin", mark: "oauth-dcr-approval" }),
      searchParams: Promise.resolve({}),
    });

    expect(redirect).toHaveBeenCalledWith(
      "https://community.example/admin/authentication/agents?tab=reviews",
      "replace",
    );
  });

  it("preserves the approval code and repeated query parameters", async () => {
    await Page({
      params: Promise.resolve({ kind: "admin", mark: "oauth-dcr-approval" }),
      searchParams: Promise.resolve({
        verification_code: "ABCD-EFGH",
        tag: ["one", "two"],
      }),
    });

    expect(redirect).toHaveBeenCalledWith(
      "https://community.example/admin/authentication/agents?tab=reviews&verification_code=ABCD-EFGH&tag=one&tag=two",
      "replace",
    );
  });

  it("rejects unknown admin destinations", async () => {
    await expect(
      Page({
        params: Promise.resolve({ kind: "admin", mark: "unknown" }),
        searchParams: Promise.resolve({}),
      }),
    ).rejects.toThrow("Not found");
    expect(redirect).not.toHaveBeenCalled();
  });

  it.each([
    ["post", "/t/locate/example"],
    ["thread", "/t/example"],
    ["reply", "/t/locate/example"],
    ["node", "/l/example"],
    ["collection", "/c/example"],
    ["profile", "/m/example"],
  ])("continues resolving %s content", async (kind, path) => {
    await Page({
      params: Promise.resolve({ kind, mark: "example" }),
      searchParams: Promise.resolve({}),
    });

    expect(redirect).toHaveBeenCalledWith(
      `https://community.example${path}`,
      "replace",
    );
  });
});
