import {
  render as renderComponent,
  screen,
  waitFor,
  within,
} from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { NuqsTestingAdapter } from "nuqs/adapters/testing";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { OAuthDynamicRegistrationAccessToken } from "@/api/openapi-schema";

import { RegistrationAccessTokens } from "./RegistrationAccessTokens";

const mocks = vi.hoisted(() => ({
  create: vi.fn(),
  revoke: vi.fn(),
  refresh: vi.fn(),
  mutate: vi.fn(),
  list: vi.fn(),
}));
let tokens: OAuthDynamicRegistrationAccessToken[];
let listError: Error | undefined;

vi.mock("@/api/openapi-client/admin", () => ({
  adminOAuthDynamicRegistrationAccessTokenCreate: mocks.create,
  adminOAuthDynamicRegistrationAccessTokenRevoke: mocks.revoke,
  getAdminOAuthDynamicRegistrationAccessTokenListKey: () => [
    "/api/admin/oauth/dcr-iats",
  ],
  useAdminOAuthDynamicRegistrationAccessTokenList: (params: {
    page: string;
  }) => {
    mocks.list(params);
    const page = Number(params.page);
    return {
      data: {
        iats: tokens.slice((page - 1) * 50, page * 50),
        total_pages: Math.ceil(tokens.length / 50),
        page_size: 50,
      },
      error: listError,
      mutate: mocks.refresh,
    };
  },
}));
vi.mock("@/auth", () => ({ useSession: () => undefined }));
vi.mock("@/config", () => ({
  getAPIAddress: () => "https://community.example",
}));
vi.mock("swr", () => ({ useSWRConfig: () => ({ mutate: mocks.mutate }) }));
vi.mock("@/components/site/Modaldrawer/Modaldrawer", () => ({
  ModalDrawer: ({
    children,
    title,
  }: {
    children: React.ReactNode;
    title: string;
  }) => (
    <section role="dialog" aria-label={title}>
      {children}
    </section>
  ),
}));

const active: OAuthDynamicRegistrationAccessToken = {
  id: "active-token",
  label: "Engineering team",
  creator_account_id: "admin",
  createdAt: "2026-09-27T12:00:00Z",
  expires_at: "2030-09-28T12:00:00Z",
  max_registrations: 3,
  registration_count: 1,
};

describe("Registration access tokens", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    tokens = [];
    listError = undefined;
    mocks.create.mockResolvedValue({
      iat: active,
      token: "one-time-registration-secret",
    });
    mocks.revoke.mockResolvedValue(undefined);
    mocks.refresh.mockResolvedValue(undefined);
    mocks.mutate.mockResolvedValue(undefined);
  });

  it("pages through 100 tokens instead of rendering all of them", async () => {
    tokens = Array.from({ length: 100 }, (_, n) => ({
      ...active,
      id: `token-${n}`,
      label: `Token ${n + 1}`,
    }));
    render(<RegistrationAccessTokens enabled />);
    expect(screen.getAllByRole("listitem")).toHaveLength(50);
    expect(screen.getByText("Token 1")).toBeVisible();
    expect(screen.queryByText("Token 51")).not.toBeInTheDocument();
    await userEvent.click(screen.getAllByRole("link", { name: "2" })[0]!);
    await waitFor(() => expect(screen.getByText("Token 51")).toBeVisible());
    expect(screen.getAllByRole("listitem")).toHaveLength(50);
    expect(screen.queryByText("Token 1")).not.toBeInTheDocument();
    expect(mocks.list).toHaveBeenLastCalledWith({ page: "2" });
  });

  it("replaces expiry with revocation time in the same metadata slot", () => {
    tokens = [{ ...active, revoked_at: "2026-09-27T13:00:00Z" }];
    render(<RegistrationAccessTokens enabled />);
    const card = screen.getByRole("listitem");
    expect(within(card).getAllByRole("term")).toHaveLength(3);
    expect(within(card).queryByText("Expires")).not.toBeInTheDocument();
    expect(
      card.querySelector('time[datetime="2026-09-27T13:00:00Z"]'),
    ).toBeInTheDocument();
    expect(
      card.querySelector('time[datetime="2030-09-28T12:00:00Z"]'),
    ).not.toBeInTheDocument();
  });

  it("creates a limited token and only keeps its secret until the dialog closes", async () => {
    const user = userEvent.setup();
    const copy = vi.spyOn(navigator.clipboard, "writeText");
    const before = Date.now();
    render(<RegistrationAccessTokens enabled />);
    await userEvent.click(screen.getByRole("button", { name: "Create token" }));
    let dialog = screen.getByRole("dialog");
    await userEvent.type(
      within(dialog).getByRole("textbox", { name: "Label" }),
      "  Test agents  ",
    );
    const limit = within(dialog).getByRole("spinbutton", {
      name: "Maximum registrations",
    });
    await userEvent.clear(limit);
    await userEvent.type(limit, "3");
    await userEvent.click(
      within(dialog).getByRole("button", { name: "Create token" }),
    );
    await waitFor(() =>
      expect(
        screen.getByRole("textbox", { name: "Registration token" }),
      ).toHaveValue("one-time-registration-secret"),
    );
    const request = mocks.create.mock.calls[0]![0];
    expect(request).toMatchObject({
      label: "Test agents",
      max_registrations: 3,
    });
    expect(Date.parse(request.expires_at)).toBeGreaterThanOrEqual(
      before + 86400000,
    );
    expect(Date.parse(request.expires_at)).toBeLessThanOrEqual(
      Date.now() + 86400000,
    );
    const matchesTokenPage = mocks.mutate.mock.calls[0]![0];
    expect(matchesTokenPage(["/api/admin/oauth/dcr-iats", { page: "1" }])).toBe(
      true,
    );
    expect(matchesTokenPage(["/api/admin/oauth/dcr-iats", { page: "2" }])).toBe(
      true,
    );
    expect(matchesTokenPage(["/api/admin/oauth/clients"])).toBe(false);
    const prompt = `Register a bot account with \`sd\` on \`https://community.example\`. Choose a handle and a local context name.

Run \`sd auth register https://community.example --handle <handle> --name <context> --registration-token-stdin --format json\`, supplying this token on stdin: \`one-time-registration-secret\`

Then use \`sd --context <context> --help\` to learn how to use the Storyden command line tool.
`;
    expect(
      screen.getByRole("textbox", { name: "Prompt for your agent" }),
    ).toHaveValue(prompt);
    await user.click(screen.getByRole("button", { name: "Copy agent prompt" }));
    expect(copy).toHaveBeenCalledWith(prompt);
    await userEvent.click(screen.getByRole("button", { name: "Done" }));
    expect(
      screen.queryByRole("textbox", { name: "Prompt for your agent" }),
    ).not.toBeInTheDocument();
    expect(
      screen.queryByRole("textbox", { name: "Registration token" }),
    ).not.toBeInTheDocument();
    await userEvent.click(screen.getByRole("button", { name: "Create token" }));
    dialog = screen.getByRole("dialog");
    expect(within(dialog).getByRole("textbox", { name: "Label" })).toHaveValue(
      "",
    );
    expect(
      screen.queryByRole("textbox", { name: "Registration token" }),
    ).not.toBeInTheDocument();
    expect(mocks.create).toHaveBeenCalledOnce();
  });

  it("keeps the dialog open during creation so the secret is not lost", async () => {
    const issued = {
      iat: active,
      token: "one-time-registration-secret",
    };
    const creation = Promise.withResolvers<typeof issued>();
    mocks.create.mockReturnValueOnce(creation.promise);
    render(<RegistrationAccessTokens enabled />);
    await userEvent.click(screen.getByRole("button", { name: "Create token" }));
    const dialog = screen.getByRole("dialog");
    await userEvent.type(
      within(dialog).getByRole("textbox", { name: "Label" }),
      "Test agents",
    );
    await userEvent.click(
      within(dialog).getByRole("button", { name: "Create token" }),
    );
    await waitFor(() => expect(mocks.create).toHaveBeenCalledOnce());
    const cancel = within(dialog).getByRole("button", { name: "Cancel" });
    expect(cancel).toBeDisabled();
    await userEvent.click(cancel);
    expect(dialog).toBeInTheDocument();

    creation.resolve(issued);
    expect(
      await screen.findByRole("textbox", { name: "Registration token" }),
    ).toHaveValue("one-time-registration-secret");
    await userEvent.click(screen.getByRole("button", { name: "Done" }));
    expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
  });

  it("validates a blank label without minting a token", async () => {
    render(<RegistrationAccessTokens enabled />);
    await userEvent.click(screen.getByRole("button", { name: "Create token" }));
    await userEvent.click(
      within(screen.getByRole("dialog")).getByRole("button", {
        name: "Create token",
      }),
    );
    expect(await screen.findByText("Label is required")).toBeVisible();
    expect(mocks.create).not.toHaveBeenCalled();
  });

  it("keeps the form and reports a failed creation", async () => {
    mocks.create.mockRejectedValue(new Error("Token could not be created"));
    render(<RegistrationAccessTokens enabled />);
    await userEvent.click(screen.getByRole("button", { name: "Create token" }));
    const dialog = screen.getByRole("dialog");
    await userEvent.type(
      within(dialog).getByRole("textbox", { name: "Label" }),
      "Test agents",
    );
    await userEvent.click(
      within(dialog).getByRole("button", { name: "Create token" }),
    );
    expect(await screen.findByText("Token could not be created")).toBeVisible();
    expect(within(dialog).getByRole("textbox", { name: "Label" })).toHaveValue(
      "Test agents",
    );
    expect(
      screen.queryByRole("textbox", { name: "Registration token" }),
    ).not.toBeInTheDocument();
    expect(mocks.mutate).not.toHaveBeenCalled();
  });

  it("shows usage and distinct active, expired, exhausted, and revoked states", () => {
    tokens = [
      active,
      {
        ...active,
        id: "expired",
        label: "Expired token",
        expires_at: "2020-01-01T00:00:00Z",
      },
      {
        ...active,
        id: "exhausted",
        label: "Used token",
        registration_count: 3,
      },
      {
        ...active,
        id: "revoked",
        label: "Revoked token",
        revoked_at: "2026-09-27T13:00:00Z",
      },
    ];
    render(<RegistrationAccessTokens enabled />);
    expect(screen.getByText("Active")).toBeVisible();
    expect(screen.getByText("Expired")).toBeVisible();
    expect(screen.getByText("Exhausted")).toBeVisible();
    expect(screen.getAllByText("Revoked").length).toBeGreaterThan(0);
    expect(screen.getAllByRole("button", { name: "Revoke" })).toHaveLength(1);
    expect(screen.getByText("3 / 3")).toBeVisible();
  });

  it("requires confirmation before revoking and refreshes metadata after success", async () => {
    tokens = [active];
    render(<RegistrationAccessTokens enabled />);
    await userEvent.click(screen.getByRole("button", { name: "Revoke" }));
    expect(mocks.revoke).not.toHaveBeenCalled();
    await userEvent.click(screen.getByRole("button", { name: "Cancel" }));
    expect(mocks.revoke).not.toHaveBeenCalled();
    await userEvent.click(screen.getByRole("button", { name: "Revoke" }));
    await userEvent.click(
      screen.getByRole("button", { name: "Confirm revoke" }),
    );
    await waitFor(() =>
      expect(mocks.revoke).toHaveBeenCalledWith("active-token"),
    );
    expect(mocks.refresh).toHaveBeenCalledOnce();
  });

  it("does not mark a token revoked when the request fails", async () => {
    tokens = [active];
    mocks.revoke.mockRejectedValue(new Error("Revocation failed"));
    render(<RegistrationAccessTokens enabled />);
    await userEvent.click(screen.getByRole("button", { name: "Revoke" }));
    await userEvent.click(
      screen.getByRole("button", { name: "Confirm revoke" }),
    );
    await waitFor(() =>
      expect(screen.getByRole("button", { name: "Revoke" })).toBeEnabled(),
    );
    expect(screen.getByText("Active")).toBeVisible();
    expect(mocks.refresh).not.toHaveBeenCalled();
  });

  it("keeps revocation available while OAuth is disabled but disables minting", () => {
    tokens = [active];
    render(<RegistrationAccessTokens enabled={false} />);
    expect(screen.getByRole("button", { name: "Create token" })).toBeDisabled();
    expect(screen.getByRole("button", { name: "Revoke" })).toBeEnabled();
  });
});

function render(ui: React.ReactNode) {
  return renderComponent(
    <NuqsTestingAdapter hasMemory>{ui}</NuqsTestingAdapter>,
  );
}
