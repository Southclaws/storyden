import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { expect, it, vi } from "vitest";

import { Button } from "../button";

import { Menu } from "./Menu";
import { Item } from "./Menu.internal";

it("portals content and reports the selected item", async () => {
  const user = userEvent.setup();
  const onSelect = vi.fn();

  render(
    <Menu
      onSelect={onSelect}
      trigger={<Button variant="outline">Open menu</Button>}
    >
      <Item value="edit">Edit</Item>
    </Menu>,
  );

  await user.click(screen.getByRole("button", { name: "Open menu" }));
  const item = await screen.findByRole("menuitem", { name: "Edit" });
  await user.click(item);

  expect(onSelect).toHaveBeenCalledWith(
    expect.objectContaining({ value: "edit" }),
  );
});
