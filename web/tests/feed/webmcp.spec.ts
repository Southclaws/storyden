import { Page, expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import { FeedConfig } from "../../src/lib/settings/feed";
import {
  createAdmin,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";
import {
  DEFAULT_ROBOT_MODEL,
  setupRobotProviderWithScript,
} from "../robot/helpers";

const PASSWORD = "TestPassword123!";
const INDEX_VIEW_TOOL_NAMES = [
  "index_page_edit_start",
  "index_page_layout_get",
] as const;
const INDEX_TOOL_NAMES = [
  ...INDEX_VIEW_TOOL_NAMES,
  "index_page_block_add",
  "index_page_block_remove",
  "index_page_block_move",
  "index_page_block_categories_update",
  "index_page_block_threads_update",
  "index_page_block_quick_share_update",
  "index_page_block_library_update",
] as const;

const INITIAL_FEED = {
  blocks: [
    { type: "title" },
    { type: "categories", layout: "list" },
    { type: "threads", source: "uncategorised" },
  ],
} satisfies FeedConfig;

type RegisteredWebMCPTool = {
  name: string;
  description: string;
  inputSchema?: unknown;
  annotations?: Record<string, unknown>;
};

async function setIndexFeed(feed: FeedConfig) {
  await withAdminAccessKey(async ({ adminSettingsUpdate }) => {
    await adminSettingsUpdate({ metadata: { feed } });
  });
}

async function loginAsFixtureAdmin(page: Page) {
  const handle = `wmcp_index_${randomUUID().replaceAll("-", "").slice(0, 16)}`;
  await createAdmin(page.context(), handle, PASSWORD);
  await login(page, handle, PASSWORD);
}

async function getIndexTools(page: Page): Promise<RegisteredWebMCPTool[]> {
  return page.evaluate(async (prefix) => {
    const modelContext = (
      document as unknown as {
        modelContext?: { getTools(): Promise<RegisteredWebMCPTool[]> };
      }
    ).modelContext;
    if (!modelContext) return [];
    const tools = await modelContext.getTools();
    return tools.filter((tool) => tool.name.startsWith(prefix));
  }, "index_page_");
}

async function executeIndexTool(
  page: Page,
  name: (typeof INDEX_TOOL_NAMES)[number],
  input: Record<string, unknown>,
): Promise<Record<string, unknown>> {
  return page.evaluate(
    async ({ name, input }) => {
      type Tool = { name: string };
      const modelContext = (
        document as unknown as {
          modelContext?: {
            getTools(): Promise<Tool[]>;
            executeTool?: (tool: Tool, input: string) => Promise<string | null>;
          };
        }
      ).modelContext;
      if (!modelContext?.executeTool) {
        throw new Error("WebMCP execution is unavailable");
      }
      const tool = (await modelContext.getTools()).find(
        (candidate) => candidate.name === name,
      );
      if (!tool) {
        throw new Error(`WebMCP tool ${name} is not registered`);
      }
      const result = await modelContext.executeTool(
        tool,
        JSON.stringify(input),
      );
      if (result === null) {
        throw new Error(`WebMCP tool ${name} returned no result`);
      }
      const parsed = JSON.parse(result) as Record<string, unknown>;
      if (parsed["isError"] === true) {
        const content = Array.isArray(parsed["content"])
          ? parsed["content"]
          : [];
        const message = content
          .map((item) =>
            item && typeof item === "object" && "text" in item
              ? item.text
              : undefined,
          )
          .find((text): text is string => typeof text === "string");
        throw new Error(message ?? `WebMCP tool ${name} failed`);
      }
      return (parsed["structuredContent"] ?? parsed) as Record<string, unknown>;
    },
    { name, input },
  );
}

async function expectIndexTools(page: Page, expected: readonly string[]) {
  await expect
    .poll(async () =>
      (await getIndexTools(page)).map((tool) => tool.name).sort(),
    )
    .toEqual([...expected].sort());
}

test.describe("Index page WebMCP", () => {
  test.beforeEach(async () => {
    await setIndexFeed(INITIAL_FEED);
  });

  test.afterEach(async () => {
    await setIndexFeed(INITIAL_FEED);
  });

  test("enters site edit mode and applies progressive changes to the mounted page", async ({
    page,
  }) => {
    await loginAsFixtureAdmin(page);
    await page.goto("/");
    await expectIndexTools(page, INDEX_VIEW_TOOL_NAMES);

    const viewTools = await getIndexTools(page);
    expect(
      viewTools.find((tool) => tool.name === "index_page_layout_get"),
    ).toMatchObject({
      description: "Retrieve the block layout of the current index page.",
      annotations: { readOnlyHint: true },
    });

    await executeIndexTool(page, "index_page_edit_start", {});
    await expect(page).toHaveURL(/editing=site/);
    await expect(
      page
        .getByRole("toolbar", { name: "Primary navigation" })
        .getByRole("button", { name: "Exit edit mode" }),
    ).toHaveAttribute("aria-pressed", "true");
    await expectIndexTools(page, INDEX_TOOL_NAMES);

    await executeIndexTool(page, "index_page_block_add", {
      block: "quick-share",
      after_block: "categories",
    });
    await executeIndexTool(page, "index_page_block_quick_share_update", {
      show_category_select: false,
    });
    await executeIndexTool(page, "index_page_block_categories_update", {
      layout: "grid",
    });
    const moved = await executeIndexTool(page, "index_page_block_move", {
      block: "title",
    });

    expect(moved["layout"]).toEqual({
      blocks: [
        { type: "categories", layout: "grid" },
        { type: "quick-share", show_category_select: false },
        { type: "threads", source: "uncategorised" },
        { type: "title" },
      ],
    });
    await expect
      .poll(() =>
        page
          .locator(".block-editor__root[data-block-type]")
          .evaluateAll((blocks) =>
            blocks.map((block) => block.getAttribute("data-block-type")),
          ),
      )
      .toEqual(["categories", "quick-share", "threads", "title"]);

    const quickShareBlock = page.locator(
      '.block-editor__root[data-block-type="quick-share"]',
    );
    await expect(quickShareBlock.locator(".quick-share")).toBeVisible();
    await expect(quickShareBlock.getByRole("combobox")).toHaveCount(0);

    await page.reload();
    await expectIndexTools(page, INDEX_TOOL_NAMES);
    const persisted = await executeIndexTool(page, "index_page_layout_get", {});
    expect(persisted["layout"]).toEqual(moved["layout"]);

    await page
      .getByRole("toolbar", { name: "Primary navigation" })
      .getByRole("button", { name: "Exit edit mode" })
      .click();
    await expectIndexTools(page, INDEX_VIEW_TOOL_NAMES);
  });

  test("lets Denbot enter edit mode and execute an index block tool", async ({
    page,
  }) => {
    await setupRobotProviderWithScript(
      "mock/../robot/scripts/robot-chat-index-webmcp.yaml",
    );
    try {
      await loginAsFixtureAdmin(page);
      await page.goto("/");
      await expectIndexTools(page, INDEX_VIEW_TOOL_NAMES);

      const robotRequests: Record<string, unknown>[] = [];
      page.on("request", (request) => {
        if (
          request.method() === "POST" &&
          request.url().endsWith("/api/robots/sessions")
        ) {
          robotRequests.push(request.postDataJSON() as Record<string, unknown>);
        }
      });

      await page.keyboard.press("ControlOrMeta+k");
      const palette = page.getByRole("dialog", { name: "Command Menu" });
      await expect(palette).toBeVisible();
      const commandInput = palette.getByRole("combobox", {
        name: "Command Menu",
      });
      await commandInput.fill("use the browser to add the quick-share block");
      await commandInput.press("Enter");

      await expect(
        palette.getByRole("group", {
          name: "Index Page Edit Start tool call",
          exact: true,
        }),
      ).toBeVisible({ timeout: 15000 });
      await expect(
        palette.getByRole("group", {
          name: "Index Page Block Add tool call",
          exact: true,
        }),
      ).toBeVisible({ timeout: 15000 });
      await expect(
        palette.getByText(
          "The browser added the quick-share block to the index page.",
          { exact: true },
        ),
      ).toBeVisible({ timeout: 15000 });

      const firstClientTools = robotRequests[0]?.["client_tools"] as
        { client_id?: string; tools?: { name: string }[] } | undefined;
      expect(firstClientTools?.client_id).toBeTruthy();
      expect(firstClientTools?.tools?.map((tool) => tool.name).sort()).toEqual(
        [...INDEX_VIEW_TOOL_NAMES].sort(),
      );
      await expect
        .poll(() =>
          robotRequests.some((request) => {
            const context = request["client_tools"] as
              { tools?: { name: string }[] } | undefined;
            return context?.tools?.some(
              (tool) => tool.name === "index_page_block_add",
            );
          }),
        )
        .toBe(true);

      await expect(
        page.locator(
          '.block-editor__root[data-block-type="quick-share"] .quick-share',
        ),
      ).toBeVisible();
      await expectIndexTools(page, INDEX_TOOL_NAMES);
      await expect(palette).toBeVisible();
    } finally {
      await setupRobotProviderWithScript(DEFAULT_ROBOT_MODEL);
    }
  });
});
