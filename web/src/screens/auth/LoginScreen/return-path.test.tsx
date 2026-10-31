import { act, renderHook } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { useLoginEmailForm } from "./LoginEmail/useLoginEmailForm";
import { useLoginHandleForm } from "./LoginHandle/useLoginHandleForm";

const mocks = vi.hoisted(() => ({
  push: vi.fn(),
  passwordSignin: vi.fn(),
  emailSignin: vi.fn(),
  mutate: vi.fn(),
  refreshFeed: vi.fn(),
}));

let params = new URLSearchParams();

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: mocks.push }),
  useSearchParams: () => params,
}));

vi.mock("@/api/client", () => ({
  handle: async (fn: () => Promise<void>) => fn(),
}));

vi.mock("@/api/openapi-client/accounts", () => ({
  useAccountGet: () => ({ mutate: mocks.mutate }),
}));

vi.mock("@/api/openapi-client/auth", () => ({
  authPasswordSignin: mocks.passwordSignin,
  authEmailPasswordSignin: mocks.emailSignin,
}));

vi.mock("@/lib/feed/refresh", () => ({
  refreshFeed: mocks.refreshFeed,
}));

describe("post-login return paths", () => {
  beforeEach(() => {
    params = new URLSearchParams();
    mocks.passwordSignin.mockResolvedValue({});
    mocks.emailSignin.mockResolvedValue({});
  });

  it.each([
    ["https://example.com", "/"],
    ["//example.com", "/"],
    ["/\\example.com", "/"],
    ["/\t/example.com", "/"],
    ["/t/welcome?sort=new#reply", "/t/welcome?sort=new#reply"],
  ])("validates %j after email login", async (returnPath, expected) => {
    params.set("return_url", returnPath);
    const { result } = renderHook(() => useLoginEmailForm());

    await act(async () => {
      result.current.form.setValue("identifier", "member@example.com");
      result.current.form.setValue("password", "test-password-123");
      await result.current.handlers.handleSubmit();
    });

    expect(mocks.emailSignin).toHaveBeenCalled();
    expect(mocks.push).toHaveBeenCalledWith(expected);
  });

  it.each([
    ["https://example.com", "/"],
    ["//example.com", "/"],
    ["/\\example.com", "/"],
    ["/\t/example.com", "/"],
    ["/t/welcome?sort=new#reply", "/t/welcome?sort=new#reply"],
  ])("validates %j after handle login", async (returnPath, expected) => {
    params.set("return_url", returnPath);
    const { result } = renderHook(() => useLoginHandleForm());

    act(() => {
      result.current.form.register("identifier", { value: "member" });
      result.current.form.register("token", { value: "test-password-123" });
    });
    await act(async () => {
      await result.current.form.handlePassword();
    });

    expect(mocks.passwordSignin).toHaveBeenCalled();
    expect(mocks.push).toHaveBeenCalledWith(expected);
  });
});
