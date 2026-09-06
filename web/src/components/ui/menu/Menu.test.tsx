import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { expect, it, vi } from "vitest";

import { Button } from "../button";

import { Menu } from "./Menu";
import * as MenuPrimitive from "./Menu.internal";

const { Item } = MenuPrimitive;

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

it("reports selections from a nested menu", async () => {
  const user = userEvent.setup();
  const onSelect = vi.fn();

  render(
    <MenuPrimitive.Root>
      <MenuPrimitive.Trigger>Open menu</MenuPrimitive.Trigger>
      <MenuPrimitive.Positioner>
        <MenuPrimitive.Content>
          <MenuPrimitive.Root onSelect={onSelect}>
            <MenuPrimitive.TriggerItem>Layout</MenuPrimitive.TriggerItem>
            <MenuPrimitive.Positioner>
              <MenuPrimitive.Content>
                <MenuPrimitive.Item value="grid">Grid</MenuPrimitive.Item>
              </MenuPrimitive.Content>
            </MenuPrimitive.Positioner>
          </MenuPrimitive.Root>
        </MenuPrimitive.Content>
      </MenuPrimitive.Positioner>
    </MenuPrimitive.Root>,
  );

  await user.click(screen.getByRole("button", { name: "Open menu" }));
  await user.click(screen.getByRole("menuitem", { name: "Layout" }));
  await user.click(await screen.findByRole("menuitem", { name: "Grid" }));

  expect(onSelect).toHaveBeenCalledWith(
    expect.objectContaining({ value: "grid" }),
  );
});
