import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { expect, it } from "vitest";

import { Button } from "../button";

import { Popover } from "./Popover";
import { CloseTrigger, Description, Title } from "./Popover.internal";

it("opens portalled content and restores focus when closed", async () => {
  const user = userEvent.setup();

  render(
    <Popover trigger={<Button>Open settings</Button>}>
      <Title>Topic settings</Title>
      <Description>Configure notifications.</Description>
      <CloseTrigger asChild>
        <Button>Close</Button>
      </CloseTrigger>
    </Popover>,
  );

  const trigger = screen.getByRole("button", { name: "Open settings" });
  await user.click(trigger);
  expect(await screen.findByRole("dialog")).toBeVisible();

  await user.click(screen.getByRole("button", { name: /close/i }));
  await waitFor(() => expect(trigger).toHaveFocus());
});
