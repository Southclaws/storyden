import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import type { Link } from "@/api/openapi-schema";

import { LinkDetails } from "./LinkDetails";

const link: Link = {
  id: "saved-link",
  slug: "saved-link",
  url: "https://reference.example/",
  domain: "reference.example",
  title: "Saved reference",
  description: "The previously saved description.",
  createdAt: "2026-01-01T10:00:00Z",
  updatedAt: "2026-01-10T10:00:00Z",
  assets: [],
  nodes: [],
  posts: [],
  recomentations: [],
};

describe("LinkDetails", () => {
  it("keeps the successful scrape date separate from a later failed attempt", () => {
    render(
      <LinkDetails
        link={{
          ...link,
          liveness_status: "alive",
          liveness_checked_at: "2026-01-10T10:00:00Z",
          scrape_status: "failed",
          scraped_at: "2026-01-02T10:00:00Z",
          scrape_attempted_at: "2026-01-09T10:00:00Z",
        }}
      />,
    );

    const success = screen.getByText("Metadata retrieved").parentElement!;
    const attempt = screen.getByText("Latest metadata attempt").parentElement!;
    const check = screen.getByText("Last checked").parentElement!;

    expect(success.querySelector("time")?.dateTime).toBe(
      "2026-01-02T10:00:00.000Z",
    );
    expect(attempt.querySelector("time")?.dateTime).toBe(
      "2026-01-09T10:00:00.000Z",
    );
    expect(check.querySelector("time")?.dateTime).toBe(
      "2026-01-10T10:00:00.000Z",
    );
    expect(screen.getByText("The previously saved description.")).toBeTruthy();
    expect(
      screen.getByRole("link", { name: "Visit site" }).getAttribute("href"),
    ).toBe(link.url);
  });

  it("does not invent check or scrape dates for older saved metadata", () => {
    render(<LinkDetails link={{ ...link, scrape_status: "pending" }} />);

    const metadata = screen.getByText("Metadata retrieved").parentElement!;
    const check = screen.getByText("Last checked").parentElement!;

    expect(metadata.querySelector("time")).toBeNull();
    expect(check.querySelector("time")).toBeNull();
    expect(screen.queryByText(/pending|fetching/i)).toBeNull();
    expect(
      screen.getByRole("heading", { name: "Saved reference" }),
    ).toBeTruthy();
    expect(screen.getByText("The previously saved description.")).toBeTruthy();
  });
});
