"use client";

import { initializeWebMCPPolyfill } from "@mcp-b/webmcp-polyfill";
import { useCallback, useMemo } from "react";
import { useWebMCP } from "usewebmcp";

import {
  WEBMCP_TOOL_DEFINITIONS,
  WebMCPToolImplementation,
} from "@/lib/webmcp/tools.generated";

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
      <LibraryPageLayoutGetWebMCPTool
        execute={implementations.library_page_layout_get}
      />
      {isAllowedToDirectEdit && (
        <LibraryPageEditStartWebMCPTool
          isDirectEditing={isDirectEditing}
          startDirectEdit={startDirectEdit}
        />
      )}
      {isDirectEditing && (
        <RegisteredLibraryPageEditWebMCPTools
          implementations={implementations}
        />
      )}
    </>
  );
}

function LibraryPageLayoutGetWebMCPTool({
  execute,
}: {
  execute: WebMCPToolImplementation<"library_page_layout_get">;
}) {
  useWebMCP({
    ...WEBMCP_TOOL_DEFINITIONS.library_page_layout_get,
    execute,
  });

  return null;
}

function LibraryPageEditStartWebMCPTool({
  isDirectEditing,
  startDirectEdit,
}: {
  isDirectEditing: boolean;
  startDirectEdit: () => void;
}) {
  const execute = useCallback<
    WebMCPToolImplementation<"library_page_edit_start">
  >(async () => {
    if (isDirectEditing) {
      return { message: "Quick edit mode is already active." };
    }

    startDirectEdit();
    await waitForLibraryPageEditTools();
    return {
      message:
        "Quick edit mode is active. Use the available Library page layout tools to make the requested change.",
    };
  }, [isDirectEditing, startDirectEdit]);

  useWebMCP({
    ...WEBMCP_TOOL_DEFINITIONS.library_page_edit_start,
    execute,
  });

  return null;
}

async function waitForLibraryPageEditTools() {
  const modelContext = document.modelContext;
  if (!modelContext) {
    throw new Error("WebMCP is unavailable in this browser.");
  }
  const expectedTool = "library_page_block_directory_update";
  const isRegistered = async () =>
    (await modelContext.getTools()).some((tool) => tool.name === expectedTool);

  if (await isRegistered()) {
    return;
  }

  await new Promise<void>((resolve, reject) => {
    const timeout = window.setTimeout(() => {
      cleanup();
      reject(
        new Error(
          "Quick edit mode opened, but the Library page edit tools did not become available.",
        ),
      );
    }, 2_000);

    const handleToolChange = () => {
      void isRegistered()
        .then((ready) => {
          if (!ready) return;
          cleanup();
          resolve();
        })
        .catch((error) => {
          cleanup();
          reject(error);
        });
    };

    const cleanup = () => {
      window.clearTimeout(timeout);
      modelContext.removeEventListener("toolchange", handleToolChange);
    };

    modelContext.addEventListener("toolchange", handleToolChange);
    handleToolChange();
  });
}

type LibraryPageEditToolImplementations = ReturnType<
  typeof createLibraryPageLayoutToolImplementations
>;

function RegisteredLibraryPageEditWebMCPTools({
  implementations,
}: {
  implementations: LibraryPageEditToolImplementations;
}) {
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
