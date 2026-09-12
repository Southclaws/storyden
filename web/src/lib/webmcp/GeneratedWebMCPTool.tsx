"use client";

import type { ToolInputSchema } from "@mcp-b/webmcp-polyfill/schema";
import type { JsonSchemaForInference } from "@mcp-b/webmcp-types";
import { useCallback } from "react";
import { type WebMCPConfig, useWebMCP } from "usewebmcp";

import {
  WEBMCP_TOOL_DEFINITIONS,
  type WebMCPToolName,
} from "./tools.generated";
import { waitForWebMCPTool } from "./waitForTool";

type GeneratedWebMCPToolProps<
  TInputSchema extends ToolInputSchema,
  TOutputSchema extends JsonSchemaForInference | undefined,
> = {
  config: WebMCPConfig<TInputSchema, TOutputSchema>;
};

export function GeneratedWebMCPTool<
  TInputSchema extends ToolInputSchema,
  TOutputSchema extends JsonSchemaForInference | undefined,
>({ config }: GeneratedWebMCPToolProps<TInputSchema, TOutputSchema>) {
  useWebMCP(config);

  return null;
}

type PageEditStartWebMCPToolProps = {
  name: "library_page_edit_start" | "index_page_edit_start";
  expectedTool: WebMCPToolName;
  isEditing: boolean;
  readyMessage: string;
  unavailableMessage: string;
  startEditing: () => void;
};

export function PageEditStartWebMCPTool({
  name,
  expectedTool,
  isEditing,
  readyMessage,
  unavailableMessage,
  startEditing,
}: PageEditStartWebMCPToolProps) {
  const execute = useCallback(async () => {
    if (isEditing) {
      return { message: "Quick edit mode is already active." };
    }

    startEditing();
    await waitForWebMCPTool(expectedTool, unavailableMessage);
    return { message: readyMessage };
  }, [expectedTool, isEditing, readyMessage, startEditing, unavailableMessage]);

  return (
    <GeneratedWebMCPTool
      config={{ ...WEBMCP_TOOL_DEFINITIONS[name], execute }}
    />
  );
}
