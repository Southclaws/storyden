import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { expect, it } from "vitest";

import { Button } from "../button";

import { Tooltip } from "./Tooltip";

it("shows portalled content when its trigger is hovered", async () => {
  const user = userEvent.setup();

  render(
    <Tooltip content="Create a discussion" openDelay={0}>
      <Button>Create</Button>
    </Tooltip>,
  );

  await user.hover(screen.getByRole("button", { name: "Create" }));

  expect(
    await screen.findByRole("tooltip", { name: "Create a discussion" }),
  ).toBeVisible();
});
