import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { expect, it, vi } from "vitest";

import { AssetUploadEditor } from "./AssetUploadEditor";

// Canvas rendering is unrelated to the editor's form and keyboard behavior.
vi.mock("react-avatar-editor", () => ({
  default: () => <canvas />,
}));

it.each([false, true])(
  "Enter submits the surrounding form through Create (image selected: %s)",
  async (selectImage) => {
    const user = userEvent.setup();
    const submit = vi.fn();
    const onUpload = vi.fn();
    const { container } = render(
      <form
        onSubmit={(event) => {
          event.preventDefault();
          submit((event.nativeEvent as SubmitEvent).submitter);
        }}
      >
        <AssetUploadEditor onUpload={onUpload} />
        <input aria-label="Name" />
        <button type="submit">Create</button>
      </form>,
    );

    const save = screen.getByRole("button", { name: "Save" });
    expect(save).toBeDisabled();

    if (selectImage) {
      const fileInput =
        container.querySelector<HTMLInputElement>('input[type="file"]')!;
      await user.upload(
        fileInput,
        new File(["image"], "cover.png", { type: "image/png" }),
      );
      expect(save).toBeEnabled();

      await user.click(save);
      expect(submit).not.toHaveBeenCalled();
    }

    await user.type(
      screen.getByRole("textbox", { name: "Name" }),
      "News{Enter}",
    );

    expect(submit).toHaveBeenCalledExactlyOnceWith(
      screen.getByRole("button", { name: "Create" }),
    );
    expect(onUpload).not.toHaveBeenCalled();
  },
);
