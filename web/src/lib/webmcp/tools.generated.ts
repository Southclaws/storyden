/* eslint-disable */
/**
 * This file was automatically generated from web/tools/schema.yaml.
 * DO NOT MODIFY IT BY HAND. Run: pnpm webmcp-schema
 */

export interface WebMCPTools {
  ToolLibraryPageEditStart?: LibraryPageEditStart;
  ToolLibraryPageLayoutGet?: LibraryPageLayoutGet;
  ToolLibraryPageBlockAdd?: LibraryPageBlockAdd;
  ToolLibraryPageBlockRemove?: LibraryPageBlockRemove;
  ToolLibraryPageBlockMove?: LibraryPageBlockMove;
  ToolLibraryPageBlockAssetsUpdate?: LibraryPageBlockAssetsUpdate;
  ToolLibraryPageBlockDirectoryUpdate?: LibraryPageBlockDirectoryUpdate;
  ToolIndexPageEditStart?: IndexPageEditStart;
  ToolIndexPageLayoutGet?: IndexPageLayoutGet;
  ToolIndexPageBlockAdd?: IndexPageBlockAdd;
  ToolIndexPageBlockRemove?: IndexPageBlockRemove;
  ToolIndexPageBlockMove?: IndexPageBlockMove;
  ToolIndexPageBlockCategoriesUpdate?: IndexPageBlockCategoriesUpdate;
  ToolIndexPageBlockThreadsUpdate?: IndexPageBlockThreadsUpdate;
  ToolIndexPageBlockQuickShareUpdate?: IndexPageBlockQuickShareUpdate;
  ToolIndexPageBlockLibraryUpdate?: IndexPageBlockLibraryUpdate;
  [k: string]: unknown;
}
/**
 * Enter quick edit mode for the current Library page so frontend-only layout tools can make the immediate change the user requested. Use this only when the user clearly asked to edit the live page. If it is unclear whether they want a quick edit or a version draft, ask them before calling this tool.
 */
export interface LibraryPageEditStart {
  input: ToolLibraryPageEditStartInput;
  output: ToolLibraryPageEditStartOutput;
  [k: string]: unknown;
}
export interface ToolLibraryPageEditStartInput {}
export interface ToolLibraryPageEditStartOutput {
  /**
   * Summary of the active edit mode.
   */
  message: string;
}
/**
 * Retrieve the block layout of the current Library page.
 */
export interface LibraryPageLayoutGet {
  input: ToolLibraryPageLayoutGetInput;
  output: ToolLibraryPageLayoutGetOutput;
  [k: string]: unknown;
}
export interface ToolLibraryPageLayoutGetInput {}
export interface ToolLibraryPageLayoutGetOutput {
  /**
   * Summary of the retrieved layout.
   */
  message: string;
  layout: LibraryPageLayout;
}
export interface LibraryPageLayout {
  /**
   * Blocks in their rendered order.
   */
  blocks: LibraryPageLayoutBlock[];
}
export interface LibraryPageLayoutBlock {
  type: "title" | "cover" | "link" | "content" | "assets" | "properties" | "directory" | "tags";
  /**
   * Layout mode for an assets or directory block.
   */
  layout?: "strip" | "grid" | "table";
  /**
   * Column count for an assets block using the grid layout.
   */
  grid_size?: number;
}
/**
 * Add a currently hidden block to the current Library page.
 */
export interface LibraryPageBlockAdd {
  input: ToolLibraryPageBlockAddInput;
  output: ToolLibraryPageBlockAddOutput;
  [k: string]: unknown;
}
export interface ToolLibraryPageBlockAddInput {
  block: "title" | "cover" | "link" | "content" | "assets" | "properties" | "directory" | "tags";
  /**
   * Existing block after which to insert the block; omit to append it.
   */
  after_block?: "title" | "cover" | "link" | "content" | "assets" | "properties" | "directory" | "tags";
}
export interface ToolLibraryPageBlockAddOutput {
  /**
   * Summary of the added block.
   */
  message: string;
  layout: LibraryPageLayout;
}
/**
 * Remove a visible block from the current Library page without deleting its content.
 */
export interface LibraryPageBlockRemove {
  input: ToolLibraryPageBlockRemoveInput;
  output: ToolLibraryPageBlockRemoveOutput;
  [k: string]: unknown;
}
export interface ToolLibraryPageBlockRemoveInput {
  block: "title" | "cover" | "link" | "content" | "assets" | "properties" | "directory" | "tags";
}
export interface ToolLibraryPageBlockRemoveOutput {
  /**
   * Summary of the removed block.
   */
  message: string;
  layout: LibraryPageLayout;
}
/**
 * Move a visible block within the current Library page layout.
 */
export interface LibraryPageBlockMove {
  input: ToolLibraryPageBlockMoveInput;
  output: ToolLibraryPageBlockMoveOutput;
  [k: string]: unknown;
}
export interface ToolLibraryPageBlockMoveInput {
  block: "title" | "cover" | "link" | "content" | "assets" | "properties" | "directory" | "tags";
  /**
   * Existing block before which to place the block; omit to move it to the end.
   */
  before_block?: "title" | "cover" | "link" | "content" | "assets" | "properties" | "directory" | "tags";
}
export interface ToolLibraryPageBlockMoveOutput {
  /**
   * Summary of the moved block.
   */
  message: string;
  layout: LibraryPageLayout;
}
/**
 * Update the presentation of the visible assets block on the current Library page.
 */
export interface LibraryPageBlockAssetsUpdate {
  input: ToolLibraryPageBlockAssetsUpdateInput;
  output: ToolLibraryPageBlockAssetsUpdateOutput;
  [k: string]: unknown;
}
export interface ToolLibraryPageBlockAssetsUpdateInput {
  layout: "strip" | "grid";
  /**
   * Number of columns used by the grid layout.
   */
  grid_size?: number;
}
export interface ToolLibraryPageBlockAssetsUpdateOutput {
  /**
   * Summary of the updated assets block.
   */
  message: string;
  layout: LibraryPageLayout;
}
/**
 * Update the presentation of the visible directory block on the current Library page.
 */
export interface LibraryPageBlockDirectoryUpdate {
  input: ToolLibraryPageBlockDirectoryUpdateInput;
  output: ToolLibraryPageBlockDirectoryUpdateOutput;
  [k: string]: unknown;
}
export interface ToolLibraryPageBlockDirectoryUpdateInput {
  layout: "table" | "grid";
}
export interface ToolLibraryPageBlockDirectoryUpdateOutput {
  /**
   * Summary of the updated directory block.
   */
  message: string;
  layout: LibraryPageLayout;
}
/**
 * Enter quick edit mode for the live index page so frontend-only block tools can make the immediate change the user requested. Use this only when the user clearly asked to edit the live page. If it is unclear whether they want a quick edit or a version draft, ask them before calling this tool.
 */
export interface IndexPageEditStart {
  input: ToolIndexPageEditStartInput;
  output: ToolIndexPageEditStartOutput;
  [k: string]: unknown;
}
export interface ToolIndexPageEditStartInput {}
export interface ToolIndexPageEditStartOutput {
  /**
   * Summary of the active edit mode.
   */
  message: string;
}
/**
 * Retrieve the block layout of the current index page.
 */
export interface IndexPageLayoutGet {
  input: ToolIndexPageLayoutGetInput;
  output: ToolIndexPageLayoutGetOutput;
  [k: string]: unknown;
}
export interface ToolIndexPageLayoutGetInput {}
export interface ToolIndexPageLayoutGetOutput {
  /**
   * Summary of the retrieved layout.
   */
  message: string;
  layout: IndexPageLayout;
}
export interface IndexPageLayout {
  /**
   * Blocks in their rendered order.
   */
  blocks: IndexPageLayoutBlock[];
}
export interface IndexPageLayoutBlock {
  type: "title" | "subtitle" | "content" | "cover" | "categories" | "threads" | "quick-share" | "library";
  /**
   * Layout mode for a categories or Library block.
   */
  layout?: "list" | "grid";
  /**
   * Thread source for a threads block.
   */
  source?: "all" | "uncategorised";
  /**
   * Whether a quick-share block shows its category picker.
   */
  show_category_select?: boolean;
  /**
   * Selected Library page ID for a library block; omitted for the Library root.
   */
  node_id?: string;
}
/**
 * Add a currently hidden block to the current index page.
 */
export interface IndexPageBlockAdd {
  input: ToolIndexPageBlockAddInput;
  output: ToolIndexPageBlockAddOutput;
  [k: string]: unknown;
}
export interface ToolIndexPageBlockAddInput {
  block: "title" | "subtitle" | "content" | "cover" | "categories" | "threads" | "quick-share" | "library";
  /**
   * Existing block after which to insert the block; omit to append it.
   */
  after_block?: "title" | "subtitle" | "content" | "cover" | "categories" | "threads" | "quick-share" | "library";
}
export interface ToolIndexPageBlockAddOutput {
  /**
   * Summary of the added block.
   */
  message: string;
  layout: IndexPageLayout;
}
/**
 * Remove a visible block from the current index page without deleting its content.
 */
export interface IndexPageBlockRemove {
  input: ToolIndexPageBlockRemoveInput;
  output: ToolIndexPageBlockRemoveOutput;
  [k: string]: unknown;
}
export interface ToolIndexPageBlockRemoveInput {
  block: "title" | "subtitle" | "content" | "cover" | "categories" | "threads" | "quick-share" | "library";
}
export interface ToolIndexPageBlockRemoveOutput {
  /**
   * Summary of the removed block.
   */
  message: string;
  layout: IndexPageLayout;
}
/**
 * Move a visible block within the current index page layout.
 */
export interface IndexPageBlockMove {
  input: ToolIndexPageBlockMoveInput;
  output: ToolIndexPageBlockMoveOutput;
  [k: string]: unknown;
}
export interface ToolIndexPageBlockMoveInput {
  block: "title" | "subtitle" | "content" | "cover" | "categories" | "threads" | "quick-share" | "library";
  /**
   * Existing block before which to place the block; omit to move it to the end.
   */
  before_block?: "title" | "subtitle" | "content" | "cover" | "categories" | "threads" | "quick-share" | "library";
}
export interface ToolIndexPageBlockMoveOutput {
  /**
   * Summary of the moved block.
   */
  message: string;
  layout: IndexPageLayout;
}
/**
 * Update the presentation of the visible categories block on the current index page.
 */
export interface IndexPageBlockCategoriesUpdate {
  input: ToolIndexPageBlockCategoriesUpdateInput;
  output: ToolIndexPageBlockCategoriesUpdateOutput;
  [k: string]: unknown;
}
export interface ToolIndexPageBlockCategoriesUpdateInput {
  layout: "list" | "grid";
}
export interface ToolIndexPageBlockCategoriesUpdateOutput {
  /**
   * Summary of the updated categories block.
   */
  message: string;
  layout: IndexPageLayout;
}
/**
 * Update the source of the visible threads block on the current index page.
 */
export interface IndexPageBlockThreadsUpdate {
  input: ToolIndexPageBlockThreadsUpdateInput;
  output: ToolIndexPageBlockThreadsUpdateOutput;
  [k: string]: unknown;
}
export interface ToolIndexPageBlockThreadsUpdateInput {
  source: "all" | "uncategorised";
}
export interface ToolIndexPageBlockThreadsUpdateOutput {
  /**
   * Summary of the updated threads block.
   */
  message: string;
  layout: IndexPageLayout;
}
/**
 * Configure the category picker on the visible quick-share block on the current index page.
 */
export interface IndexPageBlockQuickShareUpdate {
  input: ToolIndexPageBlockQuickShareUpdateInput;
  output: ToolIndexPageBlockQuickShareUpdateOutput;
  [k: string]: unknown;
}
export interface ToolIndexPageBlockQuickShareUpdateInput {
  show_category_select: boolean;
}
export interface ToolIndexPageBlockQuickShareUpdateOutput {
  /**
   * Summary of the updated quick-share block.
   */
  message: string;
  layout: IndexPageLayout;
}
/**
 * Configure the visible Library block on the current index page.
 */
export interface IndexPageBlockLibraryUpdate {
  input: ToolIndexPageBlockLibraryUpdateInput;
  output: ToolIndexPageBlockLibraryUpdateOutput;
  [k: string]: unknown;
}
export interface ToolIndexPageBlockLibraryUpdateInput {
  /**
   * Show the Library root or a specific Library page.
   */
  source: "library_root" | "library_page";
  /**
   * Required when source is library_page; omit when source is library_root.
   */
  node_id?: string;
  /**
   * Layout for the Library root; omit to preserve it. Specific pages control their own layout.
   */
  layout?: "list" | "grid";
}
export interface ToolIndexPageBlockLibraryUpdateOutput {
  /**
   * Summary of the updated Library block.
   */
  message: string;
  layout: IndexPageLayout;
}

export type WebMCPToolName = "library_page_edit_start" | "library_page_layout_get" | "library_page_block_add" | "library_page_block_remove" | "library_page_block_move" | "library_page_block_assets_update" | "library_page_block_directory_update" | "index_page_edit_start" | "index_page_layout_get" | "index_page_block_add" | "index_page_block_remove" | "index_page_block_move" | "index_page_block_categories_update" | "index_page_block_threads_update" | "index_page_block_quick_share_update" | "index_page_block_library_update";

export const WEBMCP_TOOL_NAMES = ["library_page_edit_start", "library_page_layout_get", "library_page_block_add", "library_page_block_remove", "library_page_block_move", "library_page_block_assets_update", "library_page_block_directory_update", "index_page_edit_start", "index_page_layout_get", "index_page_block_add", "index_page_block_remove", "index_page_block_move", "index_page_block_categories_update", "index_page_block_threads_update", "index_page_block_quick_share_update", "index_page_block_library_update"] as const;

export type WebMCPToolInputMap = {
  "library_page_edit_start": ToolLibraryPageEditStartInput;
  "library_page_layout_get": ToolLibraryPageLayoutGetInput;
  "library_page_block_add": ToolLibraryPageBlockAddInput;
  "library_page_block_remove": ToolLibraryPageBlockRemoveInput;
  "library_page_block_move": ToolLibraryPageBlockMoveInput;
  "library_page_block_assets_update": ToolLibraryPageBlockAssetsUpdateInput;
  "library_page_block_directory_update": ToolLibraryPageBlockDirectoryUpdateInput;
  "index_page_edit_start": ToolIndexPageEditStartInput;
  "index_page_layout_get": ToolIndexPageLayoutGetInput;
  "index_page_block_add": ToolIndexPageBlockAddInput;
  "index_page_block_remove": ToolIndexPageBlockRemoveInput;
  "index_page_block_move": ToolIndexPageBlockMoveInput;
  "index_page_block_categories_update": ToolIndexPageBlockCategoriesUpdateInput;
  "index_page_block_threads_update": ToolIndexPageBlockThreadsUpdateInput;
  "index_page_block_quick_share_update": ToolIndexPageBlockQuickShareUpdateInput;
  "index_page_block_library_update": ToolIndexPageBlockLibraryUpdateInput;
};

export type WebMCPToolOutputMap = {
  "library_page_edit_start": ToolLibraryPageEditStartOutput;
  "library_page_layout_get": ToolLibraryPageLayoutGetOutput;
  "library_page_block_add": ToolLibraryPageBlockAddOutput;
  "library_page_block_remove": ToolLibraryPageBlockRemoveOutput;
  "library_page_block_move": ToolLibraryPageBlockMoveOutput;
  "library_page_block_assets_update": ToolLibraryPageBlockAssetsUpdateOutput;
  "library_page_block_directory_update": ToolLibraryPageBlockDirectoryUpdateOutput;
  "index_page_edit_start": ToolIndexPageEditStartOutput;
  "index_page_layout_get": ToolIndexPageLayoutGetOutput;
  "index_page_block_add": ToolIndexPageBlockAddOutput;
  "index_page_block_remove": ToolIndexPageBlockRemoveOutput;
  "index_page_block_move": ToolIndexPageBlockMoveOutput;
  "index_page_block_categories_update": ToolIndexPageBlockCategoriesUpdateOutput;
  "index_page_block_threads_update": ToolIndexPageBlockThreadsUpdateOutput;
  "index_page_block_quick_share_update": ToolIndexPageBlockQuickShareUpdateOutput;
  "index_page_block_library_update": ToolIndexPageBlockLibraryUpdateOutput;
};

export type WebMCPToolImplementation<K extends WebMCPToolName> = (
  input: WebMCPToolInputMap[K],
) => WebMCPToolOutputMap[K] | Promise<WebMCPToolOutputMap[K]>;

export type WebMCPToolImplementations = {
  [K in WebMCPToolName]: WebMCPToolImplementation<K>;
};

export const WEBMCP_TOOL_DEFINITIONS = {
  "library_page_edit_start": {
    "name": "library_page_edit_start",
    "title": "Start Library Page Quick Edit",
    "description": "Enter quick edit mode for the current Library page so frontend-only layout tools can make the immediate change the user requested. Use this only when the user clearly asked to edit the live page. If it is unclear whether they want a quick edit or a version draft, ask them before calling this tool.",
    "inputSchema": {
      "type": "object",
      "properties": {},
      "additionalProperties": false
    },
    "outputSchema": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "message": {
          "type": "string",
          "description": "Summary of the active edit mode."
        }
      },
      "additionalProperties": false
    },
    "annotations": {
      "readOnlyHint": false,
      "destructiveHint": false,
      "idempotentHint": true,
      "openWorldHint": false,
      "untrustedContentHint": false
    },
    "pathPrefix": "/l/"
  },
  "library_page_layout_get": {
    "name": "library_page_layout_get",
    "title": "Get Library Page Layout",
    "description": "Retrieve the block layout of the current Library page.",
    "inputSchema": {
      "type": "object",
      "properties": {},
      "additionalProperties": false
    },
    "outputSchema": {
      "type": "object",
      "required": [
        "message",
        "layout"
      ],
      "properties": {
        "message": {
          "type": "string",
          "description": "Summary of the retrieved layout."
        },
        "layout": {
          "type": "object",
          "required": [
            "blocks"
          ],
          "properties": {
            "blocks": {
              "type": "array",
              "description": "Blocks in their rendered order.",
              "items": {
                "type": "object",
                "required": [
                  "type"
                ],
                "properties": {
                  "type": {
                    "type": "string",
                    "enum": [
                      "title",
                      "cover",
                      "link",
                      "content",
                      "assets",
                      "properties",
                      "directory",
                      "tags"
                    ]
                  },
                  "layout": {
                    "type": "string",
                    "enum": [
                      "strip",
                      "grid",
                      "table"
                    ],
                    "description": "Layout mode for an assets or directory block."
                  },
                  "grid_size": {
                    "type": "integer",
                    "minimum": 1,
                    "maximum": 4,
                    "description": "Column count for an assets block using the grid layout."
                  }
                },
                "additionalProperties": false
              }
            }
          },
          "additionalProperties": false
        }
      },
      "additionalProperties": false
    },
    "annotations": {
      "readOnlyHint": true,
      "destructiveHint": false,
      "idempotentHint": true,
      "openWorldHint": false,
      "untrustedContentHint": false
    },
    "pathPrefix": "/l/"
  },
  "library_page_block_add": {
    "name": "library_page_block_add",
    "title": "Add Library Page Block",
    "description": "Add a currently hidden block to the current Library page.",
    "inputSchema": {
      "type": "object",
      "required": [
        "block"
      ],
      "properties": {
        "block": {
          "type": "string",
          "enum": [
            "title",
            "cover",
            "link",
            "content",
            "assets",
            "properties",
            "directory",
            "tags"
          ]
        },
        "after_block": {
          "description": "Existing block after which to insert the block; omit to append it.",
          "type": "string",
          "enum": [
            "title",
            "cover",
            "link",
            "content",
            "assets",
            "properties",
            "directory",
            "tags"
          ]
        }
      },
      "additionalProperties": false
    },
    "outputSchema": {
      "type": "object",
      "required": [
        "message",
        "layout"
      ],
      "properties": {
        "message": {
          "type": "string",
          "description": "Summary of the added block."
        },
        "layout": {
          "type": "object",
          "required": [
            "blocks"
          ],
          "properties": {
            "blocks": {
              "type": "array",
              "description": "Blocks in their rendered order.",
              "items": {
                "type": "object",
                "required": [
                  "type"
                ],
                "properties": {
                  "type": {
                    "type": "string",
                    "enum": [
                      "title",
                      "cover",
                      "link",
                      "content",
                      "assets",
                      "properties",
                      "directory",
                      "tags"
                    ]
                  },
                  "layout": {
                    "type": "string",
                    "enum": [
                      "strip",
                      "grid",
                      "table"
                    ],
                    "description": "Layout mode for an assets or directory block."
                  },
                  "grid_size": {
                    "type": "integer",
                    "minimum": 1,
                    "maximum": 4,
                    "description": "Column count for an assets block using the grid layout."
                  }
                },
                "additionalProperties": false
              }
            }
          },
          "additionalProperties": false
        }
      },
      "additionalProperties": false
    },
    "annotations": {
      "readOnlyHint": false,
      "destructiveHint": false,
      "idempotentHint": true,
      "openWorldHint": false,
      "untrustedContentHint": false
    },
    "pathPrefix": "/l/"
  },
  "library_page_block_remove": {
    "name": "library_page_block_remove",
    "title": "Remove Library Page Block",
    "description": "Remove a visible block from the current Library page without deleting its content.",
    "inputSchema": {
      "type": "object",
      "required": [
        "block"
      ],
      "properties": {
        "block": {
          "type": "string",
          "enum": [
            "title",
            "cover",
            "link",
            "content",
            "assets",
            "properties",
            "directory",
            "tags"
          ]
        }
      },
      "additionalProperties": false
    },
    "outputSchema": {
      "type": "object",
      "required": [
        "message",
        "layout"
      ],
      "properties": {
        "message": {
          "type": "string",
          "description": "Summary of the removed block."
        },
        "layout": {
          "type": "object",
          "required": [
            "blocks"
          ],
          "properties": {
            "blocks": {
              "type": "array",
              "description": "Blocks in their rendered order.",
              "items": {
                "type": "object",
                "required": [
                  "type"
                ],
                "properties": {
                  "type": {
                    "type": "string",
                    "enum": [
                      "title",
                      "cover",
                      "link",
                      "content",
                      "assets",
                      "properties",
                      "directory",
                      "tags"
                    ]
                  },
                  "layout": {
                    "type": "string",
                    "enum": [
                      "strip",
                      "grid",
                      "table"
                    ],
                    "description": "Layout mode for an assets or directory block."
                  },
                  "grid_size": {
                    "type": "integer",
                    "minimum": 1,
                    "maximum": 4,
                    "description": "Column count for an assets block using the grid layout."
                  }
                },
                "additionalProperties": false
              }
            }
          },
          "additionalProperties": false
        }
      },
      "additionalProperties": false
    },
    "annotations": {
      "readOnlyHint": false,
      "destructiveHint": true,
      "idempotentHint": true,
      "openWorldHint": false,
      "untrustedContentHint": false
    },
    "pathPrefix": "/l/"
  },
  "library_page_block_move": {
    "name": "library_page_block_move",
    "title": "Move Library Page Block",
    "description": "Move a visible block within the current Library page layout.",
    "inputSchema": {
      "type": "object",
      "required": [
        "block"
      ],
      "properties": {
        "block": {
          "type": "string",
          "enum": [
            "title",
            "cover",
            "link",
            "content",
            "assets",
            "properties",
            "directory",
            "tags"
          ]
        },
        "before_block": {
          "description": "Existing block before which to place the block; omit to move it to the end.",
          "type": "string",
          "enum": [
            "title",
            "cover",
            "link",
            "content",
            "assets",
            "properties",
            "directory",
            "tags"
          ]
        }
      },
      "additionalProperties": false
    },
    "outputSchema": {
      "type": "object",
      "required": [
        "message",
        "layout"
      ],
      "properties": {
        "message": {
          "type": "string",
          "description": "Summary of the moved block."
        },
        "layout": {
          "type": "object",
          "required": [
            "blocks"
          ],
          "properties": {
            "blocks": {
              "type": "array",
              "description": "Blocks in their rendered order.",
              "items": {
                "type": "object",
                "required": [
                  "type"
                ],
                "properties": {
                  "type": {
                    "type": "string",
                    "enum": [
                      "title",
                      "cover",
                      "link",
                      "content",
                      "assets",
                      "properties",
                      "directory",
                      "tags"
                    ]
                  },
                  "layout": {
                    "type": "string",
                    "enum": [
                      "strip",
                      "grid",
                      "table"
                    ],
                    "description": "Layout mode for an assets or directory block."
                  },
                  "grid_size": {
                    "type": "integer",
                    "minimum": 1,
                    "maximum": 4,
                    "description": "Column count for an assets block using the grid layout."
                  }
                },
                "additionalProperties": false
              }
            }
          },
          "additionalProperties": false
        }
      },
      "additionalProperties": false
    },
    "annotations": {
      "readOnlyHint": false,
      "destructiveHint": false,
      "idempotentHint": true,
      "openWorldHint": false,
      "untrustedContentHint": false
    },
    "pathPrefix": "/l/"
  },
  "library_page_block_assets_update": {
    "name": "library_page_block_assets_update",
    "title": "Update Library Page Assets Block",
    "description": "Update the presentation of the visible assets block on the current Library page.",
    "inputSchema": {
      "type": "object",
      "required": [
        "layout"
      ],
      "properties": {
        "layout": {
          "type": "string",
          "enum": [
            "strip",
            "grid"
          ]
        },
        "grid_size": {
          "type": "integer",
          "minimum": 1,
          "maximum": 4,
          "description": "Number of columns used by the grid layout."
        }
      },
      "additionalProperties": false
    },
    "outputSchema": {
      "type": "object",
      "required": [
        "message",
        "layout"
      ],
      "properties": {
        "message": {
          "type": "string",
          "description": "Summary of the updated assets block."
        },
        "layout": {
          "type": "object",
          "required": [
            "blocks"
          ],
          "properties": {
            "blocks": {
              "type": "array",
              "description": "Blocks in their rendered order.",
              "items": {
                "type": "object",
                "required": [
                  "type"
                ],
                "properties": {
                  "type": {
                    "type": "string",
                    "enum": [
                      "title",
                      "cover",
                      "link",
                      "content",
                      "assets",
                      "properties",
                      "directory",
                      "tags"
                    ]
                  },
                  "layout": {
                    "type": "string",
                    "enum": [
                      "strip",
                      "grid",
                      "table"
                    ],
                    "description": "Layout mode for an assets or directory block."
                  },
                  "grid_size": {
                    "type": "integer",
                    "minimum": 1,
                    "maximum": 4,
                    "description": "Column count for an assets block using the grid layout."
                  }
                },
                "additionalProperties": false
              }
            }
          },
          "additionalProperties": false
        }
      },
      "additionalProperties": false
    },
    "annotations": {
      "readOnlyHint": false,
      "destructiveHint": false,
      "idempotentHint": true,
      "openWorldHint": false,
      "untrustedContentHint": false
    },
    "pathPrefix": "/l/"
  },
  "library_page_block_directory_update": {
    "name": "library_page_block_directory_update",
    "title": "Update Library Page Directory Block",
    "description": "Update the presentation of the visible directory block on the current Library page.",
    "inputSchema": {
      "type": "object",
      "required": [
        "layout"
      ],
      "properties": {
        "layout": {
          "type": "string",
          "enum": [
            "table",
            "grid"
          ]
        }
      },
      "additionalProperties": false
    },
    "outputSchema": {
      "type": "object",
      "required": [
        "message",
        "layout"
      ],
      "properties": {
        "message": {
          "type": "string",
          "description": "Summary of the updated directory block."
        },
        "layout": {
          "type": "object",
          "required": [
            "blocks"
          ],
          "properties": {
            "blocks": {
              "type": "array",
              "description": "Blocks in their rendered order.",
              "items": {
                "type": "object",
                "required": [
                  "type"
                ],
                "properties": {
                  "type": {
                    "type": "string",
                    "enum": [
                      "title",
                      "cover",
                      "link",
                      "content",
                      "assets",
                      "properties",
                      "directory",
                      "tags"
                    ]
                  },
                  "layout": {
                    "type": "string",
                    "enum": [
                      "strip",
                      "grid",
                      "table"
                    ],
                    "description": "Layout mode for an assets or directory block."
                  },
                  "grid_size": {
                    "type": "integer",
                    "minimum": 1,
                    "maximum": 4,
                    "description": "Column count for an assets block using the grid layout."
                  }
                },
                "additionalProperties": false
              }
            }
          },
          "additionalProperties": false
        }
      },
      "additionalProperties": false
    },
    "annotations": {
      "readOnlyHint": false,
      "destructiveHint": false,
      "idempotentHint": true,
      "openWorldHint": false,
      "untrustedContentHint": false
    },
    "pathPrefix": "/l/"
  },
  "index_page_edit_start": {
    "name": "index_page_edit_start",
    "title": "Start Index Page Quick Edit",
    "description": "Enter quick edit mode for the live index page so frontend-only block tools can make the immediate change the user requested. Use this only when the user clearly asked to edit the live page. If it is unclear whether they want a quick edit or a version draft, ask them before calling this tool.",
    "inputSchema": {
      "type": "object",
      "properties": {},
      "additionalProperties": false
    },
    "outputSchema": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "message": {
          "type": "string",
          "description": "Summary of the active edit mode."
        }
      },
      "additionalProperties": false
    },
    "annotations": {
      "readOnlyHint": false,
      "destructiveHint": false,
      "idempotentHint": true,
      "openWorldHint": false,
      "untrustedContentHint": false
    },
    "pathPrefix": "/"
  },
  "index_page_layout_get": {
    "name": "index_page_layout_get",
    "title": "Get Index Page Layout",
    "description": "Retrieve the block layout of the current index page.",
    "inputSchema": {
      "type": "object",
      "properties": {},
      "additionalProperties": false
    },
    "outputSchema": {
      "type": "object",
      "required": [
        "message",
        "layout"
      ],
      "properties": {
        "message": {
          "type": "string",
          "description": "Summary of the retrieved layout."
        },
        "layout": {
          "type": "object",
          "required": [
            "blocks"
          ],
          "properties": {
            "blocks": {
              "type": "array",
              "description": "Blocks in their rendered order.",
              "items": {
                "type": "object",
                "required": [
                  "type"
                ],
                "properties": {
                  "type": {
                    "type": "string",
                    "enum": [
                      "title",
                      "subtitle",
                      "content",
                      "cover",
                      "categories",
                      "threads",
                      "quick-share",
                      "library"
                    ]
                  },
                  "layout": {
                    "type": "string",
                    "enum": [
                      "list",
                      "grid"
                    ],
                    "description": "Layout mode for a categories or Library block."
                  },
                  "source": {
                    "type": "string",
                    "enum": [
                      "all",
                      "uncategorised"
                    ],
                    "description": "Thread source for a threads block."
                  },
                  "show_category_select": {
                    "type": "boolean",
                    "description": "Whether a quick-share block shows its category picker."
                  },
                  "node_id": {
                    "type": "string",
                    "description": "Selected Library page ID for a library block; omitted for the Library root."
                  }
                },
                "additionalProperties": false
              }
            }
          },
          "additionalProperties": false
        }
      },
      "additionalProperties": false
    },
    "annotations": {
      "readOnlyHint": true,
      "destructiveHint": false,
      "idempotentHint": true,
      "openWorldHint": false,
      "untrustedContentHint": false
    },
    "pathPrefix": "/"
  },
  "index_page_block_add": {
    "name": "index_page_block_add",
    "title": "Add Index Page Block",
    "description": "Add a currently hidden block to the current index page.",
    "inputSchema": {
      "type": "object",
      "required": [
        "block"
      ],
      "properties": {
        "block": {
          "type": "string",
          "enum": [
            "title",
            "subtitle",
            "content",
            "cover",
            "categories",
            "threads",
            "quick-share",
            "library"
          ]
        },
        "after_block": {
          "description": "Existing block after which to insert the block; omit to append it.",
          "type": "string",
          "enum": [
            "title",
            "subtitle",
            "content",
            "cover",
            "categories",
            "threads",
            "quick-share",
            "library"
          ]
        }
      },
      "additionalProperties": false
    },
    "outputSchema": {
      "type": "object",
      "required": [
        "message",
        "layout"
      ],
      "properties": {
        "message": {
          "type": "string",
          "description": "Summary of the added block."
        },
        "layout": {
          "type": "object",
          "required": [
            "blocks"
          ],
          "properties": {
            "blocks": {
              "type": "array",
              "description": "Blocks in their rendered order.",
              "items": {
                "type": "object",
                "required": [
                  "type"
                ],
                "properties": {
                  "type": {
                    "type": "string",
                    "enum": [
                      "title",
                      "subtitle",
                      "content",
                      "cover",
                      "categories",
                      "threads",
                      "quick-share",
                      "library"
                    ]
                  },
                  "layout": {
                    "type": "string",
                    "enum": [
                      "list",
                      "grid"
                    ],
                    "description": "Layout mode for a categories or Library block."
                  },
                  "source": {
                    "type": "string",
                    "enum": [
                      "all",
                      "uncategorised"
                    ],
                    "description": "Thread source for a threads block."
                  },
                  "show_category_select": {
                    "type": "boolean",
                    "description": "Whether a quick-share block shows its category picker."
                  },
                  "node_id": {
                    "type": "string",
                    "description": "Selected Library page ID for a library block; omitted for the Library root."
                  }
                },
                "additionalProperties": false
              }
            }
          },
          "additionalProperties": false
        }
      },
      "additionalProperties": false
    },
    "annotations": {
      "readOnlyHint": false,
      "destructiveHint": false,
      "idempotentHint": true,
      "openWorldHint": false,
      "untrustedContentHint": false
    },
    "pathPrefix": "/"
  },
  "index_page_block_remove": {
    "name": "index_page_block_remove",
    "title": "Remove Index Page Block",
    "description": "Remove a visible block from the current index page without deleting its content.",
    "inputSchema": {
      "type": "object",
      "required": [
        "block"
      ],
      "properties": {
        "block": {
          "type": "string",
          "enum": [
            "title",
            "subtitle",
            "content",
            "cover",
            "categories",
            "threads",
            "quick-share",
            "library"
          ]
        }
      },
      "additionalProperties": false
    },
    "outputSchema": {
      "type": "object",
      "required": [
        "message",
        "layout"
      ],
      "properties": {
        "message": {
          "type": "string",
          "description": "Summary of the removed block."
        },
        "layout": {
          "type": "object",
          "required": [
            "blocks"
          ],
          "properties": {
            "blocks": {
              "type": "array",
              "description": "Blocks in their rendered order.",
              "items": {
                "type": "object",
                "required": [
                  "type"
                ],
                "properties": {
                  "type": {
                    "type": "string",
                    "enum": [
                      "title",
                      "subtitle",
                      "content",
                      "cover",
                      "categories",
                      "threads",
                      "quick-share",
                      "library"
                    ]
                  },
                  "layout": {
                    "type": "string",
                    "enum": [
                      "list",
                      "grid"
                    ],
                    "description": "Layout mode for a categories or Library block."
                  },
                  "source": {
                    "type": "string",
                    "enum": [
                      "all",
                      "uncategorised"
                    ],
                    "description": "Thread source for a threads block."
                  },
                  "show_category_select": {
                    "type": "boolean",
                    "description": "Whether a quick-share block shows its category picker."
                  },
                  "node_id": {
                    "type": "string",
                    "description": "Selected Library page ID for a library block; omitted for the Library root."
                  }
                },
                "additionalProperties": false
              }
            }
          },
          "additionalProperties": false
        }
      },
      "additionalProperties": false
    },
    "annotations": {
      "readOnlyHint": false,
      "destructiveHint": true,
      "idempotentHint": true,
      "openWorldHint": false,
      "untrustedContentHint": false
    },
    "pathPrefix": "/"
  },
  "index_page_block_move": {
    "name": "index_page_block_move",
    "title": "Move Index Page Block",
    "description": "Move a visible block within the current index page layout.",
    "inputSchema": {
      "type": "object",
      "required": [
        "block"
      ],
      "properties": {
        "block": {
          "type": "string",
          "enum": [
            "title",
            "subtitle",
            "content",
            "cover",
            "categories",
            "threads",
            "quick-share",
            "library"
          ]
        },
        "before_block": {
          "description": "Existing block before which to place the block; omit to move it to the end.",
          "type": "string",
          "enum": [
            "title",
            "subtitle",
            "content",
            "cover",
            "categories",
            "threads",
            "quick-share",
            "library"
          ]
        }
      },
      "additionalProperties": false
    },
    "outputSchema": {
      "type": "object",
      "required": [
        "message",
        "layout"
      ],
      "properties": {
        "message": {
          "type": "string",
          "description": "Summary of the moved block."
        },
        "layout": {
          "type": "object",
          "required": [
            "blocks"
          ],
          "properties": {
            "blocks": {
              "type": "array",
              "description": "Blocks in their rendered order.",
              "items": {
                "type": "object",
                "required": [
                  "type"
                ],
                "properties": {
                  "type": {
                    "type": "string",
                    "enum": [
                      "title",
                      "subtitle",
                      "content",
                      "cover",
                      "categories",
                      "threads",
                      "quick-share",
                      "library"
                    ]
                  },
                  "layout": {
                    "type": "string",
                    "enum": [
                      "list",
                      "grid"
                    ],
                    "description": "Layout mode for a categories or Library block."
                  },
                  "source": {
                    "type": "string",
                    "enum": [
                      "all",
                      "uncategorised"
                    ],
                    "description": "Thread source for a threads block."
                  },
                  "show_category_select": {
                    "type": "boolean",
                    "description": "Whether a quick-share block shows its category picker."
                  },
                  "node_id": {
                    "type": "string",
                    "description": "Selected Library page ID for a library block; omitted for the Library root."
                  }
                },
                "additionalProperties": false
              }
            }
          },
          "additionalProperties": false
        }
      },
      "additionalProperties": false
    },
    "annotations": {
      "readOnlyHint": false,
      "destructiveHint": false,
      "idempotentHint": true,
      "openWorldHint": false,
      "untrustedContentHint": false
    },
    "pathPrefix": "/"
  },
  "index_page_block_categories_update": {
    "name": "index_page_block_categories_update",
    "title": "Update Index Page Categories Block",
    "description": "Update the presentation of the visible categories block on the current index page.",
    "inputSchema": {
      "type": "object",
      "required": [
        "layout"
      ],
      "properties": {
        "layout": {
          "type": "string",
          "enum": [
            "list",
            "grid"
          ]
        }
      },
      "additionalProperties": false
    },
    "outputSchema": {
      "type": "object",
      "required": [
        "message",
        "layout"
      ],
      "properties": {
        "message": {
          "type": "string",
          "description": "Summary of the updated categories block."
        },
        "layout": {
          "type": "object",
          "required": [
            "blocks"
          ],
          "properties": {
            "blocks": {
              "type": "array",
              "description": "Blocks in their rendered order.",
              "items": {
                "type": "object",
                "required": [
                  "type"
                ],
                "properties": {
                  "type": {
                    "type": "string",
                    "enum": [
                      "title",
                      "subtitle",
                      "content",
                      "cover",
                      "categories",
                      "threads",
                      "quick-share",
                      "library"
                    ]
                  },
                  "layout": {
                    "type": "string",
                    "enum": [
                      "list",
                      "grid"
                    ],
                    "description": "Layout mode for a categories or Library block."
                  },
                  "source": {
                    "type": "string",
                    "enum": [
                      "all",
                      "uncategorised"
                    ],
                    "description": "Thread source for a threads block."
                  },
                  "show_category_select": {
                    "type": "boolean",
                    "description": "Whether a quick-share block shows its category picker."
                  },
                  "node_id": {
                    "type": "string",
                    "description": "Selected Library page ID for a library block; omitted for the Library root."
                  }
                },
                "additionalProperties": false
              }
            }
          },
          "additionalProperties": false
        }
      },
      "additionalProperties": false
    },
    "annotations": {
      "readOnlyHint": false,
      "destructiveHint": false,
      "idempotentHint": true,
      "openWorldHint": false,
      "untrustedContentHint": false
    },
    "pathPrefix": "/"
  },
  "index_page_block_threads_update": {
    "name": "index_page_block_threads_update",
    "title": "Update Index Page Threads Block",
    "description": "Update the source of the visible threads block on the current index page.",
    "inputSchema": {
      "type": "object",
      "required": [
        "source"
      ],
      "properties": {
        "source": {
          "type": "string",
          "enum": [
            "all",
            "uncategorised"
          ]
        }
      },
      "additionalProperties": false
    },
    "outputSchema": {
      "type": "object",
      "required": [
        "message",
        "layout"
      ],
      "properties": {
        "message": {
          "type": "string",
          "description": "Summary of the updated threads block."
        },
        "layout": {
          "type": "object",
          "required": [
            "blocks"
          ],
          "properties": {
            "blocks": {
              "type": "array",
              "description": "Blocks in their rendered order.",
              "items": {
                "type": "object",
                "required": [
                  "type"
                ],
                "properties": {
                  "type": {
                    "type": "string",
                    "enum": [
                      "title",
                      "subtitle",
                      "content",
                      "cover",
                      "categories",
                      "threads",
                      "quick-share",
                      "library"
                    ]
                  },
                  "layout": {
                    "type": "string",
                    "enum": [
                      "list",
                      "grid"
                    ],
                    "description": "Layout mode for a categories or Library block."
                  },
                  "source": {
                    "type": "string",
                    "enum": [
                      "all",
                      "uncategorised"
                    ],
                    "description": "Thread source for a threads block."
                  },
                  "show_category_select": {
                    "type": "boolean",
                    "description": "Whether a quick-share block shows its category picker."
                  },
                  "node_id": {
                    "type": "string",
                    "description": "Selected Library page ID for a library block; omitted for the Library root."
                  }
                },
                "additionalProperties": false
              }
            }
          },
          "additionalProperties": false
        }
      },
      "additionalProperties": false
    },
    "annotations": {
      "readOnlyHint": false,
      "destructiveHint": false,
      "idempotentHint": true,
      "openWorldHint": false,
      "untrustedContentHint": false
    },
    "pathPrefix": "/"
  },
  "index_page_block_quick_share_update": {
    "name": "index_page_block_quick_share_update",
    "title": "Update Index Page Quick Share Block",
    "description": "Configure the category picker on the visible quick-share block on the current index page.",
    "inputSchema": {
      "type": "object",
      "required": [
        "show_category_select"
      ],
      "properties": {
        "show_category_select": {
          "type": "boolean"
        }
      },
      "additionalProperties": false
    },
    "outputSchema": {
      "type": "object",
      "required": [
        "message",
        "layout"
      ],
      "properties": {
        "message": {
          "type": "string",
          "description": "Summary of the updated quick-share block."
        },
        "layout": {
          "type": "object",
          "required": [
            "blocks"
          ],
          "properties": {
            "blocks": {
              "type": "array",
              "description": "Blocks in their rendered order.",
              "items": {
                "type": "object",
                "required": [
                  "type"
                ],
                "properties": {
                  "type": {
                    "type": "string",
                    "enum": [
                      "title",
                      "subtitle",
                      "content",
                      "cover",
                      "categories",
                      "threads",
                      "quick-share",
                      "library"
                    ]
                  },
                  "layout": {
                    "type": "string",
                    "enum": [
                      "list",
                      "grid"
                    ],
                    "description": "Layout mode for a categories or Library block."
                  },
                  "source": {
                    "type": "string",
                    "enum": [
                      "all",
                      "uncategorised"
                    ],
                    "description": "Thread source for a threads block."
                  },
                  "show_category_select": {
                    "type": "boolean",
                    "description": "Whether a quick-share block shows its category picker."
                  },
                  "node_id": {
                    "type": "string",
                    "description": "Selected Library page ID for a library block; omitted for the Library root."
                  }
                },
                "additionalProperties": false
              }
            }
          },
          "additionalProperties": false
        }
      },
      "additionalProperties": false
    },
    "annotations": {
      "readOnlyHint": false,
      "destructiveHint": false,
      "idempotentHint": true,
      "openWorldHint": false,
      "untrustedContentHint": false
    },
    "pathPrefix": "/"
  },
  "index_page_block_library_update": {
    "name": "index_page_block_library_update",
    "title": "Update Index Page Library Block",
    "description": "Configure the visible Library block on the current index page.",
    "inputSchema": {
      "type": "object",
      "required": [
        "source"
      ],
      "properties": {
        "source": {
          "type": "string",
          "enum": [
            "library_root",
            "library_page"
          ],
          "description": "Show the Library root or a specific Library page."
        },
        "node_id": {
          "type": "string",
          "minLength": 1,
          "description": "Required when source is library_page; omit when source is library_root."
        },
        "layout": {
          "type": "string",
          "enum": [
            "list",
            "grid"
          ],
          "description": "Layout for the Library root; omit to preserve it. Specific pages control their own layout."
        }
      },
      "additionalProperties": false
    },
    "outputSchema": {
      "type": "object",
      "required": [
        "message",
        "layout"
      ],
      "properties": {
        "message": {
          "type": "string",
          "description": "Summary of the updated Library block."
        },
        "layout": {
          "type": "object",
          "required": [
            "blocks"
          ],
          "properties": {
            "blocks": {
              "type": "array",
              "description": "Blocks in their rendered order.",
              "items": {
                "type": "object",
                "required": [
                  "type"
                ],
                "properties": {
                  "type": {
                    "type": "string",
                    "enum": [
                      "title",
                      "subtitle",
                      "content",
                      "cover",
                      "categories",
                      "threads",
                      "quick-share",
                      "library"
                    ]
                  },
                  "layout": {
                    "type": "string",
                    "enum": [
                      "list",
                      "grid"
                    ],
                    "description": "Layout mode for a categories or Library block."
                  },
                  "source": {
                    "type": "string",
                    "enum": [
                      "all",
                      "uncategorised"
                    ],
                    "description": "Thread source for a threads block."
                  },
                  "show_category_select": {
                    "type": "boolean",
                    "description": "Whether a quick-share block shows its category picker."
                  },
                  "node_id": {
                    "type": "string",
                    "description": "Selected Library page ID for a library block; omitted for the Library root."
                  }
                },
                "additionalProperties": false
              }
            }
          },
          "additionalProperties": false
        }
      },
      "additionalProperties": false
    },
    "annotations": {
      "readOnlyHint": false,
      "destructiveHint": false,
      "idempotentHint": true,
      "openWorldHint": false,
      "untrustedContentHint": false
    },
    "pathPrefix": "/"
  }
} as const;
