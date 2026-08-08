import { fireEvent, render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { LayoutEditModeButton } from "./LayoutEditModeButton";

const { mockUseSiteEditorState } = vi.hoisted(() => ({
  mockUseSiteEditorState: vi.fn(),
}));

vi.mock("@/lib/settings/site-editor-client", () => ({
  useSiteEditorState: mockUseSiteEditorState,
}));

describe("LayoutEditModeButton", () => {
  beforeEach(() => {
    mockUseSiteEditorState.mockReset();
  });

  it("renders the site-wide edit toggle when editing is allowed", () => {
    const handleToggleEditing = vi.fn();
    mockUseSiteEditorState.mockReturnValue({
      isEditingEnabled: true,
      isEditing: false,
      handleToggleEditing,
    });

    render(<LayoutEditModeButton />);

    const button = screen.getByRole("button", { name: "Edit site" });
    expect(button).toHaveAttribute("aria-pressed", "false");

    fireEvent.click(button);
    expect(handleToggleEditing).toHaveBeenCalledOnce();
  });

  it("announces the active edit state", () => {
    mockUseSiteEditorState.mockReturnValue({
      isEditingEnabled: true,
      isEditing: true,
      handleToggleEditing: vi.fn(),
    });

    render(<LayoutEditModeButton />);

    expect(
      screen.getByRole("button", { name: "Exit edit mode" }),
    ).toHaveAttribute("aria-pressed", "true");
  });

  it("renders no privileged action when editing is not allowed", () => {
    mockUseSiteEditorState.mockReturnValue({
      isEditingEnabled: false,
      isEditing: false,
      handleToggleEditing: vi.fn(),
    });

    const { container } = render(<LayoutEditModeButton />);

    expect(container).toBeEmptyDOMElement();
  });
});
