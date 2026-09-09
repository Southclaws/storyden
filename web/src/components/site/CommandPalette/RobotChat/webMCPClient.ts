import { initializeWebMCPPolyfill } from "@mcp-b/webmcp-polyfill";
import type {
  ChromeModelContext,
  RegisteredTool,
} from "@mcp-b/webmcp-types";

import {
  RobotChatContext,
  RobotClientTool,
  RobotClientToolAnnotations,
  RobotClientToolContext,
} from "@/api/openapi-schema";
import { generateXid } from "@/utils/xid";

initializeWebMCPPolyfill();

const CLIENT_ID_STORAGE_KEY = "storyden.webmcp.client-id";

export type WebMCPToolCallMetadata = {
  clientID: string;
  scope?: unknown;
};

export type WebMCPExecutionResult =
  { status: "unavailable" } | { status: "completed"; output: unknown };

let fallbackClientID: string | undefined;

export function getWebMCPClientID(): string {
  if (typeof window === "undefined") {
    return (fallbackClientID ??= generateXid());
  }

  try {
    const stored = window.sessionStorage.getItem(CLIENT_ID_STORAGE_KEY);
    if (stored) {
      return stored;
    }
    const generated = generateXid();
    window.sessionStorage.setItem(CLIENT_ID_STORAGE_KEY, generated);
    return generated;
  } catch {
    return (fallbackClientID ??= generateXid());
  }
}

export async function getWebMCPClientToolContext(
  clientID: string,
): Promise<RobotClientToolContext | undefined> {
  const modelContext = getBrowserModelContext();
  if (!modelContext) {
    return undefined;
  }

  let registered: RegisteredTool[];
  try {
    registered = await modelContext.getTools();
  } catch (error) {
    console.warn("[WebMCP] Could not inspect registered tools", error);
    return undefined;
  }
  const tools = registered.map(normalizeRegisteredTool).filter(isDefined);
  if (tools.length === 0) {
    return undefined;
  }

  return { client_id: clientID, tools };
}

export function getCurrentWebMCPClientToolContext() {
  return getWebMCPClientToolContext(getWebMCPClientID());
}

export async function executeWebMCPTool(
  name: string,
  input: unknown,
  signal?: AbortSignal,
): Promise<WebMCPExecutionResult> {
  const modelContext = getBrowserModelContext();
  if (!modelContext) {
    return { status: "unavailable" };
  }

  const tool = (await modelContext.getTools()).find(
    (candidate) => candidate.name === name,
  );
  if (!tool) {
    return { status: "unavailable" };
  }
  if (!modelContext.executeTool) {
    throw new Error("This browser does not support WebMCP tool execution.");
  }

  const encodedInput = JSON.stringify(input ?? {});
  const result = await modelContext.executeTool(tool, encodedInput, { signal });
  if (result === null) {
    throw new Error(
      `WebMCP tool "${name}" navigated without returning a result.`,
    );
  }

  return { status: "completed", output: parseToolResult(result) };
}

export function readWebMCPToolCallMetadata(
  providerMetadata: unknown,
): WebMCPToolCallMetadata | undefined {
  const storyden = readObject(providerMetadata)?.["storyden"];
  const metadata = readObject(storyden);
  if (metadata?.["source"] !== "webmcp") {
    return undefined;
  }
  const clientID = metadata["client_id"];
  if (typeof clientID !== "string" || !clientID) {
    return undefined;
  }
  return { clientID, scope: metadata["scope"] };
}

export function webMCPToolScopeMatches(
  expected: unknown,
  current: RobotChatContext,
): boolean {
  const expectedScope = readObject(expected);
  if (!expectedScope) {
    return true;
  }

  const expectedItem = readObject(expectedScope["datagraph_item"]);
  if (expectedItem) {
    return expectedItem["id"] === current.datagraph_item?.id;
  }

  const expectedPageType = expectedScope["page_type"];
  if (typeof expectedPageType === "string") {
    return expectedPageType === current.page_type;
  }

  return Object.keys(expectedScope).length === 0;
}

export function subscribeToWebMCPToolChanges(listener: () => void): () => void {
  const modelContext = getBrowserModelContext();
  if (!modelContext) {
    return () => undefined;
  }
  const eventListener = () => listener();
  modelContext.addEventListener("toolchange", eventListener);
  return () => modelContext.removeEventListener("toolchange", eventListener);
}

function getBrowserModelContext(): ChromeModelContext | undefined {
  if (typeof document === "undefined") {
    return undefined;
  }
  return document.modelContext;
}

function normalizeRegisteredTool(
  tool: RegisteredTool,
): RobotClientTool | undefined {
  const inputSchema = normalizeInputSchema(tool.inputSchema);
  if (!inputSchema) {
    console.warn(
      `[WebMCP] Ignoring tool with invalid input schema: ${tool.name}`,
    );
    return undefined;
  }

  const annotations = normalizeAnnotations(tool.annotations);
  return {
    name: tool.name,
    title: tool.title || undefined,
    description: tool.description,
    inputSchema,
    ...(annotations ? { annotations } : {}),
  };
}

function normalizeInputSchema(
  input: RegisteredTool["inputSchema"],
): Record<string, unknown> | undefined {
  if (input === undefined) {
    return { type: "object", properties: {} };
  }
  if (typeof input === "string") {
    try {
      return readObject(JSON.parse(input));
    } catch {
      return undefined;
    }
  }
  return readObject(input);
}

function normalizeAnnotations(
  input: RegisteredTool["annotations"],
): RobotClientToolAnnotations | undefined {
  const annotations: RobotClientToolAnnotations = {};
  for (const key of [
    "readOnlyHint",
    "untrustedContentHint",
    "consequentialHint",
  ] as const) {
    if (typeof input?.[key] === "boolean") {
      annotations[key] = input[key];
    }
  }
  return Object.keys(annotations).length > 0 ? annotations : undefined;
}

function parseToolResult(result: string): unknown {
  let parsed: unknown;
  try {
    parsed = JSON.parse(result);
  } catch {
    return result;
  }

  const envelope = readObject(parsed);
  if (!envelope) {
    return parsed;
  }
  if (envelope["isError"] === true) {
    const content = Array.isArray(envelope["content"])
      ? envelope["content"]
      : [];
    const message = content
      .map((item) => readObject(item)?.["text"])
      .find((text): text is string => typeof text === "string");
    throw new Error(message ?? "WebMCP tool execution failed.");
  }
  if ("structuredContent" in envelope) {
    return envelope["structuredContent"];
  }

  return parsed;
}

function readObject(value: unknown): Record<string, unknown> | undefined {
  if (!value || typeof value !== "object" || Array.isArray(value)) {
    return undefined;
  }
  return value as Record<string, unknown>;
}

function isDefined<T>(value: T | undefined): value is T {
  return value !== undefined;
}
