import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { NotificationsTrigger } from "./NotificationsTrigger";

describe("NotificationsTrigger", () => {
  it("announces unread notifications without exposing the visual indicator", () => {
    const { container } = render(<NotificationsTrigger hideLabel unread />);

    expect(
      screen.getByRole("button", { name: "Notifications, unread" }),
    ).toBeInTheDocument();
    expect(
      container.querySelector(".notifications-trigger__unread-indicator"),
    ).toHaveAttribute("aria-hidden", "true");
  });

  it("uses the standard label when there are no unread notifications", () => {
    render(<NotificationsTrigger hideLabel />);

    expect(
      screen.getByRole("button", { name: "Notifications" }),
    ).toBeInTheDocument();
  });
});
