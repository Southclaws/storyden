/* eslint-disable */
/**
 * This file was automatically generated from web/tools/schema.yaml.
 * DO NOT MODIFY IT BY HAND. Run: pnpm webmcp-schema
 */

export interface WebMCPTools {
  ToolLibraryPageLayoutGet?: LibraryPageLayoutGet;
  ToolLibraryPageBlockAdd?: LibraryPageBlockAdd;
  ToolLibraryPageBlockRemove?: LibraryPageBlockRemove;
  ToolLibraryPageBlockMove?: LibraryPageBlockMove;
  ToolLibraryPageBlockAssetsUpdate?: LibraryPageBlockAssetsUpdate;
  ToolLibraryPageBlockDirectoryUpdate?: LibraryPageBlockDirectoryUpdate;
  [k: string]: unknown;
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

export type WebMCPToolName = "library_page_layout_get" | "library_page_block_add" | "library_page_block_remove" | "library_page_block_move" | "library_page_block_assets_update" | "library_page_block_directory_update";

export const WEBMCP_TOOL_NAMES = ["library_page_layout_get", "library_page_block_add", "library_page_block_remove", "library_page_block_move", "library_page_block_assets_update", "library_page_block_directory_update"] as const;

export type WebMCPToolInputMap = {
  "library_page_layout_get": ToolLibraryPageLayoutGetInput;
  "library_page_block_add": ToolLibraryPageBlockAddInput;
  "library_page_block_remove": ToolLibraryPageBlockRemoveInput;
  "library_page_block_move": ToolLibraryPageBlockMoveInput;
  "library_page_block_assets_update": ToolLibraryPageBlockAssetsUpdateInput;
  "library_page_block_directory_update": ToolLibraryPageBlockDirectoryUpdateInput;
};

export type WebMCPToolOutputMap = {
  "library_page_layout_get": ToolLibraryPageLayoutGetOutput;
  "library_page_block_add": ToolLibraryPageBlockAddOutput;
  "library_page_block_remove": ToolLibraryPageBlockRemoveOutput;
  "library_page_block_move": ToolLibraryPageBlockMoveOutput;
  "library_page_block_assets_update": ToolLibraryPageBlockAssetsUpdateOutput;
  "library_page_block_directory_update": ToolLibraryPageBlockDirectoryUpdateOutput;
};

export type WebMCPToolImplementation<K extends WebMCPToolName> = (
  input: WebMCPToolInputMap[K],
) => WebMCPToolOutputMap[K] | Promise<WebMCPToolOutputMap[K]>;

export type WebMCPToolImplementations = {
  [K in WebMCPToolName]: WebMCPToolImplementation<K>;
};

export const WEBMCP_TOOL_DEFINITIONS = {
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
  }
} as const;
