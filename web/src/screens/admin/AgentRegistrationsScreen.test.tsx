import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { NuqsTestingAdapter } from "nuqs/adapters/testing";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { DefaultSettings } from "@/lib/settings/settings";

import { AgentRegistrationsScreen } from "./AgentRegistrationsScreen";

const mocks = vi.hoisted(() => ({
  settings: vi.fn(),
  back: vi.fn(),
  goBack: vi.fn(),
}));
let administrator = true;
let oauthEnabled = true;
let roleID: string | null = "bot-role";
let roles: { id: string; name: string }[] = [];

vi.mock("@/api/openapi-client/accounts", () => ({
  useAccountGet: () => ({
    data: {
      roles: [
        { permissions: [administrator ? "ADMINISTRATOR" : "MANAGE_SETTINGS"] },
      ],
    },
  }),
}));
vi.mock("@/api/openapi-client/admin", () => ({
  useAdminSettingsGet: (options: unknown) => {
    mocks.settings(options);
    return {
      data: {
        services: {
          oauth: {
            dynamic_registration_enabled: true,
            autonomous_registration_mode: "approval",
            autonomous_registration_role_id: roleID,
          },
        },
      },
    };
  },
}));
vi.mock("@/api/openapi-client/roles", () => ({
  useRoleList: () => ({ data: { roles } }),
}));
vi.mock("@/lib/settings/settings-client", () => ({
  useSettings: () => ({
    ready: true,
    settings: {
      ...DefaultSettings,
      capabilities: oauthEnabled ? ["oauth"] : [],
    },
  }),
}));
vi.mock("@/lib/navigation/smart-back", () => ({
  useSmartBack: (fallback: string) => {
    mocks.back(fallback);
    return mocks.goBack;
  },
}));
vi.mock(
  "@/components/admin/AgentRegistrations/RegistrationAccessTokens",
  () => ({
    RegistrationAccessTokens: ({ enabled }: { enabled: boolean }) => (
      <button disabled={!enabled}>Create token</button>
    ),
  }),
);
vi.mock(
  "@/components/admin/AgentRegistrations/AgentRegistrationApprovals",
  () => ({
    AgentRegistrationApprovals: ({
      enabled,
      defaultRoleName,
    }: {
      enabled: boolean;
      defaultRoleName: string;
    }) => (
      <section>
        <button disabled={!enabled}>Review registration</button>
        <p>{defaultRoleName}</p>
      </section>
    ),
  }),
);

describe("AgentRegistrationsScreen", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    administrator = true;
    oauthEnabled = true;
    roleID = "bot-role";
    roles = [{ id: "bot-role", name: "Bot" }];
  });
  it("uses the shared page header and returns to Authentication", async () => {
    renderScreen();
    expect(
      screen.getByRole("heading", { name: "Agent registrations", level: 1 }),
    ).toBeVisible();
    expect(mocks.back).toHaveBeenCalledWith("/admin/authentication");
    await userEvent.click(screen.getByRole("button", { name: "Back" }));
    expect(mocks.goBack).toHaveBeenCalledOnce();
    await userEvent.click(
      screen.getByRole("tab", { name: "Review registrations" }),
    );
    expect(screen.getByText("Bot + Member")).toBeVisible();
  });
  it("keeps a deleted default role safe to view", async () => {
    roles = [];
    renderScreen();
    await userEvent.click(
      screen.getByRole("tab", { name: "Review registrations" }),
    );
    expect(
      screen.getByText("Member only (saved role was deleted)"),
    ).toBeVisible();
  });
  it("disables admission actions when OAuth is unavailable", async () => {
    oauthEnabled = false;
    renderScreen();
    expect(screen.getByRole("button", { name: "Create token" })).toBeDisabled();
    await userEvent.click(
      screen.getByRole("tab", { name: "Review registrations" }),
    );
    expect(
      screen.getByRole("button", { name: "Review registration" }),
    ).toBeDisabled();
  });
  it("does not expose admission management to settings-only members", () => {
    administrator = false;
    renderScreen();
    expect(
      screen.getByText("Only administrators can manage agent registrations."),
    ).toBeVisible();
    expect(
      screen.queryByRole("button", { name: "Create token" }),
    ).not.toBeInTheDocument();
    expect(mocks.settings).toHaveBeenCalledWith({ swr: { enabled: false } });
  });
  it("starts with tokens and switches between separate views", async () => {
    renderScreen();
    expect(
      screen.getByRole("tab", { name: "Registration tokens" }),
    ).toHaveAttribute("aria-selected", "true");
    expect(screen.getByRole("button", { name: "Create token" })).toBeVisible();
    expect(
      screen.queryByRole("button", { name: "Review registration" }),
    ).not.toBeInTheDocument();
    await userEvent.click(
      screen.getByRole("tab", { name: "Review registrations" }),
    );
    expect(
      screen.getByRole("button", { name: "Review registration" }),
    ).toBeVisible();
    await waitFor(() =>
      expect(
        screen.queryByRole("button", { name: "Create token" }),
      ).not.toBeInTheDocument(),
    );
  });
  it("opens approval links directly in the reviews tab", () => {
    renderScreen("?verification_code=ABCD-EFGH");
    expect(
      screen.getByRole("tab", { name: "Review registrations" }),
    ).toHaveAttribute("aria-selected", "true");
    expect(
      screen.getByRole("button", { name: "Review registration" }),
    ).toBeVisible();
  });
});

function renderScreen(searchParams = "") {
  return render(
    <NuqsTestingAdapter searchParams={searchParams} hasMemory>
      <AgentRegistrationsScreen />
    </NuqsTestingAdapter>,
  );
}
