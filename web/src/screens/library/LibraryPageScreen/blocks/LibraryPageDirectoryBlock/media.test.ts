import { describe, expect, it, vi } from "vitest";

vi.mock("@/config", () => ({ API_ADDRESS: "https://api.storyden.test" }));

import { getDirectoryImageProps } from "./media";

describe("Directory image presentation", () => {
  it.each([
    [1600, 900, "cover"],
    [600, 900, "cover"],
    [2400, 400, "cover"],
    [400, 1200, "cover"],
    [400, 1201, "contain"],
    [400, 1600, "contain"],
    [1, 1000, "contain"],
    [1000, 1, "cover"],
  ] as const)("Gallery fits %i × %i using %s", (width, height, fit) => {
    const image = { path: "/api/assets/photo", width, height };
    expect(getDirectoryImageProps(image, "gallery").objectFit).toBe(fit);
    expect(getDirectoryImageProps(image, "grid").objectFit).toBe("cover");
  });

  it("uses full-resolution candidates for the narrow stock portrait", () => {
    const props = getDirectoryImageProps(
      { path: "/api/assets/tall", width: 400, height: 1600 },
      "gallery",
    );
    expect(props.srcSet).toBe(
      "https://api.storyden.test/api/assets/tall?variant=s 80w, https://api.storyden.test/api/assets/tall?variant=m 240w, https://api.storyden.test/api/assets/tall?variant=l 400w",
    );
    expect(props.htmlWidth).toBe(400);
    expect(props.htmlHeight).toBe(1600);
  });

  it("omits missing dimensions and responsive candidates while processing", () => {
    const props = getDirectoryImageProps(
      { path: "/api/assets/new", width: 0, height: 0 },
      "gallery",
    );
    expect(props.objectFit).toBe("cover");
    expect(props.htmlWidth).toBeUndefined();
    expect(props.htmlHeight).toBeUndefined();
    expect(props.srcSet).toBeUndefined();
    expect(props.src).toBe(
      "https://api.storyden.test/api/assets/new?variant=m",
    );
  });
});
