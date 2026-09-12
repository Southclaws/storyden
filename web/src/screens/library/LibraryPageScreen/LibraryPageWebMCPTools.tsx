"use client";

import { initializeWebMCPPolyfill } from "@mcp-b/webmcp-polyfill";
import { useMemo } from "react";

import {
  GeneratedWebMCPTool,
  PageEditStartWebMCPTool,
} from "@/lib/webmcp/GeneratedWebMCPTool";
import { WEBMCP_TOOL_DEFINITIONS } from "@/lib/webmcp/tools.generated";

import { useLibraryPageContext } from "./Context";
import { createLibraryPageLayoutToolImplementations } from "./libraryPageLayoutTools";
import { useLibraryPagePermissions } from "./permissions";
import { useEditState } from "./useEditState";

initializeWebMCPPolyfill();

export function LibraryPageWebMCPTools() {
  const { store } = useLibraryPageContext();
  const { isDirectEditing, startDirectEdit } = useEditState();
  const { isAllowedToDirectEdit } = useLibraryPagePermissions();
  const implementations = useMemo(
    () => createLibraryPageLayoutToolImplementations(store),
    [store],
  );

  return (
    <>
      <GeneratedWebMCPTool
        config={{
          ...WEBMCP_TOOL_DEFINITIONS.library_page_layout_get,
          execute: implementations.library_page_layout_get,
        }}
      />
      {isAllowedToDirectEdit && (
        <PageEditStartWebMCPTool
          name="library_page_edit_start"
          expectedTool="library_page_block_directory_update"
          isEditing={isDirectEditing}
          readyMessage="Quick edit mode is active. Use the available Library page layout tools to make the requested change."
          unavailableMessage="Quick edit mode opened, but the Library page edit tools did not become available."
          startEditing={startDirectEdit}
        />
      )}
      {isDirectEditing && (
        <>
          <GeneratedWebMCPTool
            config={{
              ...WEBMCP_TOOL_DEFINITIONS.library_page_block_add,
              execute: implementations.library_page_block_add,
            }}
          />
          <GeneratedWebMCPTool
            config={{
              ...WEBMCP_TOOL_DEFINITIONS.library_page_block_remove,
              execute: implementations.library_page_block_remove,
            }}
          />
          <GeneratedWebMCPTool
            config={{
              ...WEBMCP_TOOL_DEFINITIONS.library_page_block_move,
              execute: implementations.library_page_block_move,
            }}
          />
          <GeneratedWebMCPTool
            config={{
              ...WEBMCP_TOOL_DEFINITIONS.library_page_block_assets_update,
              execute: implementations.library_page_block_assets_update,
            }}
          />
          <GeneratedWebMCPTool
            config={{
              ...WEBMCP_TOOL_DEFINITIONS.library_page_block_directory_update,
              execute: implementations.library_page_block_directory_update,
            }}
          />
        </>
      )}
    </>
  );
}
