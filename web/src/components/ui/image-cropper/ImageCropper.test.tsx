import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { expect, it, vi } from "vitest";

import * as Menu from "../menu";

import { ImageCropper } from ".";

vi.mock("@/config", () => ({ API_ADDRESS: "https://api.storyden.test" }));

it("opens and closes tools explicitly with pointer and keyboard input", async () => {
  render(
    <ImageCropper
      src="/photo.png"
      width={1920}
      height={480}
      onRemove={async () => {}}
    />,
  );
  const user = userEvent.setup();
  const tools = screen.getByRole("button", { name: "Image tools" });
  await user.hover(tools);
  expect(tools).toHaveAttribute("aria-expanded", "false");
  expect(
    screen.queryByRole("button", { name: "Remove image" }),
  ).not.toBeInTheDocument();
  tools.focus();
  await user.keyboard("{Enter}");
  expect(tools).toHaveAttribute("aria-expanded", "true");
  expect(screen.getByRole("button", { name: "Replace image" })).toBeVisible();
  expect(screen.getByRole("button", { name: "Remove image" })).toBeVisible();
  await user.keyboard("{Enter}");
  expect(tools).toHaveAttribute("aria-expanded", "false");
  expect(
    screen.queryByRole("button", { name: "Replace image" }),
  ).not.toBeInTheDocument();

  await user.click(tools);
  expect(tools).toHaveAttribute("aria-expanded", "true");
  await user.click(tools);
  expect(tools).toHaveAttribute("aria-expanded", "false");
  expect(
    screen.queryByRole("button", { name: "Remove image" }),
  ).not.toBeInTheDocument();
});

it.each(["fixed", "resizable"] as const)(
  "extends the %s cropper menu without adding controls around the image",
  async (cropMode) => {
    render(
      <ImageCropper
        src="/photo.png"
        width={1600}
        height={900}
        cropMode={cropMode}
      >
        <Menu.Item value="caption">Edit caption</Menu.Item>
      </ImageCropper>,
    );
    const user = userEvent.setup();
    expect(screen.getByRole("button", { name: "Image tools" })).toHaveAttribute(
      "aria-expanded",
      "false",
    );
    expect(
      screen.queryByRole("button", { name: "Image options" }),
    ).not.toBeInTheDocument();
    await user.click(screen.getByRole("button", { name: "Image tools" }));
    await user.click(screen.getByRole("button", { name: "Image options" }));
    expect(
      screen.getByRole("menuitem", { name: "Edit caption" }),
    ).toBeVisible();
    if (cropMode === "resizable") {
      expect(
        screen.getByRole("menuitem", { name: "Crop ratio" }),
      ).toBeVisible();
    } else {
      expect(
        screen.queryByRole("menuitem", { name: "Crop ratio" }),
      ).not.toBeInTheDocument();
    }
    expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
    expect(
      screen.queryByRole("button", { name: "Apply crop" }),
    ).not.toBeInTheDocument();
  },
);

it("announces saving and locks editing actions when the tools are opened", async () => {
  render(
    <ImageCropper
      width={1600}
      height={900}
      saving
      src="/photo.png"
      onRemove={async () => {}}
    >
      <Menu.Item value="caption">Edit caption</Menu.Item>
    </ImageCropper>,
  );
  expect(screen.getByRole("group", { name: "Crop area" })).toHaveAttribute(
    "aria-busy",
    "true",
  );
  await userEvent
    .setup()
    .click(screen.getByRole("button", { name: "Image tools" }));
  expect(screen.getByRole("group", { name: "Crop area" })).toHaveAttribute(
    "aria-disabled",
    "true",
  );
  expect(screen.getByRole("button", { name: "Replace image" })).toBeDisabled();
  expect(screen.getByRole("button", { name: "Remove image" })).toBeDisabled();
  expect(screen.getByRole("button", { name: "Image options" })).toBeDisabled();
});
