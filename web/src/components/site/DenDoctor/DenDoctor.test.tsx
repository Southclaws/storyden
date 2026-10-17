import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { DefaultSettings, Settings } from "@/lib/settings/settings";

import { DenDoctor } from "./DenDoctor";
import { FreshDenDoctor } from "./FreshDenDoctor";

vi.mock("@/config", () => ({
  API_ADDRESS: "https://localhost:8000",
  WEB_ADDRESS: window.location.origin,
}));

let status: Settings["onboarding_status"] | undefined;
vi.mock("@/lib/settings/settings-client", () => ({
  useSettings: () => ({
    settings: status ? { onboarding_status: status } : undefined,
  }),
}));

const healthy = {
  apiAddress: "https://api.example.com",
  webAddress: "https://example.com",
  browserOrigin: "https://example.com",
};

describe("Den Doctor", () => {
  beforeEach(() => {
    status = undefined;
  });

  it("stays quiet for a healthy fresh installation", () => {
    const { container } = render(<DenDoctor {...healthy} />);
    expect(container).toBeEmptyDOMElement();
  });

  it("offers a persistent configuration review in admin settings", async () => {
    render(<DenDoctor {...healthy} showHealthy />);
    expect(screen.getByText(/No address problems detected/)).toBeVisible();
    await userEvent.click(screen.getByText("Show configuration details"));
    expect(
      screen.getByText("PUBLIC_WEB_ADDRESS", { selector: "code" }),
    ).toBeVisible();
  });

  it("explains a problem and its remedy without requiring the API", () => {
    render(<DenDoctor {...healthy} apiAddress="http://api.example.com" />);
    expect(screen.getByRole("heading", { name: "Den Doctor" })).toBeVisible();
    expect(
      screen.getByText("Your website and API use different connection types"),
    ).toBeVisible();
    expect(screen.getByText(/Set the website and API addresses/)).toBeVisible();
    expect(screen.getByText("http://", { selector: "code" })).toBeVisible();
    expect(screen.getByText("https://", { selector: "code" })).toBeVisible();
    expect(
      screen.getByRole("link", { name: "Read troubleshooting guide" }),
    ).toHaveAttribute(
      "href",
      "https://www.storyden.org/docs/introduction/troubleshooting",
    );
  });

  it("offers the configured website when the browser uses another address", () => {
    render(<DenDoctor {...healthy} browserOrigin="https://www.example.com" />);

    expect(
      screen.getByRole("link", { name: "Open configured website" }),
    ).toHaveAttribute("href", "https://example.com");
  });

  it("renders paths, setting names and configured values as inline code", async () => {
    render(
      <DenDoctor {...healthy} apiAddress="https://api.example.com/backend" />,
    );

    expect(
      screen.getByText("/backend/api", { selector: "code" }),
    ).toBeVisible();

    await userEvent.click(screen.getByText("Show configuration details"));

    expect(
      screen.getByText("PUBLIC_API_ADDRESS", { selector: "code" }),
    ).toBeVisible();
    expect(
      screen.getByText("https://api.example.com/backend", { selector: "code" }),
    ).toBeVisible();
  });

  it("shows initial setup problems and removes them after account creation", async () => {
    const initial = {
      ...DefaultSettings,
      onboarding_status: "requires_first_account" as const,
    };
    const { rerender } = render(
      <FreshDenDoctor settings={initial} hasAddressIssues />,
    );
    expect(
      await screen.findByRole("heading", { name: "Den Doctor" }),
    ).toBeVisible();
    status = "requires_category";
    rerender(<FreshDenDoctor settings={initial} hasAddressIssues />);
    expect(screen.queryByRole("heading", { name: "Den Doctor" })).toBeNull();
  });

  it("does not show setup guidance to visitors of an existing community", () => {
    render(<FreshDenDoctor settings={DefaultSettings} hasAddressIssues />);
    expect(screen.queryByRole("heading", { name: "Den Doctor" })).toBeNull();
  });

  it("keeps a healthy fresh setup quiet when the browser uses the configured origin", () => {
    const settings = {
      ...DefaultSettings,
      onboarding_status: "requires_first_account" as const,
    };
    const { container } = render(
      <FreshDenDoctor settings={settings} hasAddressIssues={false} />,
    );

    expect(container).toBeEmptyDOMElement();
  });

  it("does not treat unavailable settings as a fresh installation", () => {
    render(<FreshDenDoctor hasAddressIssues />);

    expect(screen.queryByRole("heading", { name: "Den Doctor" })).toBeNull();
  });

  it("shows setup guidance only once a retry confirms the first account is needed", async () => {
    const { rerender } = render(<FreshDenDoctor hasAddressIssues />);

    expect(screen.queryByRole("heading", { name: "Den Doctor" })).toBeNull();

    status = "requires_first_account";
    rerender(<FreshDenDoctor hasAddressIssues />);

    expect(
      await screen.findByRole("heading", { name: "Den Doctor" }),
    ).toBeVisible();
  });
});
