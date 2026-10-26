import { describe, expect, it } from "vitest";

import { parseNodeMetadata } from "./metadata";

describe("Main image metadata", () => {
  it("keeps an existing cover block and old crop metadata without adding config", () => {
    const metadata = {
      coverImage: { top: 40, left: 20 },
      layout: {
        blocks: [{ type: "title" }, { type: "cover" }, { type: "content" }],
      },
    };
    expect(parseNodeMetadata(metadata)).toEqual(metadata);
  });

  it.each(["banner", "image", "artwork"])(
    "retains additive %s config and neighbouring blocks",
    (layout) => {
      const metadata = {
        coverImage: { top: 0, left: 0 },
        layout: {
          blocks: [
            {
              type: "cover",
              config: {
                layout,
                alt: "A landscape",
                crop: { left: 0.25, top: 0.1, width: 0.5, height: 0.8 },
              },
            },
            { type: "content" },
          ],
        },
      };
      expect(parseNodeMetadata(metadata)).toEqual(metadata);
    },
  );
});
