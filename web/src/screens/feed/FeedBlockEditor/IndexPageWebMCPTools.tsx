"use client";

import { initializeWebMCPPolyfill } from "@mcp-b/webmcp-polyfill";
import { useMemo } from "react";

import {
  GeneratedWebMCPTool,
  PageEditStartWebMCPTool,
} from "@/lib/webmcp/GeneratedWebMCPTool";
import { WEBMCP_TOOL_DEFINITIONS } from "@/lib/webmcp/tools.generated";

import { useFeedBlockEditor } from "./Context";
import { createIndexPageLayoutToolImplementations } from "./indexPageLayoutTools";

initializeWebMCPPolyfill();

type Props = {
  isEditingEnabled: boolean;
  isEditing: boolean;
  startEditing: () => void;
};

export function IndexPageWebMCPTools({
  isEditingEnabled,
  isEditing,
  startEditing,
}: Props) {
  const editor = useFeedBlockEditor();
  const implementations = useMemo(
    () => createIndexPageLayoutToolImplementations(editor),
    [editor],
  );

  return (
    <>
      <GeneratedWebMCPTool
        config={{
          ...WEBMCP_TOOL_DEFINITIONS.index_page_layout_get,
          execute: implementations.index_page_layout_get,
        }}
      />
      {isEditingEnabled && (
        <PageEditStartWebMCPTool
          name="index_page_edit_start"
          expectedTool="index_page_block_library_update"
          isEditing={isEditing}
          readyMessage="Quick edit mode is active. Use the available index page layout tools to make the requested change."
          unavailableMessage="Quick edit mode opened, but the index page edit tools did not become available."
          startEditing={startEditing}
        />
      )}
      {isEditing && (
        <>
          <GeneratedWebMCPTool
            config={{
              ...WEBMCP_TOOL_DEFINITIONS.index_page_block_add,
              execute: implementations.index_page_block_add,
            }}
          />
          <GeneratedWebMCPTool
            config={{
              ...WEBMCP_TOOL_DEFINITIONS.index_page_block_remove,
              execute: implementations.index_page_block_remove,
            }}
          />
          <GeneratedWebMCPTool
            config={{
              ...WEBMCP_TOOL_DEFINITIONS.index_page_block_move,
              execute: implementations.index_page_block_move,
            }}
          />
          <GeneratedWebMCPTool
            config={{
              ...WEBMCP_TOOL_DEFINITIONS.index_page_block_categories_update,
              execute: implementations.index_page_block_categories_update,
            }}
          />
          <GeneratedWebMCPTool
            config={{
              ...WEBMCP_TOOL_DEFINITIONS.index_page_block_threads_update,
              execute: implementations.index_page_block_threads_update,
            }}
          />
          <GeneratedWebMCPTool
            config={{
              ...WEBMCP_TOOL_DEFINITIONS.index_page_block_quick_share_update,
              execute: implementations.index_page_block_quick_share_update,
            }}
          />
          <GeneratedWebMCPTool
            config={{
              ...WEBMCP_TOOL_DEFINITIONS.index_page_block_library_update,
              execute: implementations.index_page_block_library_update,
            }}
          />
        </>
      )}
    </>
  );
}
