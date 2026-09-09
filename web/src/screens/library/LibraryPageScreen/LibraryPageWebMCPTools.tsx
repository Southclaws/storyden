"use client";

import { initializeWebMCPPolyfill } from "@mcp-b/webmcp-polyfill";
import { useMemo } from "react";
import { useWebMCP } from "usewebmcp";

import { WEBMCP_TOOL_DEFINITIONS } from "@/lib/webmcp/tools.generated";

import { useLibraryPageContext } from "./Context";
import { createLibraryPageLayoutToolImplementations } from "./libraryPageLayoutTools";
import { useEditState } from "./useEditState";

initializeWebMCPPolyfill();

export function LibraryPageWebMCPTools() {
  const { isDirectEditing } = useEditState();

  return isDirectEditing ? <RegisteredLibraryPageWebMCPTools /> : null;
}

function RegisteredLibraryPageWebMCPTools() {
  const { store } = useLibraryPageContext();
  const implementations = useMemo(
    () => createLibraryPageLayoutToolImplementations(store),
    [store],
  );

  useWebMCP({
    ...WEBMCP_TOOL_DEFINITIONS.library_page_layout_get,
    execute: implementations.library_page_layout_get,
  });
  useWebMCP({
    ...WEBMCP_TOOL_DEFINITIONS.library_page_block_add,
    execute: implementations.library_page_block_add,
  });
  useWebMCP({
    ...WEBMCP_TOOL_DEFINITIONS.library_page_block_remove,
    execute: implementations.library_page_block_remove,
  });
  useWebMCP({
    ...WEBMCP_TOOL_DEFINITIONS.library_page_block_move,
    execute: implementations.library_page_block_move,
  });
  useWebMCP({
    ...WEBMCP_TOOL_DEFINITIONS.library_page_block_assets_update,
    execute: implementations.library_page_block_assets_update,
  });
  useWebMCP({
    ...WEBMCP_TOOL_DEFINITIONS.library_page_block_directory_update,
    execute: implementations.library_page_block_directory_update,
  });

  return null;
}
