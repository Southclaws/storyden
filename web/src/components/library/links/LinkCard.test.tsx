import { render, screen, within } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import type { LinkReference } from "@/api/openapi-schema";

import { LinkCard } from "./LinkCard";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

const link: LinkReference = {
  id: "saved-reference",
  slug: "saved-reference",
  url: "https://reference.example/",
  domain: "reference.example",
  createdAt: "2026-01-01T10:00:00Z",
  updatedAt: "2026-01-01T10:00:00Z",
  title: "A useful reference",
  description: "Notes saved with the original page.",
};

describe("LinkCard", () => {
  it("keeps an older preview usable when newer status fields still say pending", () => {
    render(<LinkCard link={{ ...link, scrape_status: "pending" }} />);

    const card = screen.getByRole("article", { name: "A useful reference" });
    expect(
      within(card)
        .getByRole("link", { name: "A useful reference" })
        .getAttribute("href"),
    ).toBe("/links/saved-reference");
    expect(
      within(card).getByText("Notes saved with the original page."),
    ).toBeTruthy();
    expect(within(card).queryByText(/pending|fetching/i)).toBeNull();
    expect(within(card).queryByText(/unavailable|failed/i)).toBeNull();
  });

  it("does not imply active fetching while an empty preview is pending", () => {
    render(
      <LinkCard
        link={{
          ...link,
          title: undefined,
          description: undefined,
          scrape_status: "pending",
        }}
      />,
    );

    expect(screen.getByText("Preview pending")).toBeTruthy();
    expect(screen.queryByText(/fetching/i)).toBeNull();
    expect(
      screen.getByRole("link", { name: link.url }).getAttribute("href"),
    ).toBe("/links/saved-reference");
  });

  it("recognises a saved image even when the title and scrape timestamp are missing", () => {
    const { container } = render(
      <LinkCard
        link={{
          ...link,
          title: undefined,
          description: undefined,
          scrape_status: "pending",
          primary_image: {
            id: "original-preview",
            filename: "preview.png",
            path: "https://reference.example/preview.png",
            mime_type: "image/png",
            width: 1200,
            height: 630,
          },
        }}
      />,
    );

    expect(
      container.querySelector(
        'img[src="https://reference.example/preview.png"]',
      ),
    ).toBeTruthy();
    expect(screen.queryByText(/pending|fetching|unavailable/i)).toBeNull();
  });

  it("distinguishes a failed metadata refresh from an unreachable destination", () => {
    render(
      <LinkCard
        link={{ ...link, liveness_status: "alive", scrape_status: "failed" }}
      />,
    );

    expect(
      screen.getByText("Notes saved with the original page."),
    ).toBeTruthy();
    expect(screen.getByText("Preview refresh failed")).toBeTruthy();
    expect(screen.queryByText(/couldn't reach/i)).toBeNull();
    expect(
      screen
        .getByRole("link", { name: "reference.example" })
        .getAttribute("href"),
    ).toBe(link.url);
  });

  it("retains saved information when a reachability check fails", () => {
    render(
      <LinkCard
        link={{ ...link, liveness_status: "offline", scrape_status: "success" }}
      />,
    );

    expect(screen.getByText("Couldn't reach site")).toBeTruthy();
    expect(
      screen.getByText("Notes saved with the original page."),
    ).toBeTruthy();
    expect(screen.queryByText(/dead|deleted/i)).toBeNull();
  });
});
