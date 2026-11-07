"use client";

import { useCallback } from "react";
import { type WebMCPConfig, useWebMCP } from "usewebmcp";

import {
  WEBMCP_TOOL_DEFINITIONS,
  type WebMCPToolName,
  type WebMCPToolOutputMap,
} from "./tools.generated";
import { waitForWebMCPTool } from "./waitForTool";

type GeneratedWebMCPToolProps<TName extends WebMCPToolName> = {
  config: WebMCPConfig<
    (typeof WEBMCP_TOOL_DEFINITIONS)[TName]["inputSchema"],
    WebMCPToolOutputMap[TName]
  > & { name: TName };
};

export function GeneratedWebMCPTool<TName extends WebMCPToolName>({
  config,
}: GeneratedWebMCPToolProps<TName>) {
  useWebMCP({
    ...config,
    execute: async (input, options) => {
      try {
        return await config.execute(input, options);
      } catch (error) {
        if (options.signal.aborted) {
          throw error;
        }

        // WebMCP hides thrown callback errors; return useful edit failures as JSON.
        return {
          isError: true,
          content: [
            {
              type: "text",
              text: error instanceof Error ? error.message : String(error),
            },
          ],
        };
      }
    },
  });

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
