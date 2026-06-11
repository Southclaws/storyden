import type { JSONContent } from "@tiptap/core";
import { test } from "uvu";
import * as assert from "uvu/assert";

import { countDiffMarks, diffTipTapJSON } from "./diff";

test("identical documents produce no diff marks", () => {
  const original: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "Hello World" }],
      },
    ],
  };

  const modified: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "Hello World" }],
      },
    ],
  };

  const result = diffTipTapJSON(original, modified);
  const counts = countDiffMarks(result);

  assert.equal(counts.insertions, 0);
  assert.equal(counts.deletions, 0);
});

test("simple text insertion", () => {
  const original: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "Hello World" }],
      },
    ],
  };

  const modified: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "Hello World Community" }],
      },
    ],
  };

  const result = diffTipTapJSON(original, modified);
  const counts = countDiffMarks(result);

  assert.equal(counts.insertions, 1);
  assert.equal(counts.deletions, 0);

  // Check the merged content has both unchanged and inserted text
  const paragraph = result.content?.[0];
  assert.ok(paragraph?.content);
  assert.ok(paragraph.content.length > 1, "Should have multiple text nodes");

  // Find the insertion node
  const insertionNode = paragraph.content.find((node) =>
    node.marks?.some((mark) => mark.type === "diffInsertion"),
  );
  assert.ok(insertionNode, "Should have insertion mark");
  assert.ok(insertionNode?.text?.includes("Community"));
});

test("simple text deletion", () => {
  const original: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "Hello World Community" }],
      },
    ],
  };

  const modified: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "Hello World" }],
      },
    ],
  };

  const result = diffTipTapJSON(original, modified);
  const counts = countDiffMarks(result);

  assert.equal(counts.insertions, 0);
  assert.equal(counts.deletions, 1);

  // Check the merged content has both unchanged and deleted text
  const paragraph = result.content?.[0];
  assert.ok(paragraph?.content);

  // Find the deletion node
  const deletionNode = paragraph.content.find((node) =>
    node.marks?.some((mark) => mark.type === "diffDeletion"),
  );
  assert.ok(deletionNode, "Should have deletion mark");
  assert.ok(deletionNode?.text?.includes("Community"));
});

test("text replacement shows both deletion and insertion", () => {
  const original: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "Hello World" }],
      },
    ],
  };

  const modified: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "Goodbye Universe" }],
      },
    ],
  };

  const result = diffTipTapJSON(original, modified);
  const counts = countDiffMarks(result);

  assert.ok(counts.insertions > 0, "Should have insertions");
  assert.ok(counts.deletions > 0, "Should have deletions");
});

test("paragraph addition", () => {
  const original: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "First paragraph" }],
      },
    ],
  };

  const modified: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "First paragraph" }],
      },
      {
        type: "paragraph",
        content: [{ type: "text", text: "Second paragraph" }],
      },
    ],
  };

  const result = diffTipTapJSON(original, modified);
  const counts = countDiffMarks(result);

  assert.equal(counts.insertions, 1);
  assert.equal(counts.deletions, 0);
  assert.equal(result.content?.length, 2, "Should have 2 paragraphs");
});

test("paragraph deletion", () => {
  const original: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "First paragraph" }],
      },
      {
        type: "paragraph",
        content: [{ type: "text", text: "Second paragraph" }],
      },
    ],
  };

  const modified: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "First paragraph" }],
      },
    ],
  };

  const result = diffTipTapJSON(original, modified);
  const counts = countDiffMarks(result);

  assert.equal(counts.insertions, 0);
  assert.equal(counts.deletions, 1);
  assert.equal(result.content?.length, 2, "Should show both paragraphs");
});

test("preserves text marks (bold, italic)", () => {
  const original: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [
          {
            type: "text",
            text: "Hello ",
            marks: [{ type: "bold" }],
          },
          {
            type: "text",
            text: "World",
          },
        ],
      },
    ],
  };

  const modified: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [
          {
            type: "text",
            text: "Hello ",
            marks: [{ type: "bold" }],
          },
          {
            type: "text",
            text: "World Community",
          },
        ],
      },
    ],
  };

  const result = diffTipTapJSON(original, modified);

  // Find the paragraph content
  const paragraph = result.content?.[0];
  assert.ok(paragraph?.content);

  // The first text node should still have bold mark
  const boldNode = paragraph.content.find((node) =>
    node.marks?.some((mark) => mark.type === "bold"),
  );
  assert.ok(boldNode, "Should preserve bold mark");
});

test("linkPreview node with different href shows both versions", () => {
  const original: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [
          {
            type: "linkPreview",
            attrs: {
              href: "https://example.com",
              display: "card",
            },
          },
        ],
      },
    ],
  };

  const modified: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [
          {
            type: "linkPreview",
            attrs: {
              href: "https://barney.is",
              display: "card",
            },
          },
        ],
      },
    ],
  };

  const result = diffTipTapJSON(original, modified);

  // Should have 2 linkPreview nodes in the paragraph
  const paragraph = result.content?.[0];
  assert.ok(paragraph?.content);
  assert.equal(
    paragraph.content.length,
    2,
    "Should show both old and new linkPreview",
  );

  const linkPreviews = paragraph.content.filter(
    (node) => node.type === "linkPreview",
  );
  assert.equal(linkPreviews.length, 2);

  // Check that hrefs are preserved
  const hrefs = linkPreviews.map((node) => node.attrs?.["href"]);
  assert.ok(hrefs.includes("https://example.com"));
  assert.ok(hrefs.includes("https://barney.is"));

  // Check that data-diff attributes are set
  const deletionNode = linkPreviews.find(
    (node) => node.attrs?.["data-diff"] === "deletion",
  );
  const insertionNode = linkPreviews.find(
    (node) => node.attrs?.["data-diff"] === "insertion",
  );
  assert.ok(deletionNode, "Should have deletion data-diff attribute");
  assert.ok(insertionNode, "Should have insertion data-diff attribute");
  assert.equal(deletionNode?.attrs?.["href"], "https://example.com");
  assert.equal(insertionNode?.attrs?.["href"], "https://barney.is");
});

test("linkPreview node with same href does not duplicate", () => {
  const original: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [
          {
            type: "linkPreview",
            attrs: {
              href: "https://example.com",
              display: "card",
            },
          },
        ],
      },
    ],
  };

  const modified: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [
          {
            type: "linkPreview",
            attrs: {
              href: "https://example.com",
              display: "card",
            },
          },
        ],
      },
    ],
  };

  const result = diffTipTapJSON(original, modified);

  // Should have only 1 linkPreview node since they're identical
  const paragraph = result.content?.[0];
  assert.ok(paragraph?.content);
  assert.equal(paragraph.content.length, 1, "Should not duplicate same link");
});

test("image node with different src shows both versions", () => {
  const original: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [
          {
            type: "image",
            attrs: {
              src: "https://example.com/image1.jpg",
            },
          },
        ],
      },
    ],
  };

  const modified: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [
          {
            type: "image",
            attrs: {
              src: "https://example.com/image2.jpg",
            },
          },
        ],
      },
    ],
  };

  const result = diffTipTapJSON(original, modified);

  // Should have 2 image nodes in the paragraph
  const paragraph = result.content?.[0];
  assert.ok(paragraph?.content);
  assert.equal(
    paragraph.content.length,
    2,
    "Should show both old and new image",
  );

  const images = paragraph.content.filter((node) => node.type === "image");
  assert.equal(images.length, 2);

  // Check that data-diff attributes are set
  const deletionNode = images.find(
    (node) => node.attrs?.["data-diff"] === "deletion",
  );
  const insertionNode = images.find(
    (node) => node.attrs?.["data-diff"] === "insertion",
  );
  assert.ok(deletionNode, "Should have deletion data-diff attribute");
  assert.ok(insertionNode, "Should have insertion data-diff attribute");
  assert.equal(deletionNode?.attrs?.["src"], "https://example.com/image1.jpg");
  assert.equal(insertionNode?.attrs?.["src"], "https://example.com/image2.jpg");
});

test("heading changes preserve structure", () => {
  const original: JSONContent = {
    type: "doc",
    content: [
      {
        type: "heading",
        attrs: { level: 1 },
        content: [{ type: "text", text: "Original Title" }],
      },
    ],
  };

  const modified: JSONContent = {
    type: "doc",
    content: [
      {
        type: "heading",
        attrs: { level: 1 },
        content: [{ type: "text", text: "Modified Title" }],
      },
    ],
  };

  const result = diffTipTapJSON(original, modified);

  // Should still be a heading
  assert.equal(result.content?.[0]?.type, "heading");
  assert.equal(result.content?.[0]?.attrs?.["level"], 1);

  const counts = countDiffMarks(result);
  assert.ok(counts.insertions > 0);
  assert.ok(counts.deletions > 0);
});

test("complex document with multiple changes", () => {
  const original: JSONContent = {
    type: "doc",
    content: [
      {
        type: "heading",
        attrs: { level: 1 },
        content: [{ type: "text", text: "Welcome" }],
      },
      {
        type: "paragraph",
        content: [{ type: "text", text: "This is the introduction." }],
      },
      {
        type: "paragraph",
        content: [
          {
            type: "linkPreview",
            attrs: {
              href: "https://example.com",
              display: "card",
            },
          },
        ],
      },
    ],
  };

  const modified: JSONContent = {
    type: "doc",
    content: [
      {
        type: "heading",
        attrs: { level: 1 },
        content: [{ type: "text", text: "Welcome to Storyden" }],
      },
      {
        type: "paragraph",
        content: [{ type: "text", text: "This is the updated introduction." }],
      },
      {
        type: "paragraph",
        content: [
          {
            type: "linkPreview",
            attrs: {
              href: "https://barney.is",
              display: "card",
            },
          },
        ],
      },
    ],
  };

  const result = diffTipTapJSON(original, modified);
  const counts = countDiffMarks(result);

  // Should have insertions (text added to heading and paragraph, plus linkPreview change)
  assert.ok(counts.insertions > 0, "Should have insertions");

  // Should preserve document structure
  assert.equal(result.content?.length, 3, "Should have 3 top-level nodes");

  // Check that linkPreview nodes are both present
  const lastParagraph = result.content?.[2];
  const linkPreviews = lastParagraph?.content?.filter(
    (node) => node.type === "linkPreview",
  );
  assert.equal(
    linkPreviews?.length,
    2,
    "Should have both old and new linkPreview nodes",
  );
});

test("entire blockquote addition has data-diff attribute", () => {
  const original: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "Some text" }],
      },
    ],
  };

  const modified: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "Some text" }],
      },
      {
        type: "blockquote",
        content: [
          {
            type: "paragraph",
            content: [{ type: "text", text: "This is a new quote" }],
          },
        ],
      },
    ],
  };

  const result = diffTipTapJSON(original, modified);

  // Should have 2 nodes: paragraph (unchanged) and blockquote (added)
  assert.equal(result.content?.length, 2);

  const blockquote = result.content?.[1];
  assert.equal(blockquote?.type, "blockquote");
  assert.equal(blockquote?.attrs?.["data-diff"], "insertion");
});

test("entire blockquote removal has data-diff attribute", () => {
  const original: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "Some text" }],
      },
      {
        type: "blockquote",
        content: [
          {
            type: "paragraph",
            content: [{ type: "text", text: "This quote will be removed" }],
          },
        ],
      },
    ],
  };

  const modified: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "Some text" }],
      },
    ],
  };

  const result = diffTipTapJSON(original, modified);

  // Should have 2 nodes: paragraph (unchanged) and blockquote (deleted)
  assert.equal(result.content?.length, 2);

  const blockquote = result.content?.[1];
  assert.equal(blockquote?.type, "blockquote");
  assert.equal(blockquote?.attrs?.["data-diff"], "deletion");
});

test("inserting paragraph above linkPreview doesn't mark link as changed", () => {
  const original: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "Content 1" }],
      },
      {
        type: "linkPreview",
        attrs: {
          href: "https://barney.is/",
          "data-display": "card",
        },
      },
      {
        type: "paragraph",
        content: [],
      },
    ],
  };

  const modified: JSONContent = {
    type: "doc",
    content: [
      {
        type: "paragraph",
        content: [{ type: "text", text: "Content 2" }],
      },
      {
        type: "paragraph",
        content: [{ type: "text", text: "next" }],
      },
      {
        type: "linkPreview",
        attrs: {
          href: "https://barney.is/",
          "data-display": "card",
        },
      },
      {
        type: "paragraph",
        content: [],
      },
    ],
  };

  const result = diffTipTapJSON(original, modified);

  // Should have 4 nodes in result
  assert.equal(result.content?.length, 4);

  // First paragraph should have text changes (Content 1 -> Content 2)
  const firstPara = result.content?.[0];
  assert.equal(firstPara?.type, "paragraph");

  // Second paragraph should be marked as insertion (new "next" paragraph)
  const secondPara = result.content?.[1];
  assert.equal(secondPara?.type, "paragraph");
  assert.ok(secondPara?.attrs?.["data-diff"] === "insertion");

  // LinkPreview should NOT be marked as changed (same href)
  const linkPreview = result.content?.[2];
  assert.equal(linkPreview?.type, "linkPreview");
  assert.equal(linkPreview?.attrs?.["href"], "https://barney.is/");
  assert.ok(
    !linkPreview?.attrs?.["data-diff"],
    "LinkPreview should not have data-diff attribute since it didn't change",
  );
});

test.run();
