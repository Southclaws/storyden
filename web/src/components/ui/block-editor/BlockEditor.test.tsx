import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { createRef } from "react";
import { expect, it, vi } from "vitest";

import * as Menu from "@/components/ui/menu";

import * as BlockEditor from ".";

function renderMenuHandle({
  onAction = vi.fn(),
  onPointerDown = vi.fn(),
  ref = createRef<HTMLButtonElement>(),
} = {}) {
  render(
    <>
      <button type="button">Outside</button>
      <BlockEditor.Root>
        <BlockEditor.Gutter>
          <BlockEditor.MenuHandle
            ref={ref}
            aria-describedby="drag-instructions"
            onPointerDown={onPointerDown}
          >
            <Menu.ItemGroup>
              <Menu.ItemGroupLabel>Block settings</Menu.ItemGroupLabel>
              <Menu.Separator />
              <Menu.Item value="action" onClick={onAction}>
                Run action
              </Menu.Item>
            </Menu.ItemGroup>
          </BlockEditor.MenuHandle>
        </BlockEditor.Gutter>
        <BlockEditor.Content>Block content</BlockEditor.Content>
      </BlockEditor.Root>
    </>,
  );

  return { onAction, onPointerDown, ref };
}

it("keeps the menu open when non-interactive menu chrome is clicked", async () => {
  const user = userEvent.setup();
  renderMenuHandle();

  await user.click(
    screen.getByRole("button", { name: "Move or configure block" }),
  );
  await user.click(await screen.findByText("Block settings"));

  expect(screen.getByRole("menuitem", { name: "Run action" })).toBeVisible();
});

it("closes only after a genuine outside interaction", async () => {
  const user = userEvent.setup();
  renderMenuHandle();

  await user.click(
    screen.getByRole("button", { name: "Move or configure block" }),
  );
  expect(
    await screen.findByRole("menuitem", { name: "Run action" }),
  ).toBeVisible();

  await new Promise((resolve) => setTimeout(resolve, 0));
  await user.click(screen.getByRole("button", { name: "Outside" }));

  await waitFor(() =>
    expect(
      screen.queryByRole("menuitem", { name: "Run action" }),
    ).not.toBeInTheDocument(),
  );
});

it("runs a menu action once and closes the menu", async () => {
  const user = userEvent.setup();
  const onAction = vi.fn();
  renderMenuHandle({ onAction });

  await user.click(
    screen.getByRole("button", { name: "Move or configure block" }),
  );
  await user.click(await screen.findByRole("menuitem", { name: "Run action" }));

  expect(onAction).toHaveBeenCalledOnce();
  await waitFor(() =>
    expect(
      screen.queryByRole("menuitem", { name: "Run action" }),
    ).not.toBeInTheDocument(),
  );
});

it("forwards drag activator props and the button ref to the real menu trigger", async () => {
  const user = userEvent.setup();
  const onPointerDown = vi.fn();
  const ref = createRef<HTMLButtonElement>();
  renderMenuHandle({ onPointerDown, ref });

  const trigger = screen.getByRole("button", {
    name: "Move or configure block",
  });
  await user.pointer({ target: trigger, keys: "[MouseLeft]" });

  expect(onPointerDown).toHaveBeenCalledOnce();
  expect(ref.current).toBe(trigger);
  expect(trigger).toHaveAttribute("aria-describedby", "drag-instructions");
});
