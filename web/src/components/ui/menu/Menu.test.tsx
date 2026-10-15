import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { renderToString } from "react-dom/server";
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

it("omits closed menu content from the server HTML by default", () => {
  const html = renderToString(
    <Menu trigger={<Button>Open menu</Button>}>
      <Item value="edit">Edit</Item>
    </Menu>,
  );
  const document = new DOMParser().parseFromString(html, "text/html");
  expect(document.querySelector('[role="menu"]')).toBeNull();
  expect(document.querySelector("button")?.textContent).toBe("Open menu");
});

it("allows content to be mounted eagerly when explicitly requested", () => {
  const html = renderToString(
    <Menu lazyMount={false} trigger={<Button>Open menu</Button>}>
      <Item value="edit">Edit</Item>
    </Menu>,
  );
  const document = new DOMParser().parseFromString(html, "text/html");
  expect(document.querySelector('[role="menuitem"]')?.textContent).toBe("Edit");
});

it.each(["{Enter}", " ", "{ArrowDown}"])(
  "opens on %s and retains content after Escape",
  async (key) => {
    const user = userEvent.setup();
    render(
      <Menu trigger={<Button>Open menu</Button>}>
        <Item value="edit">Edit</Item>
      </Menu>,
    );
    expect(screen.queryByText("Edit")).not.toBeInTheDocument();
    const trigger = screen.getByRole("button", { name: "Open menu" });
    await user.tab();
    expect(trigger).toHaveFocus();
    await user.keyboard(key);
    const item = await screen.findByRole("menuitem", { name: "Edit" });
    const menu = screen.getByRole("menu");
    await waitFor(() => expect(menu).toHaveFocus());
    expect(menu).toHaveAttribute("aria-activedescendant", item.id);
    expect(item).toHaveAttribute("data-highlighted");
    await user.keyboard("{Escape}");
    await waitFor(() => expect(trigger).toHaveFocus());
    expect(screen.getByText("Edit")).toBeInTheDocument();
    await waitFor(() =>
      expect(screen.queryByRole("menu")).not.toBeInTheDocument(),
    );
  },
);
