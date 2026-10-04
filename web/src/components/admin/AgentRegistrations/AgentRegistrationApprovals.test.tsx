import { render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { NuqsTestingAdapter } from "nuqs/adapters/testing";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { OAuthRegistrationApprovalReview } from "@/api/openapi-schema";

import { AgentRegistrationApprovals } from "./AgentRegistrationApprovals";

const mocks = vi.hoisted(() => ({
  submit: vi.fn(),
  bulkSubmit: vi.fn(),
  list: vi.fn(),
  mutate: vi.fn(),
  lookup: vi.fn(),
  lookupMutate: vi.fn(),
}));
let registrations: OAuthRegistrationApprovalReview[];
let lookupError: Error | undefined;
let lookupData: OAuthRegistrationApprovalReview | undefined;
vi.mock("@/api/openapi-client/admin", () => ({
  useAdminOAuthRegistrationApprovalList: (params: { page: string }) => {
    mocks.list(params);
    const page = Number(params.page);
    return {
      data: {
        registrations: registrations.slice((page - 1) * 50, page * 50),
        total_pages: Math.ceil(registrations.length / 50),
        page_size: 50,
        snapshot_at: "2026-09-27T12:00:00Z",
      },
      mutate: mocks.mutate,
    };
  },
  useAdminOAuthRegistrationApprovalGet: (params: unknown, options: unknown) => {
    mocks.lookup(params, options);
    return { data: lookupData, error: lookupError, mutate: mocks.lookupMutate };
  },
  adminOAuthRegistrationApprovalSubmit: mocks.submit,
  adminOAuthRegistrationApprovalBulkSubmit: mocks.bulkSubmit,
}));
vi.mock("@/auth", () => ({ useSession: () => undefined }));

const pending: OAuthRegistrationApprovalReview = {
  verification_code: "ABCD-EFGH",
  created_at: "2026-09-27T12:00:00Z",
  expires_at: "2030-09-27T12:10:00Z",
  metadata: { client_name: "library-bot", grant_types: ["client_credentials"] },
};

function renderSettings(searchParams = "") {
  return render(
    <NuqsTestingAdapter searchParams={searchParams} hasMemory>
      <AgentRegistrationApprovals enabled defaultRoleName="Bot + Member" />
    </NuqsTestingAdapter>,
  );
}

describe("OAuth registration approvals", () => {
  beforeEach(() => {
    registrations = [pending];
    lookupError = undefined;
    lookupData = undefined;
    mocks.submit.mockResolvedValue(undefined);
    mocks.bulkSubmit.mockResolvedValue({ updated: 100 });
    mocks.mutate.mockImplementation(async () => {
      registrations = [];
    });
  });

  it.each([
    [
      true,
      "Approve all",
      "Confirm approve all",
      "100 pending registrations approved.",
    ],
    [
      false,
      "Deny all",
      "Confirm deny all",
      "100 pending registrations denied.",
    ],
  ] as const)(
    "confirms bulk decision %s across all pages of 100 requests",
    async (approved, action, confirmation, outcome) => {
      registrations = Array.from({ length: 100 }, (_, n) => ({
        ...pending,
        verification_code: `code-${n}`,
        metadata: { client_name: `bot-${n}` },
      }));
      renderSettings();
      expect(screen.getAllByRole("listitem")).toHaveLength(50);
      await userEvent.click(screen.getAllByRole("link", { name: "2" })[0]!);
      await waitFor(() => expect(screen.getByText("@bot-50")).toBeVisible());
      expect(screen.queryByText("@bot-0")).not.toBeInTheDocument();
      await userEvent.click(screen.getByRole("button", { name: action }));
      expect(mocks.bulkSubmit).not.toHaveBeenCalled();
      expect(
        screen.getByText(/all pending requests across all pages/),
      ).toBeVisible();
      await userEvent.click(screen.getByRole("button", { name: confirmation }));
      expect(mocks.bulkSubmit).toHaveBeenCalledWith({
        approved,
        created_before: "2026-09-27T12:00:00Z",
      });
      await waitFor(() => expect(screen.getByText(outcome)).toBeVisible());
      expect(mocks.list).toHaveBeenLastCalledWith({ page: "1" });
      expect(screen.queryAllByRole("listitem")).toHaveLength(0);
    },
  );

  it("cancels bulk confirmation and retains the queue on a failed decision", async () => {
    renderSettings();
    await userEvent.click(screen.getByRole("button", { name: "Deny all" }));
    await userEvent.click(screen.getByRole("button", { name: "Cancel" }));
    expect(mocks.bulkSubmit).not.toHaveBeenCalled();
    mocks.bulkSubmit.mockRejectedValue(new Error("Decision failed"));
    await userEvent.click(screen.getByRole("button", { name: "Approve all" }));
    await userEvent.click(
      screen.getByRole("button", { name: "Confirm approve all" }),
    );
    await waitFor(() =>
      expect(screen.getByText("Decision failed")).toBeVisible(),
    );
    expect(screen.getByText("@library-bot")).toBeVisible();
    expect(mocks.mutate).not.toHaveBeenCalled();
  });

  it("opens an approval link using queued metadata without fetching it again", () => {
    renderSettings("?verification_code=ABCD-EFGH");
    expect(screen.getByText("@library-bot")).toBeVisible();
    expect(screen.getByText("Bot + Member")).toBeVisible();
    expect(mocks.lookup).toHaveBeenLastCalledWith(
      { verification_code: "ABCD-EFGH" },
      {
        swr: {
          enabled: false,
          shouldRetryOnError: false,
          keepPreviousData: false,
        },
      },
    );
    expect(mocks.submit).not.toHaveBeenCalled();
    expect(screen.getAllByText("@library-bot")).toHaveLength(1);
    expect(
      screen.queryByRole("button", { name: "Review" }),
    ).not.toBeInTheDocument();
  });

  it("submits the decision for the chosen card without selecting a review", async () => {
    registrations = [
      pending,
      {
        ...pending,
        verification_code: "JKLM-NPQR",
        metadata: { client_name: "second-bot" },
      },
    ];
    renderSettings();
    const card = screen.getByText("@second-bot").closest("li")!;
    expect(within(card).getByText("Pending")).toBeVisible();
    expect(within(card).getByText("Bot + Member")).toBeVisible();
    await userEvent.click(within(card).getByRole("button", { name: "Deny" }));
    expect(mocks.submit).toHaveBeenCalledWith({
      verification_code: "JKLM-NPQR",
      approved: false,
    });
  });

  it("includes a linked request outside the current queue page as a single card", async () => {
    registrations = [];
    lookupData = pending;
    renderSettings("?verification_code=ABCD-EFGH");
    expect(screen.getAllByText("@library-bot")).toHaveLength(1);
    await userEvent.click(screen.getByRole("button", { name: "Approve" }));
    await waitFor(() =>
      expect(screen.queryByText("@library-bot")).not.toBeInTheDocument(),
    );
    expect(mocks.submit).toHaveBeenCalledWith({
      verification_code: "ABCD-EFGH",
      approved: true,
    });
  });

  it.each([
    [true, "Approve", "Registration approved"],
    [false, "Deny", "Registration denied"],
  ] as const)(
    "submits the explicit decision %s and refreshes the queue",
    async (approved, action, outcome) => {
      renderSettings();
      await userEvent.click(screen.getByRole("button", { name: action }));
      await waitFor(() => expect(screen.getByText(outcome)).toBeVisible());
      expect(mocks.submit).toHaveBeenCalledWith({
        verification_code: "ABCD-EFGH",
        approved,
      });
      expect(mocks.mutate).toHaveBeenCalledOnce();
      expect(
        screen.queryByRole("button", { name: action }),
      ).not.toBeInTheDocument();
    },
  );

  it("retains the review and shows a failed decision instead of success", async () => {
    mocks.submit.mockRejectedValue(new Error("Could not save this decision"));
    renderSettings("?verification_code=ABCD-EFGH");
    await userEvent.click(screen.getByRole("button", { name: "Approve" }));
    await waitFor(() =>
      expect(screen.getByText("Decision could not be saved")).toBeVisible(),
    );
    expect(screen.getByRole("button", { name: "Approve" })).toBeEnabled();
    expect(screen.queryByText("Registration approved")).not.toBeInTheDocument();
    expect(mocks.mutate).not.toHaveBeenCalled();
  });

  it("does not offer a decision on stale metadata when a review becomes unavailable", () => {
    registrations = [];
    lookupData = pending;
    lookupError = new Error("This request has already been reviewed");
    renderSettings("?verification_code=ABCD-EFGH");
    expect(screen.getByText("Registration unavailable")).toBeVisible();
    expect(
      screen.queryByRole("button", { name: "Approve" }),
    ).not.toBeInTheDocument();
  });

  it("shows an unavailable request without offering approval", () => {
    registrations = [];
    lookupError = new Error("Expired");
    renderSettings("?verification_code=EXPD-CODE");
    expect(screen.getByText("Registration unavailable")).toBeVisible();
    expect(
      screen.queryByRole("button", { name: "Approve" }),
    ).not.toBeInTheDocument();
  });
});
