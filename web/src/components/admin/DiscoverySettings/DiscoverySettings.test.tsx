import { act, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { SWRConfig } from "swr";
import { beforeEach, expect, it, vi } from "vitest";

import { generateRobotsTxt } from "@/lib/settings/robots-txt";
import { type AdminSettings, DefaultSettings } from "@/lib/settings/settings";

import { DiscoverySettingsForm } from "./DiscoverySettings";

vi.mock("@/lib/settings/robots-txt", () => ({
  generateRobotsTxt: vi.fn(),
}));

const mocks = vi.hoisted(() => ({ update: vi.fn() }));
vi.mock("@/lib/settings/mutation", () => ({
  useSettingsMutation: () => ({ updateSettings: mocks.update }),
}));
vi.mock("@/api/client", () => ({
  handle: async (fn: () => Promise<unknown>) => fn(),
}));

const generatedRobots =
  "User-agent: *\nAllow: /\n\nUser-agent: GPTBot\nDisallow: /\n";

beforeEach(() => {
  mocks.update.mockResolvedValue(undefined);
  vi.mocked(generateRobotsTxt).mockResolvedValue(generatedRobots);
  vi.spyOn(globalThis, "fetch").mockImplementation(
    async () => new Response(generatedRobots),
  );
});

function renderForm(settings: AdminSettings = DefaultSettings) {
  return render(
    <SWRConfig
      value={{
        provider: () => new Map(),
        dedupingInterval: 0,
        errorRetryCount: 0,
      }}
    >
      <DiscoverySettingsForm settings={settings} />
    </SWRConfig>,
  );
}

it("saves independent document controls and custom robots.txt", async () => {
  renderForm();
  await waitFor(() =>
    expect(screen.getByLabelText("robots.txt")).toHaveValue(generatedRobots),
  );
  await userEvent.click(screen.getByRole("checkbox", { name: "llms.txt" }));
  await userEvent.clear(screen.getByLabelText("robots.txt"));
  await userEvent.type(
    screen.getByLabelText("robots.txt"),
    "User-agent: *\nDisallow: /",
  );
  expect(
    screen.getByRole("checkbox", { name: "Allow AI training crawlers" }),
  ).toBeEnabled();
  await userEvent.click(screen.getByRole("button", { name: "Save" }));
  await waitFor(() =>
    expect(mocks.update).toHaveBeenCalledWith({
      services: {
        discovery: {
          search_indexing: true,
          training_crawling: false,
          llms_txt: true,
          auth_md: false,
          ai_catalog: false,
          api_catalog: false,
          mcp_server_card: false,
          robots_txt: "User-agent: *\nDisallow: /",
        },
      },
    }),
  );
});

it("restores the generated crawler policy without resetting document switches", async () => {
  renderForm({
    ...DefaultSettings,
    services: {
      discovery: {
        robots_txt: "User-agent: *\nDisallow: /",
        auth_md: false,
      },
    },
  });
  await userEvent.click(
    screen.getByRole("button", { name: "Use default robots.txt" }),
  );
  await waitFor(() =>
    expect(screen.getByLabelText("robots.txt")).toHaveValue(generatedRobots),
  );
  expect(generateRobotsTxt).toHaveBeenCalledWith(true, false);
  expect(mocks.update).not.toHaveBeenCalled();
  expect(
    screen.getByRole("checkbox", { name: "Allow AI training crawlers" }),
  ).toBeEnabled();
  await userEvent.click(screen.getByRole("button", { name: "Save" }));
  await waitFor(() =>
    expect(mocks.update).toHaveBeenCalledWith({
      services: {
        discovery: expect.objectContaining({ robots_txt: "", auth_md: false }),
      },
    }),
  );
  await waitFor(() =>
    expect(screen.getByLabelText("robots.txt")).toHaveValue(generatedRobots),
  );
});

it("shows the served document after saving presets without persisting it as an override", async () => {
  renderForm();
  const textarea = screen.getByLabelText("robots.txt");
  await waitFor(() => expect(textarea).toHaveValue(generatedRobots));
  await userEvent.click(
    screen.getByRole("checkbox", { name: "Enable search indexing" }),
  );

  const updatedRobots =
    "User-agent: *\nDisallow: /\n\nUser-agent: GPTBot\nDisallow: /\n";
  vi.mocked(fetch).mockImplementation(async () => new Response(updatedRobots));
  await userEvent.click(screen.getByRole("button", { name: "Save" }));

  await waitFor(() => expect(textarea).toHaveValue(updatedRobots));
  expect(mocks.update).toHaveBeenCalledWith({
    services: {
      discovery: expect.objectContaining({
        search_indexing: false,
        robots_txt: "",
      }),
    },
  });
  expect(fetch).toHaveBeenLastCalledWith("/robots.txt", { cache: "no-store" });
  expect(vi.mocked(fetch).mock.invocationCallOrder[1]).toBeGreaterThan(
    mocks.update.mock.invocationCallOrder[0]!,
  );

  await userEvent.click(
    screen.getByRole("checkbox", { name: "Allow AI training crawlers" }),
  );
  await userEvent.click(screen.getByRole("button", { name: "Save" }));
  await waitFor(() => expect(textarea).toBeEnabled());
  expect(mocks.update).toHaveBeenLastCalledWith({
    services: {
      discovery: expect.objectContaining({
        training_crawling: true,
        robots_txt: "",
      }),
    },
  });
});

it("keeps handwritten rules when presets change and the form is saved", async () => {
  renderForm();
  const textarea = screen.getByLabelText("robots.txt");
  await waitFor(() => expect(textarea).toHaveValue(generatedRobots));
  const customRobots = `${generatedRobots}\nDisallow: /private-path/\n`;
  await userEvent.type(textarea, "\nDisallow: /private-path/\n");
  await userEvent.click(
    screen.getByRole("checkbox", { name: "Allow AI training crawlers" }),
  );
  expect(textarea).toHaveValue(customRobots);
  expect(
    screen.getByText(/Preset changes won’t overwrite your edits/),
  ).toBeInTheDocument();

  vi.mocked(fetch).mockImplementation(async () => new Response(customRobots));
  await userEvent.click(screen.getByRole("button", { name: "Save" }));
  await waitFor(() => expect(textarea).toBeEnabled());
  expect(textarea).toHaveValue(customRobots);
  expect(mocks.update).toHaveBeenCalledWith({
    services: {
      discovery: expect.objectContaining({
        training_crawling: true,
        robots_txt: customRobots,
      }),
    },
  });
});

it("restores preset control when handwritten text is cleared and saved", async () => {
  renderForm({
    ...DefaultSettings,
    services: { discovery: { robots_txt: "User-agent: *\nDisallow: /" } },
  });
  const textarea = screen.getByLabelText("robots.txt");
  await userEvent.clear(textarea);
  expect(textarea).toHaveValue("");
  await userEvent.click(screen.getByRole("button", { name: "Save" }));

  await waitFor(() => expect(textarea).toHaveValue(generatedRobots));
  expect(mocks.update).toHaveBeenCalledWith({
    services: { discovery: expect.objectContaining({ robots_txt: "" }) },
  });
});

it("does not replace typing with a delayed initial document", async () => {
  const pending = Promise.withResolvers<Response>();
  vi.mocked(fetch).mockReturnValueOnce(pending.promise);
  renderForm();
  const textarea = screen.getByLabelText("robots.txt");
  await userEvent.type(textarea, "User-agent: *\nDisallow: /handwritten/");

  await act(async () => pending.resolve(new Response(generatedRobots)));
  expect(textarea).toHaveValue("User-agent: *\nDisallow: /handwritten/");
});

it("reports a failed document refresh without showing the error body as robots.txt", async () => {
  renderForm();
  const textarea = screen.getByLabelText("robots.txt");
  await waitFor(() => expect(textarea).toHaveValue(generatedRobots));
  vi.mocked(fetch).mockImplementation(
    async () => new Response("temporarily unavailable", { status: 503 }),
  );
  await userEvent.click(screen.getByRole("button", { name: "Save" }));

  await waitFor(() =>
    expect(screen.getByText(/Could not load robots.txt/)).toBeInTheDocument(),
  );
  expect(textarea).toHaveValue("");
  expect(textarea).toBeEnabled();
});

it("prevents edits while saving and refreshing the served document", async () => {
  renderForm();
  const textarea = screen.getByLabelText("robots.txt");
  await waitFor(() => expect(textarea).toHaveValue(generatedRobots));
  const pending = Promise.withResolvers<void>();
  mocks.update.mockReturnValueOnce(pending.promise);
  await userEvent.click(screen.getByRole("button", { name: "Save" }));

  expect(textarea).toBeDisabled();
  await userEvent.click(
    screen.getByRole("checkbox", { name: "Enable search indexing" }),
  );
  expect(
    screen.getByRole("checkbox", { name: "Enable search indexing" }),
  ).toBeChecked();
  await userEvent.click(
    screen.getByRole("checkbox", { name: "Allow AI training crawlers" }),
  );
  expect(
    screen.getByRole("checkbox", { name: "Allow AI training crawlers" }),
  ).not.toBeChecked();
  await act(async () => pending.resolve());
  await waitFor(() => expect(textarea).toBeEnabled());
});

it("keeps custom text visible while default rules are loading", async () => {
  const customRobots = "User-agent: *\nDisallow: /custom-only/\n";
  renderForm({
    ...DefaultSettings,
    services: { discovery: { robots_txt: customRobots } },
  });
  const pending = Promise.withResolvers<string>();
  vi.mocked(generateRobotsTxt).mockReturnValueOnce(pending.promise);

  await userEvent.click(
    screen.getByRole("button", { name: "Use default robots.txt" }),
  );
  const textarea = screen.getByLabelText("robots.txt");
  expect(textarea).toHaveValue(customRobots);
  expect(textarea).toBeDisabled();
  await act(async () => pending.resolve(generatedRobots));
  await waitFor(() => expect(textarea).toHaveValue(generatedRobots));
  expect(textarea).toBeEnabled();
});

it("keeps the override when default rules cannot be loaded and allows a retry", async () => {
  const customRobots = "User-agent: *\nDisallow: /custom-only/\n";
  renderForm({
    ...DefaultSettings,
    services: { discovery: { robots_txt: customRobots } },
  });
  vi.mocked(generateRobotsTxt).mockRejectedValueOnce(new Error("unavailable"));

  await userEvent.click(
    screen.getByRole("button", { name: "Use default robots.txt" }),
  );
  await waitFor(() =>
    expect(
      screen.getByText(/Could not load default robots.txt/),
    ).toBeInTheDocument(),
  );
  const textarea = screen.getByLabelText("robots.txt");
  expect(textarea).toHaveValue(customRobots);
  expect(textarea).toBeEnabled();

  await userEvent.click(
    screen.getByRole("button", { name: "Use default robots.txt" }),
  );
  await waitFor(() => expect(textarea).toHaveValue(generatedRobots));
  expect(
    screen.queryByText(/Could not load default robots.txt/),
  ).not.toBeInTheDocument();
});

it("previews the selected presets before saving", async () => {
  renderForm({
    ...DefaultSettings,
    services: { discovery: { robots_txt: "User-agent: *\nAllow: /\n" } },
  });
  await userEvent.click(
    screen.getByRole("checkbox", { name: "Enable search indexing" }),
  );
  await userEvent.click(
    screen.getByRole("checkbox", { name: "Allow AI training crawlers" }),
  );
  vi.mocked(generateRobotsTxt).mockResolvedValue(
    "User-agent: *\nDisallow: /\n",
  );
  await userEvent.click(
    screen.getByRole("button", { name: "Use default robots.txt" }),
  );

  await waitFor(() =>
    expect(screen.getByLabelText("robots.txt")).toHaveValue(
      "User-agent: *\nDisallow: /\n",
    ),
  );
  expect(generateRobotsTxt).toHaveBeenCalledWith(false, true);
  expect(mocks.update).not.toHaveBeenCalled();
});
