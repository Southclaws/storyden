import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, test, vi } from "vitest";

import { CreatePageFromURLAction } from "./CreatePageFromURL";
import { importFromURLGenerator } from "./import";

vi.mock("@/lib/settings/capabilities", () => ({ useCapability: () => false }));
vi.mock("./import", async (original) => ({
  ...(await original<typeof import("./import")>()),
  importFromURLGenerator: vi.fn(),
}));

describe("CreatePageFromURLAction", () => {
  test("updates the progress label and completes the import", async () => {
    let finish!: () => void;
    const pending = new Promise<void>((resolve) => {
      finish = resolve;
    });
    vi.mocked(importFromURLGenerator).mockImplementation(async function* () {
      yield { step: "fetching_link" };
      await pending;
      yield { step: "complete" };
      return { step: "complete" };
    });
    const onComplete = vi.fn();
    render(<CreatePageFromURLAction onComplete={onComplete} />);
    fireEvent.click(screen.getByRole("button", { name: "Create from URL" }));
    const input = screen.getByPlaceholderText("Enter URL to import...");
    expect(
      await screen.findByRole("button", { name: "Import" }),
    ).toBeDisabled();
    fireEvent.change(input, {
      target: { value: "https://example.com/article" },
    });
    fireEvent.click(await screen.findByRole("button", { name: "Import" }));
    expect(await screen.findByText("Fetching page...")).toBeVisible();
    expect(onComplete).not.toHaveBeenCalled();
    finish();
    expect(await screen.findByText("Complete")).toBeVisible();
    await waitFor(() => expect(onComplete).toHaveBeenCalledOnce());
    expect(importFromURLGenerator).toHaveBeenCalledWith({
      url: "https://example.com/article",
      parentSlug: undefined,
      genaiAvailable: false,
    });
  });
});
