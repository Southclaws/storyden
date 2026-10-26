import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { expect, it, vi } from "vitest";

import * as Menu from "../menu";

import { ImageCropper } from ".";

vi.mock("@/config", () => ({ API_ADDRESS: "https://api.storyden.test" }));

it("keeps the default toolbar collapsible and reveals actions with the keyboard", async () => {
  render(
    <ImageCropper
      src="/photo.png"
      width={1920}
      height={480}
      onRemove={async () => {}}
    />,
  );
  const tools = screen.getByRole("button", { name: "Image tools" });
  expect(tools).toHaveAttribute("aria-expanded", "false");
  expect(
    screen.queryByRole("button", { name: "Remove image" }),
  ).not.toBeInTheDocument();
  tools.focus();
  await userEvent.setup().keyboard("{Enter}");
  expect(tools).toHaveAttribute("aria-expanded", "true");
  expect(screen.getByRole("button", { name: "Replace image" })).toBeVisible();
  expect(screen.getByRole("button", { name: "Remove image" })).toBeVisible();
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
        toolbar="expanded"
      >
        <Menu.Item value="caption">Edit caption</Menu.Item>
      </ImageCropper>,
    );
    await userEvent
      .setup()
      .click(screen.getByRole("button", { name: "Image options" }));
    expect(
      screen.getByRole("menuitem", { name: "Edit caption" }),
    ).toBeVisible();
    if (cropMode === "resizable")
      expect(
        screen.getByRole("menuitem", { name: "Crop ratio" }),
      ).toBeVisible();
    else
      expect(
        screen.queryByRole("menuitem", { name: "Crop ratio" }),
      ).not.toBeInTheDocument();
    expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
    expect(
      screen.queryByRole("button", { name: "Apply crop" }),
    ).not.toBeInTheDocument();
  },
);

it("announces saving and locks all expanded editing actions", () => {
  render(
    <ImageCropper
      width={1600}
      height={900}
      saving
      toolbar="expanded"
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
  expect(screen.getByRole("group", { name: "Crop area" })).toHaveAttribute(
    "aria-disabled",
    "true",
  );
  expect(screen.getByRole("button", { name: "Replace image" })).toBeDisabled();
  expect(screen.getByRole("button", { name: "Remove image" })).toBeDisabled();
  expect(screen.getByRole("button", { name: "Image options" })).toBeDisabled();
});
