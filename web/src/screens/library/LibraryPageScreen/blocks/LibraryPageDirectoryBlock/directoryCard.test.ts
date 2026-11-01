import { describe, expect, it } from "vitest";

import { Asset, LinkReference, NodeWithChildren } from "@/api/openapi-schema";

import { ColumnDefinition, ColumnValue } from "./column";
import {
  applyDirectoryDraft,
  getDirectoryBadgeColumns,
  getRichCardImageAspectRatio,
  getRichCardShape,
  isRichCardFieldVisible,
} from "./directoryCard";

const image: Asset = {
  id: "cover",
  path: "/assets/cover",
  filename: "cover.png",
  mime_type: "image/png",
  width: 600,
  height: 1200,
};

const sourceLink: LinkReference = {
  id: "link",
  url: "https://example.com/font",
  slug: "font",
  domain: "example.com",
  title: "Typeface website",
  createdAt: "2026-11-01T00:00:00Z",
  updatedAt: "2026-11-01T00:00:00Z",
};

const source: NodeWithChildren = {
  id: "font",
  slug: "font",
  name: "Published font",
  description: "Published description",
  content: "<p>Full page content.</p>",
  primary_image: image,
  assets: [image],
  createdAt: "2026-11-01T00:00:00Z",
  updatedAt: "2026-11-01T00:00:00Z",
  owner: {
    id: "author",
    handle: "author",
    name: "Author",
    joined: "2026-11-01T00:00:00Z",
    roles: [],
  },
  visibility: "published",
  hide_child_tree: false,
  meta: {},
  tags: [{ name: "typography", colour: "#808080", item_count: 1 }],
  children: [],
  child_property_schema: [],
  recomentations: [],
  properties: [],
  link: sourceLink,
};

describe("Directory draft overlay", () => {
  it("uses the source directly when there is no local draft", () => {
    expect(applyDirectoryDraft(source, undefined)).toBe(source);
  });

  it("applies editable fields without losing media, rich content, or metadata", () => {
    const properties = [
      {
        fid: "style",
        name: "Style",
        type: "text",
        sort: "0",
        value: "Sans-serif",
      },
    ] as const;
    const draft: NodeWithChildren = {
      ...source,
      name: "Draft font",
      description: "",
      properties: [...properties],
      primary_image: undefined,
      content: undefined,
      link: undefined,
      assets: [],
      tags: [],
    };

    const result = applyDirectoryDraft(source, draft);

    expect(result.name).toBe("Draft font");
    expect(result.description).toBe("");
    expect(result.properties).toEqual(properties);
    expect(result.primary_image).toBe(image);
    expect(result.assets).toBe(source.assets);
    expect(result.tags).toBe(source.tags);
    expect(result.link).toBe(source.link);
    expect(result.content).toBe("<p>Full page content.</p>");
    expect(source.name).toBe("Published font");
    expect(source.description).toBe("Published description");
  });

  it("retains an explicitly cleared draft URL", () => {
    const draft = { ...source, link: { ...sourceLink, url: "" } };

    expect(applyDirectoryDraft(source, draft).link?.url).toBe("");
    expect(source.link?.url).toBe("https://example.com/font");
  });
});

describe("Rich card shape", () => {
  it.each([
    ["rows", "minimal", true, "row"],
    ["grid", "minimal", true, "minimal"],
    ["grid", "minimal", false, "minimal"],
    ["gallery", "minimal", true, "minimal"],
    ["gallery", undefined, false, "minimal"],
    ["rows", "fill", true, "row"],
    ["rows", "box", false, "row"],
    ["gallery", "fill", true, "box"],
    ["gallery", "row", false, "box"],
    ["grid", "row", true, "row"],
    ["grid", "row", false, "row"],
    ["grid", "fill", true, "fill"],
    ["grid", "fill", false, "box"],
    ["grid", "box", true, "box"],
    ["grid", undefined, true, "minimal"],
  ] as const)(
    "renders %s/%s with cover=%s as %s",
    (layout, style, hasImage, expected) => {
      expect(getRichCardShape(layout, style, hasImage)).toBe(expected);
    },
  );
});

describe("Rich card image ratio", () => {
  it("keeps Editorial grid covers at 4:3", () => {
    expect(getRichCardImageAspectRatio("grid", "box", image)).toBeCloseTo(
      1.3333333333,
    );
  });

  it("uses natural gallery ratios", () => {
    expect(getRichCardImageAspectRatio("gallery", "box", image)).toBe(0.5);
    expect(
      getRichCardImageAspectRatio("gallery", "box", {
        width: 1200,
        height: 600,
      }),
    ).toBe(2);
  });

  it.each([
    [0, 1200],
    [600, 0],
    [-1, 1200],
    [600, -1],
  ])(
    "uses 4:3 while gallery dimensions %i × %i are invalid",
    (width, height) => {
      expect(
        getRichCardImageAspectRatio("gallery", "box", { width, height }),
      ).toBeCloseTo(1.3333333333);
    },
  );

  it("does not reserve a box ratio for rows, overlays, or missing images", () => {
    expect(getRichCardImageAspectRatio("rows", "row", image)).toBeUndefined();
    expect(getRichCardImageAspectRatio("grid", "fill", image)).toBeUndefined();
    expect(
      getRichCardImageAspectRatio("gallery", "box", undefined),
    ).toBeUndefined();
  });
});

describe("Rich card field visibility", () => {
  it("shows unconfigured fields and respects explicit hiding independently", () => {
    const config = {
      columns: [
        { fid: "fixed:primary_image", hidden: true },
        { fid: "fixed:link", hidden: false },
      ],
    };

    expect(isRichCardFieldVisible(config, "name")).toBe(true);
    expect(isRichCardFieldVisible(config, "description")).toBe(true);
    expect(isRichCardFieldVisible(config, "link")).toBe(true);
    expect(isRichCardFieldVisible(config, "primary_image")).toBe(false);
  });
});

describe("Minimal card field visibility", () => {
  it("defaults to a cover without summary fields", () => {
    const config = { columns: [] };

    expect(isRichCardFieldVisible(config, "name", "minimal")).toBe(false);
    expect(isRichCardFieldVisible(config, "description", "minimal")).toBe(
      false,
    );
    expect(isRichCardFieldVisible(config, "link", "minimal")).toBe(false);
    expect(isRichCardFieldVisible(config, "primary_image", "minimal")).toBe(
      true,
    );
  });

  it.each(["name", "description", "link", "primary_image"] as const)(
    "respects explicit visibility for %s",
    (field) => {
      const fid = `fixed:${field}`;
      const visible = { columns: [{ fid, hidden: false }] };
      const hidden = { columns: [{ fid, hidden: true }] };

      expect(isRichCardFieldVisible(visible, field, "minimal")).toBe(true);
      expect(isRichCardFieldVisible(hidden, field, "minimal")).toBe(false);
    },
  );

  it("uses fixed Grid proportions and natural Gallery sizing", () => {
    expect(getRichCardImageAspectRatio("grid", "minimal", image)).toBeCloseTo(
      4 / 3,
    );
    expect(
      getRichCardImageAspectRatio("gallery", "minimal", image),
    ).toBeUndefined();
    expect(
      getRichCardImageAspectRatio("grid", "minimal", undefined),
    ).toBeUndefined();
  });
});

describe("Directory badge columns", () => {
  const columns: ColumnDefinition[] = [
    {
      fid: "status",
      name: "Status",
      type: "text",
      hidden: false,
      fixed: false,
      placement: "header",
    },
    {
      fid: "empty",
      name: "Empty",
      type: "text",
      hidden: false,
      fixed: false,
      placement: "header",
    },
    {
      fid: "missing",
      name: "Missing",
      type: "text",
      hidden: false,
      fixed: false,
      placement: "header",
    },
    {
      fid: "count",
      name: "Count",
      type: "text",
      hidden: false,
      fixed: false,
      placement: "footer-end",
    },
  ];
  const values: Record<string, ColumnValue> = {
    status: { fid: "status", value: "Open" },
    empty: { fid: "empty", value: "" },
    count: { fid: "count", value: "0" },
  };

  it("omits empty values in view mode and keeps placements separate", () => {
    expect(
      getDirectoryBadgeColumns(columns, values, false, "header").map(
        (column) => column.fid,
      ),
    ).toEqual(["status"]);
    expect(
      getDirectoryBadgeColumns(columns, values, false, "footer-end").map(
        (column) => column.fid,
      ),
    ).toEqual(["count"]);
    expect(
      getDirectoryBadgeColumns(columns, values, false, "footer-start"),
    ).toEqual([]);
  });

  it("keeps empty and missing values editable in their original order", () => {
    expect(
      getDirectoryBadgeColumns(columns, values, true, "header").map(
        (column) => column.fid,
      ),
    ).toEqual(["status", "empty", "missing"]);
  });
});
