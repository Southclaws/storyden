import { afterEach, describe, expect, it, vi } from "vitest";

import {
  executeWebMCPTool,
  getWebMCPClientToolContext,
  readWebMCPToolCallMetadata,
  webMCPToolScopeMatches,
} from "./webMCPClient";

interface ModelContextStub {
  getTools: ReturnType<typeof vi.fn>;
  executeTool?: ReturnType<typeof vi.fn>;
  addEventListener: ReturnType<typeof vi.fn>;
  removeEventListener: ReturnType<typeof vi.fn>;
}

const originalModelContext = Object.getOwnPropertyDescriptor(
  document,
  "modelContext",
);

afterEach(() => {
  if (originalModelContext) {
    Object.defineProperty(document, "modelContext", originalModelContext);
  } else {
    Reflect.deleteProperty(document, "modelContext");
  }
  vi.restoreAllMocks();
});

describe("getWebMCPClientToolContext", () => {
  it("snapshots browser tools using the durable client-tool wire shape", async () => {
    installModelContext({
      getTools: vi.fn().mockResolvedValue([
        {
          name: "library_page_block_add",
          title: "Add library block",
          description: "Add one block to the current library page.",
          inputSchema: JSON.stringify({
            type: "object",
            properties: { block: { type: "string" } },
            required: ["block"],
          }),
          annotations: {
            readOnlyHint: false,
            consequentialHint: true,
            unknownHint: true,
          },
        },
      ]),
    });

    await expect(getWebMCPClientToolContext("client-1")).resolves.toEqual({
      client_id: "client-1",
      tools: [
        {
          name: "library_page_block_add",
          title: "Add library block",
          description: "Add one block to the current library page.",
          inputSchema: {
            type: "object",
            properties: { block: { type: "string" } },
            required: ["block"],
          },
          annotations: {
            readOnlyHint: false,
            consequentialHint: true,
          },
        },
      ],
    });
  });

  it("omits tools whose browser-provided input schema is invalid", async () => {
    vi.spyOn(console, "warn").mockImplementation(() => undefined);
    installModelContext({
      getTools: vi.fn().mockResolvedValue([
        {
          name: "broken_tool",
          description: "Cannot be declared.",
          inputSchema: "not-json",
        },
      ]),
    });

    await expect(
      getWebMCPClientToolContext("client-1"),
    ).resolves.toBeUndefined();
  });
});

describe("executeWebMCPTool", () => {
  it("executes the currently registered browser tool and decodes JSON output", async () => {
    const tool = {
      name: "library_page_layout_get",
      description: "Read the current layout.",
      inputSchema: { type: "object" },
    };
    const executeTool = vi.fn().mockResolvedValue('{"blocks":[]}');
    installModelContext({
      getTools: vi.fn().mockResolvedValue([tool]),
      executeTool,
    });

    await expect(
      executeWebMCPTool("library_page_layout_get", {}),
    ).resolves.toEqual({ status: "completed", output: { blocks: [] } });
    expect(executeTool).toHaveBeenCalledWith(tool, "{}", {
      signal: undefined,
    });
  });

  it("unwraps structured content from the WebMCP result envelope", async () => {
    installModelContext({
      getTools: vi.fn().mockResolvedValue([
        {
          name: "library_page_layout_get",
          description: "Read the current layout.",
        },
      ]),
      executeTool: vi.fn().mockResolvedValue(
        JSON.stringify({
          content: [{ type: "text", text: '{"blocks":[]}' }],
          structuredContent: { blocks: [] },
          isError: false,
        }),
      ),
    });

    await expect(
      executeWebMCPTool("library_page_layout_get", {}),
    ).resolves.toEqual({ status: "completed", output: { blocks: [] } });
  });

  it("surfaces a WebMCP error envelope as an execution failure", async () => {
    installModelContext({
      getTools: vi.fn().mockResolvedValue([
        {
          name: "library_page_block_add",
          description: "Add a block.",
        },
      ]),
      executeTool: vi.fn().mockResolvedValue(
        JSON.stringify({
          content: [
            {
              type: "text",
              text: "The assets block is already present.",
            },
          ],
          isError: true,
        }),
      ),
    });

    await expect(
      executeWebMCPTool("library_page_block_add", { block: "assets" }),
    ).rejects.toThrow("The assets block is already present.");
  });

  it("keeps a disconnected call pending when its tool is not registered", async () => {
    installModelContext({ getTools: vi.fn().mockResolvedValue([]) });

    await expect(
      executeWebMCPTool("library_page_block_add", {}),
    ).resolves.toEqual({ status: "unavailable" });
  });
});

describe("WebMCP tool provenance", () => {
  it("accepts only WebMCP metadata with a client id", () => {
    expect(
      readWebMCPToolCallMetadata({
        storyden: {
          source: "webmcp",
          client_id: "client-1",
          scope: { page_type: "library" },
        },
      }),
    ).toEqual({
      clientID: "client-1",
      scope: { page_type: "library" },
    });
    expect(
      readWebMCPToolCallMetadata({
        storyden: { source: "builtin", client_id: "client-1" },
      }),
    ).toBeUndefined();
  });

  it("matches calls to their originating item before client execution", () => {
    expect(
      webMCPToolScopeMatches(
        { datagraph_item: { id: "library-1" } },
        { datagraph_item: { id: "library-1", kind: "node" } },
      ),
    ).toBe(true);
    expect(
      webMCPToolScopeMatches(
        { datagraph_item: { id: "library-1" } },
        { datagraph_item: { id: "library-2", kind: "node" } },
      ),
    ).toBe(false);
  });
});

function installModelContext(
  input: Pick<ModelContextStub, "getTools"> &
    Partial<Omit<ModelContextStub, "getTools">>,
) {
  Object.defineProperty(document, "modelContext", {
    configurable: true,
    value: {
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      ...input,
    },
  });
}
