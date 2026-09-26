import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { Asset } from "@/api/openapi-schema";

import { FullPageChatInput } from "./FullPageChatInput";

const mocks = vi.hoisted(() => ({
  sendMessage: vi.fn(),
  upload: vi.fn(),
  toastError: vi.fn(),
}));

vi.mock("@/components/site/CommandPalette/RobotChat/RobotChatContext", () => ({
  useRobotChat: () => ({
    activeRobotName: "Denbot",
    sendMessage: mocks.sendMessage,
    cancelActiveTurn: vi.fn(),
    canCancelActiveTurn: false,
    isCancelling: false,
    status: "ready",
    queuedMessageCount: 0,
  }),
}));

vi.mock("@/components/content/useImageUpload", () => ({
  useImageUpload: () => ({ upload: mocks.upload }),
}));

vi.mock("sonner", () => ({
  toast: { error: mocks.toastError },
}));

vi.mock("@/components/asset/AssetThumbnail", () => ({
  AssetThumbnail: ({
    asset,
    deleteLabel,
    handleDelete,
  }: {
    asset: Asset;
    deleteLabel: string;
    handleDelete: () => void;
  }) => (
    <div>
      <img src={asset.path} alt={asset.filename} />
      <button type="button" aria-label={deleteLabel} onClick={handleDelete} />
    </div>
  ),
}));

const uploadedAsset: Asset = {
  id: "d7g0nnsrvimc3c4bg6sg",
  filename: "diagram.png",
  path: "/api/assets/d7g0nnsrvimc3c4bg6sg-diagram.png",
  mime_type: "image/png",
  width: 800,
  height: 800,
};

function imageFile(name = "diagram.png") {
  return new File(["image"], name, { type: "image/png" });
}

function selectFiles(files: File[]) {
  const input = document.querySelector<HTMLInputElement>('input[type="file"]');
  expect(input).not.toBeNull();
  fireEvent.change(input!, { target: { files } });
}

describe("FullPageChatInput image attachments", () => {
  beforeEach(() => {
    mocks.sendMessage.mockReset().mockResolvedValue(undefined);
    mocks.upload.mockReset().mockResolvedValue(uploadedAsset);
    mocks.toastError.mockReset();
    vi.spyOn(URL, "createObjectURL").mockReturnValue("blob:diagram-preview");
    vi.spyOn(URL, "revokeObjectURL").mockImplementation(() => undefined);
  });

  it("keeps the established subtle send action", () => {
    render(<FullPageChatInput />);

    expect(screen.getByRole("button", { name: "Send message" })).toHaveClass(
      "button--variant_subtle",
    );
  });

  it("uploads selected images and sends their asset IDs with an image-only message", async () => {
    const user = userEvent.setup();
    render(<FullPageChatInput />);

    selectFiles([imageFile()]);

    expect(screen.getByLabelText("Uploading diagram.png")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Send message" })).toBeDisabled();
    await screen.findByRole("img", { name: "diagram.png" });

    await user.click(screen.getByRole("button", { name: "Send message" }));

    expect(mocks.upload).toHaveBeenCalledWith(expect.any(File), {
      filename: "diagram.png",
    });
    expect(mocks.sendMessage).toHaveBeenCalledWith({
      text: "",
      assets: [uploadedAsset],
    });
    expect(screen.queryByLabelText("Attached images")).not.toBeInTheDocument();
  });

  it("accepts pasted and dropped images in display order", async () => {
    render(<FullPageChatInput />);
    const textarea = screen.getByRole("textbox", { name: "Message" });
    const pasted = imageFile("pasted.png");
    const dropped = imageFile("dropped.png");

    mocks.upload
      .mockResolvedValueOnce({
        ...uploadedAsset,
        id: "asset-pasted",
        filename: pasted.name,
      })
      .mockResolvedValueOnce({
        ...uploadedAsset,
        id: "asset-dropped",
        filename: dropped.name,
      });

    fireEvent.paste(textarea, {
      clipboardData: {
        items: [{ kind: "file", type: pasted.type, getAsFile: () => pasted }],
      },
    });
    fireEvent.drop(
      screen.getByRole("form", { name: "Send message to Robot" }),
      {
        dataTransfer: {
          types: ["Files"],
          files: [dropped],
          dropEffect: "none",
        },
      },
    );

    await screen.findByRole("img", { name: "pasted.png" });
    await screen.findByRole("img", { name: "dropped.png" });
    expect(
      screen.getAllByRole("img").map((image) => image.getAttribute("alt")),
    ).toEqual(["pasted.png", "dropped.png"]);
  });

  it("removes thumbnails and collapses the attachment shelf", async () => {
    const user = userEvent.setup();
    render(<FullPageChatInput />);

    selectFiles([imageFile()]);
    await screen.findByRole("img", { name: "diagram.png" });
    await user.click(
      screen.getByRole("button", { name: "Remove diagram.png" }),
    );

    expect(screen.queryByLabelText("Attached images")).not.toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Send message" })).toBeDisabled();
  });

  it("restores text and attachments when sending fails", async () => {
    const user = userEvent.setup();
    mocks.sendMessage.mockRejectedValueOnce(new Error("offline"));
    render(<FullPageChatInput />);

    selectFiles([imageFile()]);
    await screen.findByRole("img", { name: "diagram.png" });
    await user.type(
      screen.getByRole("textbox", { name: "Message" }),
      "Explain this",
    );
    await user.click(screen.getByRole("button", { name: "Send message" }));

    await waitFor(() => {
      expect(screen.getByRole("textbox", { name: "Message" })).toHaveValue(
        "Explain this",
      );
      expect(
        screen.getByRole("img", { name: "diagram.png" }),
      ).toBeInTheDocument();
    });
  });

  it("rejects unsupported image formats before upload", () => {
    render(<FullPageChatInput />);

    selectFiles([
      new File(["vector"], "diagram.svg", { type: "image/svg+xml" }),
    ]);

    expect(mocks.upload).not.toHaveBeenCalled();
    expect(mocks.toastError).toHaveBeenCalledWith(
      "Images must be GIF, JPEG, PNG, or WebP files.",
    );
  });
});
