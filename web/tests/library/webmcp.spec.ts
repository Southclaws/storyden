import { Page, expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

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
const LIBRARY_TOOL_NAMES = [
  "library_page_layout_get",
  "library_page_block_add",
  "library_page_block_remove",
  "library_page_block_move",
  "library_page_block_assets_update",
  "library_page_block_directory_update",
] as const;

type RegisteredWebMCPTool = {
  name: string;
  description: string;
  inputSchema?: unknown;
  annotations?: Record<string, unknown>;
};

type LibraryPageFixture = {
  id: string;
  slug: string;
};

async function createLibraryPageFixture(): Promise<LibraryPageFixture> {
  const seed = randomUUID();
  let created: LibraryPageFixture | undefined;

  await withAdminAccessKey(async ({ nodeCreate }) => {
    const node = await nodeCreate({
      name: `WebMCP Library page ${seed}`,
      slug: `webmcp-library-page-${seed}`,
      content: `WebMCP Library page content ${seed}`,
      visibility: "published",
      meta: {
        layout: {
          blocks: [
            { type: "title" },
            {
              type: "directory",
              config: { layout: "table", columns: [] },
            },
            { type: "content" },
          ],
        },
      },
    });
    created = { id: node.id, slug: node.slug };
  });

  if (!created) {
    throw new Error("Library page fixture was not created");
  }
  return created;
}

async function loginAsFixtureAdmin(page: Page) {
  const handle = `wmcp_${randomUUID().replaceAll("-", "").slice(0, 20)}`;
  await createAdmin(page.context(), handle, PASSWORD);
  await login(page, handle, PASSWORD);
}

async function enterQuickEdit(page: Page, slug: string) {
  await page.goto(`/l/${slug}`);
  await page.getByRole("button", { name: "Edit", exact: true }).click();
  await page.getByRole("menuitem", { name: "Quick edit" }).click();
  await expect(
    page.getByRole("button", { name: "View", exact: true }),
  ).toBeVisible();
}

async function getLibraryTools(page: Page): Promise<RegisteredWebMCPTool[]> {
  return page.evaluate(async (prefix) => {
    const modelContext = (
      document as unknown as {
        modelContext?: { getTools(): Promise<RegisteredWebMCPTool[]> };
      }
    ).modelContext;
    if (!modelContext) return [];
    const tools = await modelContext.getTools();
    return tools.filter((tool) => tool.name.startsWith(prefix));
  }, "library_page_");
}

async function executeLibraryTool(
  page: Page,
  name: (typeof LIBRARY_TOOL_NAMES)[number],
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

async function expectLibraryTools(page: Page, expected: readonly string[]) {
  await expect
    .poll(async () =>
      (await getLibraryTools(page)).map((tool) => tool.name).sort(),
    )
    .toEqual([...expected].sort());
}

test.describe("Library page WebMCP", () => {
  test("mounts only during quick edit and safely describes blocks without configuration", async ({
    page,
  }) => {
    const libraryPage = await createLibraryPageFixture();
    await loginAsFixtureAdmin(page);
    await page.goto(`/l/${libraryPage.slug}`);

    await expectLibraryTools(page, []);

    await page.getByRole("button", { name: "Edit", exact: true }).click();
    await page.getByRole("menuitem", { name: "Quick edit" }).click();
    await expectLibraryTools(page, LIBRARY_TOOL_NAMES);

    const tools = await getLibraryTools(page);
    const layoutGet = tools.find(
      (tool) => tool.name === "library_page_layout_get",
    );
    expect(layoutGet).toMatchObject({
      description: "Retrieve the block layout of the current Library page.",
      annotations: { readOnlyHint: true },
    });

    const result = await executeLibraryTool(
      page,
      "library_page_layout_get",
      {},
    );
    expect(result).toEqual({
      message: "Retrieved 3 Library page layout blocks.",
      layout: {
        blocks: [
          { type: "title" },
          { type: "directory", layout: "table" },
          { type: "content" },
        ],
      },
    });

    await page.getByRole("button", { name: "View", exact: true }).click();
    await expectLibraryTools(page, []);
  });

  test("applies progressive edits, rejects stale instructions, and persists the resulting order", async ({
    page,
  }) => {
    const libraryPage = await createLibraryPageFixture();
    await loginAsFixtureAdmin(page);
    await enterQuickEdit(page, libraryPage.slug);
    await expectLibraryTools(page, LIBRARY_TOOL_NAMES);

    await expect(
      executeLibraryTool(page, "library_page_block_add", {
        block: "assets",
        after_block: "tags",
      }),
    ).rejects.toThrow(/reference block tags is not present/i);

    await executeLibraryTool(page, "library_page_block_add", {
      block: "assets",
      after_block: "directory",
    });
    await executeLibraryTool(page, "library_page_block_assets_update", {
      layout: "grid",
      grid_size: 3,
    });
    await executeLibraryTool(page, "library_page_block_directory_update", {
      layout: "grid",
    });
    await executeLibraryTool(page, "library_page_block_move", {
      block: "title",
    });

    const edited = await executeLibraryTool(
      page,
      "library_page_layout_get",
      {},
    );
    expect(edited["layout"]).toEqual({
      blocks: [
        { type: "directory", layout: "grid" },
        { type: "assets", layout: "grid", grid_size: 3 },
        { type: "content" },
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
      .toEqual(["directory", "assets", "content", "title"]);

    await expect(
      executeLibraryTool(page, "library_page_block_add", { block: "assets" }),
    ).rejects.toThrow(/assets block is already present/i);
    await expect(
      executeLibraryTool(page, "library_page_block_move", {
        block: "title",
        before_block: "title",
      }),
    ).rejects.toThrow(/cannot be moved before itself/i);

    await page.reload();
    await expectLibraryTools(page, LIBRARY_TOOL_NAMES);
    const persisted = await executeLibraryTool(
      page,
      "library_page_layout_get",
      {},
    );
    expect(persisted["layout"]).toEqual(edited["layout"]);
  });

  test("executes an ADK client tool through the Command Palette and returns its result to the Robot", async ({
    page,
  }) => {
    await setupRobotProviderWithScript(
      "mock/../robot/scripts/robot-chat-webmcp.yaml",
    );
    try {
      const libraryPage = await createLibraryPageFixture();
      await loginAsFixtureAdmin(page);
      await enterQuickEdit(page, libraryPage.slug);
      await expectLibraryTools(page, LIBRARY_TOOL_NAMES);

      let firstRobotRequest: Record<string, unknown> | undefined;
      page.on("request", (request) => {
        if (
          !firstRobotRequest &&
          request.method() === "POST" &&
          request.url().endsWith("/api/robots/sessions")
        ) {
          firstRobotRequest = request.postDataJSON() as Record<string, unknown>;
        }
      });

      await page.keyboard.press("ControlOrMeta+k");
      const palette = page.getByRole("dialog", { name: "Command Menu" });
      await expect(palette).toBeVisible();
      const commandInput = palette.getByRole("combobox", {
        name: "Command Menu",
      });
      await commandInput.fill("use the browser to add the assets block");
      await commandInput.press("Enter");

      await expect(
        palette.getByRole("group", {
          name: "Library Page Block Add tool call",
          exact: true,
        }),
      ).toBeVisible({ timeout: 15000 });
      await expect(
        palette.getByText("The browser added the assets block.", {
          exact: true,
        }),
      ).toBeVisible({ timeout: 15000 });

      const requestContext = firstRobotRequest?.["context"] as
        { datagraph_item?: { id?: string; slug?: string } } | undefined;
      const clientTools = firstRobotRequest?.["client_tools"] as
        { client_id?: string; tools?: { name: string }[] } | undefined;
      expect(requestContext?.datagraph_item).toMatchObject({
        id: libraryPage.id,
        slug: libraryPage.slug,
      });
      expect(clientTools?.client_id).toBeTruthy();
      expect(clientTools?.tools?.map((tool) => tool.name).sort()).toEqual(
        [...LIBRARY_TOOL_NAMES].sort(),
      );

      await expect(
        page.locator("#block-assets_content #gallery-strip"),
      ).toBeVisible();
      await expect(palette).toBeVisible();
      const persisted = await executeLibraryTool(
        page,
        "library_page_layout_get",
        {},
      );
      expect(persisted["layout"]).toEqual({
        blocks: [
          { type: "title" },
          { type: "directory", layout: "table" },
          { type: "assets" },
          { type: "content" },
        ],
      });
    } finally {
      await setupRobotProviderWithScript(DEFAULT_ROBOT_MODEL);
    }
  });
});
