import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { EmailVerificationScreen } from "./EmailVerificationScreen";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

describe("email verification return link", () => {
  it.each([
    ["https://example.com", "/"],
    ["//example.com", "/"],
    ["/\\example.com", "/"],
    ["/\t/example.com", "/"],
    ["/settings/account?tab=email", "/settings/account?tab=email"],
  ])("validates the return path %j", (returnURL, expected) => {
    render(<EmailVerificationScreen returnURL={returnURL} />);

    expect(
      screen.getByRole("link", { name: "Back to previous page" }),
    ).toHaveAttribute("href", expected);
  });

  it("omits the return link when no return path was supplied", () => {
    render(<EmailVerificationScreen />);

    expect(
      screen.queryByRole("link", { name: "Back to previous page" }),
    ).not.toBeInTheDocument();
  });
});
