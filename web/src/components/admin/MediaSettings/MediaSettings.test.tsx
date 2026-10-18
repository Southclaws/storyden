import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { DefaultSettings } from "@/lib/settings/settings";

import { MediaSettingsForm } from "./MediaSettings";
import { type Props } from "./useMediaSettings";

const mocks = vi.hoisted(() => ({
  update: vi.fn(),
  revalidate: vi.fn(),
}));

vi.mock("@/lib/settings/mutation", () => ({
  useSettingsMutation: () => ({
    updateSettings: mocks.update,
    revalidate: mocks.revalidate,
  }),
}));

vi.mock("sonner", () => ({
  toast: { success: vi.fn() },
}));

function props(): Props {
  return {
    settings: {
      ...DefaultSettings,
      services: {
        moderation: { thread_body_length_max: 60_000 },
        media: {
          image_max_size: 33_554_432,
          image_max_pixels: 12_000_000,
          image_memory_budget: 268_435_456,
          image_estimated_memory: 402_653_184,
        },
      },
    },
  };
}

describe("MediaSettingsForm", () => {
  beforeEach(() => {
    mocks.update.mockResolvedValue(undefined);
    mocks.revalidate.mockResolvedValue(undefined);
  });

  it("displays the policy supplied by the API in MiB and megapixels", () => {
    const input = props();
    const { rerender } = render(<MediaSettingsForm {...input} />);

    expect(
      screen.getByRole("spinbutton", { name: "Maximum image size (MiB)" }),
    ).toHaveValue("32");
    expect(
      screen.getByRole("spinbutton", {
        name: "Maximum image resolution (megapixels)",
      }),
    ).toHaveValue("12");

    input.settings.services!.media!.image_max_size = 16_777_216;
    input.settings.services!.media!.image_max_pixels = 6_000_000;
    rerender(<MediaSettingsForm key="different-policy" {...input} />);

    expect(
      screen.getByRole("spinbutton", { name: "Maximum image size (MiB)" }),
    ).toHaveValue("16");
    expect(
      screen.getByRole("spinbutton", {
        name: "Maximum image resolution (megapixels)",
      }),
    ).toHaveValue("6");
  });

  it("saves only editable media policy using byte and pixel counts", async () => {
    const user = userEvent.setup();
    render(<MediaSettingsForm {...props()} />);

    const size = screen.getByRole("spinbutton", {
      name: "Maximum image size (MiB)",
    });
    const resolution = screen.getByRole("spinbutton", {
      name: "Maximum image resolution (megapixels)",
    });

    await user.clear(size);
    await user.type(size, "16.5");
    await user.clear(resolution);
    await user.type(resolution, "8.25");
    await user.click(screen.getByRole("button", { name: "Save" }));

    await waitFor(() => {
      expect(mocks.update).toHaveBeenCalledWith({
        services: {
          media: {
            image_max_size: 17_301_504,
            image_max_pixels: 8_250_000,
          },
        },
      });
    });
  });

  it("warns when the saved policy can exceed the worker budget", () => {
    const input = props();
    const { rerender } = render(<MediaSettingsForm {...input} />);

    expect(screen.getByRole("status")).toHaveTextContent(
      "Some images may exceed the processing budget",
    );
    expect(screen.getByRole("status")).toHaveTextContent("384 MiB");
    expect(screen.getByText(/Worker memory budget:/)).toHaveTextContent(
      "256 MiB per server process",
    );

    input.settings.services!.media = {
      ...input.settings.services!.media,
      image_estimated_memory: 268_435_456,
    };
    rerender(<MediaSettingsForm {...input} />);

    expect(screen.queryByRole("status")).not.toBeInTheDocument();
  });

  it("does not invent policy defaults when fields are absent", async () => {
    const input = props();
    input.settings.services = {};
    render(<MediaSettingsForm {...input} />);

    const size = screen.getByRole("spinbutton", {
      name: "Maximum image size (MiB)",
    });
    expect(size).toHaveValue("");
    expect(screen.queryByRole("status")).not.toBeInTheDocument();

    fireEvent.submit(size.closest("form")!);

    await waitFor(() => expect(size).toHaveAttribute("aria-invalid", "true"));
    expect(mocks.update).not.toHaveBeenCalled();
  });

  it("blocks submission when a required limit is cleared", async () => {
    const user = userEvent.setup();
    render(<MediaSettingsForm {...props()} />);

    const resolution = screen.getByRole("spinbutton", {
      name: "Maximum image resolution (megapixels)",
    });
    await user.clear(resolution);
    fireEvent.submit(resolution.closest("form")!);

    await waitFor(() => {
      expect(resolution).toHaveAttribute("aria-invalid", "true");
    });
    expect(mocks.update).not.toHaveBeenCalled();
  });

  it("keeps edited limits after a failed save and allows retry", async () => {
    const user = userEvent.setup();
    mocks.update.mockRejectedValueOnce(
      new Error("Settings could not be saved"),
    );
    render(<MediaSettingsForm {...props()} />);

    const size = screen.getByRole("spinbutton", {
      name: "Maximum image size (MiB)",
    });
    await user.clear(size);
    await user.type(size, "24");
    await user.click(screen.getByRole("button", { name: "Save" }));

    expect(
      await screen.findByText("Settings could not be saved"),
    ).toBeVisible();
    expect(size).toHaveValue("24");

    await user.click(screen.getByRole("button", { name: "Save" }));

    await waitFor(() => {
      expect(
        screen.queryByText("Settings could not be saved"),
      ).not.toBeInTheDocument();
    });
  });
});
