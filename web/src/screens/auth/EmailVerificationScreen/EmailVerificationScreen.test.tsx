import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { EmailVerificationScreen } from "./EmailVerificationScreen";

const mocks = vi.hoisted(() => ({
  push: vi.fn(),
  verify: vi.fn(),
}));

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: mocks.push }),
}));

vi.mock("@/api/client", () => ({
  handle: async (fn: () => Promise<void>) => fn(),
}));

vi.mock("@/api/openapi-client/auth", () => ({
  authEmailVerify: mocks.verify,
}));

describe("email verification return link", () => {
  beforeEach(() => {
    mocks.verify.mockResolvedValue({});
  });

  it("renders safely when the return parameter is repeated", () => {
    render(<EmailVerificationScreen returnURL={["/settings", "/profile"]} />);

    expect(
      screen.getByRole("link", { name: "Back to previous page" }),
    ).toHaveAttribute("href", "/");
  });

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

  it.each([
    ["https://example.com", "/"],
    ["/settings/account?tab=email", "/settings/account?tab=email"],
  ])(
    "navigates safely after verifying with %j",
    async (returnURL, expected) => {
      render(<EmailVerificationScreen returnURL={returnURL} />);

      const user = userEvent.setup();
      const firstDigit = await screen.findByLabelText("pin code 1 of 6");
      await user.type(
        screen.getByPlaceholderText("Email address..."),
        "member@example.com",
      );
      await user.click(firstDigit);
      await user.paste("123456");
      await waitFor(() => {
        expect(screen.getByRole("button", { name: "Verify" })).toBeEnabled();
      });
      await user.click(screen.getByRole("button", { name: "Verify" }));

      await waitFor(() => {
        expect(mocks.verify).toHaveBeenCalledWith({
          email: "member@example.com",
          code: "123456",
        });
        expect(mocks.push).toHaveBeenCalledWith(expected);
      });
    },
  );
});
