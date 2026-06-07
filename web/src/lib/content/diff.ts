import type { JSONContent } from "@tiptap/core";

import {
  diff_match_patch as DiffMatchPatch,
  DIFF_DELETE,
  DIFF_INSERT,
  DIFF_EQUAL,
} from "@/utils/diff-match-patch";

// NOTE: This one was fully AI written using Claude including the tests.

/**
 * Create a merged diff document showing both insertions and deletions inline.
 * This uses text-level diffing within matching structural nodes.
 */
export function diffTipTapJSON(
  original: JSONContent,
  modified: JSONContent,
): JSONContent {
  // Create a merged document by recursively comparing nodes
  const merged = mergeNodes(original, modified);
  return merged;
}

/**
 * Recursively merge two nodes, creating a document with diff marks
 */
function mergeNodes(
  originalNode: JSONContent,
  modifiedNode: JSONContent,
): JSONContent {
  // If node types differ, show as block replacement
  if (originalNode.type !== modifiedNode.type) {
    // Return modified with all content marked as insertion
    const result = JSON.parse(JSON.stringify(modifiedNode));
    markAllText(result, "insertion");
    return result;
  }

  // Handle text nodes with character-level diff
  if (originalNode.type === "text" && modifiedNode.type === "text") {
    return mergeTextNodes(originalNode, modifiedNode);
  }

  // Handle nodes with content (like paragraphs, headings, etc.)
  if (originalNode.content || modifiedNode.content) {
    const origContent = originalNode.content || [];
    const modContent = modifiedNode.content || [];

    // Special case: if both have single text children, merge them inline
    const firstOrig = origContent[0];
    const firstMod = modContent[0];

    if (
      origContent.length === 1 &&
      modContent.length === 1 &&
      firstOrig &&
      firstMod &&
      firstOrig.type === "text" &&
      firstMod.type === "text"
    ) {
      const mergedTextNodes = mergeTextNodesInline(firstOrig, firstMod);
      return {
        ...modifiedNode,
        content: mergedTextNodes,
      };
    }

    // Otherwise, merge content arrays
    const mergedContent = mergeContentArrays(origContent, modContent);

    return {
      ...modifiedNode,
      content: mergedContent,
    };
  }

  // Leaf nodes without content - return modified
  return modifiedNode;
}

/**
 * Merge two text nodes using character-level diff.
 * Returns a SINGLE node (for use when text nodes are isolated).
 */
function mergeTextNodes(
  originalNode: JSONContent,
  modifiedNode: JSONContent,
): JSONContent {
  const originalText = originalNode.text || "";
  const modifiedText = modifiedNode.text || "";

  // If texts are identical, return as-is
  if (originalText === modifiedText) {
    return modifiedNode;
  }

  // For isolated text nodes, just mark the whole thing as changed
  // This is a simplification
  return {
    type: "text",
    text: modifiedText,
    marks: [
      ...(modifiedNode.marks || []),
      { type: "diffInsertion", attrs: { "data-diff": "insertion" } },
    ],
  };
}

/**
 * Merge two text nodes inline, returning an ARRAY of text nodes with diff marks.
 * This handles cases like "Welcome to Storyden" → "Welcome to Storyden Community"
 * Result: ["Welcome to Storyden" (deletion), " Community" (insertion)]
 */
function mergeTextNodesInline(
  originalNode: JSONContent,
  modifiedNode: JSONContent,
): JSONContent[] {
  const originalText = originalNode.text || "";
  const modifiedText = modifiedNode.text || "";

  // If texts are identical, return as single node
  if (originalText === modifiedText) {
    return [modifiedNode];
  }

  // Use diff-match-patch for character-level diff
  const dmp = new DiffMatchPatch();
  const diffs = dmp.diff_main(originalText, modifiedText);
  dmp.diff_cleanupSemantic(diffs);

  // Convert diffs to an array of text nodes with appropriate marks
  const resultNodes: JSONContent[] = [];

  // Diff format is an array of [operation, text] tuples
  for (let i = 0; i < diffs.length; i++) {
    const diff = diffs[i];
    if (!diff) continue;

    const operation = diff[0];
    const text = diff[1];

    if (operation === DIFF_EQUAL) {
      // Unchanged text - preserve marks
      resultNodes.push({
        type: "text",
        text,
        marks: modifiedNode.marks || [],
      });
    } else if (operation === DIFF_DELETE) {
      // Deleted text - show with deletion mark
      resultNodes.push({
        type: "text",
        text,
        marks: [
          ...(originalNode.marks || []),
          { type: "diffDeletion", attrs: { "data-diff": "deletion" } },
        ],
      });
    } else if (operation === DIFF_INSERT) {
      // Inserted text - show with insertion mark
      resultNodes.push({
        type: "text",
        text,
        marks: [
          ...(modifiedNode.marks || []),
          { type: "diffInsertion", attrs: { "data-diff": "insertion" } },
        ],
      });
    }
  }

  return resultNodes;
}

/**
 * Merge two content arrays (arrays of child nodes) using LCS-based diffing.
 * This finds the best alignment between nodes, handling insertions and deletions correctly.
 */
function mergeContentArrays(
  originalContent: JSONContent[],
  modifiedContent: JSONContent[],
): JSONContent[] {
  const result: JSONContent[] = [];

  // Build LCS matrix to find optimal alignment
  const lcs = computeLCS(originalContent, modifiedContent);

  let i = 0; // pointer in original
  let j = 0; // pointer in modified

  while (i < originalContent.length || j < modifiedContent.length) {
    const origChild = originalContent[i];
    const modChild = modifiedContent[j];

    // Check if current nodes match in the LCS
    const isMatch =
      origChild &&
      modChild &&
      !nodesAreDifferent(origChild, modChild) &&
      (lcs[i]?.[j] ?? 0) === ((lcs[i + 1]?.[j + 1] ?? 0) + 1);

    if (isMatch) {
      // Nodes match - recursively merge them
      result.push(mergeNodes(origChild, modChild));
      i++;
      j++;
    } else if (
      origChild &&
      (!modChild || (lcs[i + 1]?.[j] ?? 0) >= (lcs[i]?.[j + 1] ?? 0))
    ) {
      // Node was deleted from original
      const marked = JSON.parse(JSON.stringify(origChild));
      markAllText(marked, "deletion");
      result.push(marked);
      i++;
    } else if (modChild) {
      // Node was added in modified
      const marked = JSON.parse(JSON.stringify(modChild));
      markAllText(marked, "insertion");
      result.push(marked);
      j++;
    }
  }

  return result;
}

/**
 * Compute LCS (Longest Common Subsequence) matrix for two content arrays.
 * This helps find the optimal alignment between nodes.
 */
function computeLCS(
  arr1: JSONContent[],
  arr2: JSONContent[],
): number[][] {
  const m = arr1.length;
  const n = arr2.length;
  const dp: number[][] = Array(m + 1)
    .fill(0)
    .map(() => Array(n + 1).fill(0));

  for (let i = m - 1; i >= 0; i--) {
    for (let j = n - 1; j >= 0; j--) {
      const node1 = arr1[i];
      const node2 = arr2[j];
      const currentRow = dp[i];
      const nextRow = dp[i + 1];

      if (!currentRow) continue;

      if (node1 && node2 && nodeSignature(node1) === nodeSignature(node2)) {
        currentRow[j] = 1 + (nextRow?.[j + 1] ?? 0);
      } else {
        currentRow[j] = Math.max(nextRow?.[j] ?? 0, currentRow[j + 1] ?? 0);
      }
    }
  }

  return dp;
}

/**
 * Generate a signature for a node to determine if two nodes are "the same" for LCS purposes.
 * This uses type and key attributes to identify nodes, focusing on structural identity
 * rather than content similarity.
 */
function nodeSignature(node: JSONContent): string {
  if (!node) return "";

  // For text nodes, use the text content (they need exact match)
  if (node.type === "text") {
    return `text:${node.text || ""}`;
  }

  // For custom nodes with unique identifiers (these are content-based), use those
  if (node.type === "linkPreview" && node.attrs?.["href"]) {
    return `linkPreview:${node.attrs["href"]}`;
  }

  if (node.type === "image" && node.attrs?.["src"]) {
    return `image:${node.attrs["src"]}`;
  }

  // For headings, include the level in signature
  if (node.type === "heading" && node.attrs?.["level"]) {
    return `heading:${node.attrs["level"]}`;
  }

  // For other structural nodes (paragraph, blockquote, list items, etc.),
  // just use the type - we'll diff their content, not their structure
  return node.type || "";
}

/**
 * Check if two nodes are significantly different.
 * Used to decide whether to show both versions or merge them.
 */
function nodesAreDifferent(node1: JSONContent, node2: JSONContent): boolean {
  // Different types = definitely different
  if (node1.type !== node2.type) {
    return true;
  }

  // For certain node types, check if attrs differ significantly
  // This handles cases like linkPreview with different hrefs
  if (node1.type === "linkPreview" || node1.type === "image") {
    const attrs1 = node1.attrs || {};
    const attrs2 = node2.attrs || {};

    // Compare key attributes
    if (node1.type === "linkPreview") {
      return attrs1["href"] !== attrs2["href"];
    }
    if (node1.type === "image") {
      return attrs1["src"] !== attrs2["src"];
    }
  }

  // For text nodes, never consider them "different" here - we handle them specially
  if (node1.type === "text") {
    return false;
  }

  // For structural nodes (paragraph, heading, etc.), try to merge
  return false;
}

/**
 * Mark all text nodes within a tree with a diff mark.
 * For non-text nodes (like linkPreview, image), add data-diff attribute.
 */
function markAllText(node: JSONContent, type: "insertion" | "deletion"): void {
  if (node.type === "text") {
    if (!node.marks) {
      node.marks = [];
    }

    node.marks.push({
      type: type === "insertion" ? "diffInsertion" : "diffDeletion",
      attrs: { "data-diff": type },
    });
  } else {
    // For non-text nodes (custom nodes like linkPreview, image, etc.)
    // add data-diff attribute so they can be styled
    if (!node.attrs) {
      node.attrs = {};
    }
    node.attrs["data-diff"] = type;
  }

  if (Array.isArray(node.content)) {
    node.content.forEach((child) => markAllText(child, type));
  }
}

/**
 * Count diff marks in a document
 */
export function countDiffMarks(doc: JSONContent): {
  insertions: number;
  deletions: number;
} {
  let insertions = 0;
  let deletions = 0;

  function traverse(node: JSONContent) {
    if (node.marks) {
      for (const mark of node.marks) {
        if (mark.type === "diffInsertion") insertions++;
        if (mark.type === "diffDeletion") deletions++;
      }
    }

    if (Array.isArray(node.content)) {
      node.content.forEach(traverse);
    }
  }

  traverse(doc);

  return { insertions, deletions };
}
