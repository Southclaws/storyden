import { describe, expect, it, vi } from "vitest";

vi.mock("@/config", () => ({ API_ADDRESS: "https://api.storyden.test" }));

import { getDirectoryCardImage, getDirectoryImageProps } from "./media";

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

describe("Directory image props for Surface", () => {
  const image = {
    id: "portrait",
    path: "/api/assets/portrait",
    filename: "portrait.png",
    mime_type: "image/png",
    width: 400,
    height: 1600,
  };

  it("supplies responsive sources and natural sizing for Minimal Gallery", () => {
    const props = getDirectoryCardImage(image, "gallery", "minimal");

    expect(props.image).toBe(
      "https://api.storyden.test/api/assets/portrait?variant=m",
    );
    expect(props.imageAspectRatio).toBeUndefined();
    expect(props.imageFit).toBe("contain");
    expect(props.imageProps).toEqual({
      srcSet:
        "https://api.storyden.test/api/assets/portrait?variant=s 80w, https://api.storyden.test/api/assets/portrait?variant=m 240w, https://api.storyden.test/api/assets/portrait?variant=l 400w",
      sizes: "408px",
      width: 400,
      height: 1600,
      loading: "lazy",
      decoding: "async",
    });
  });

  it("keeps Minimal Grid cropped to 4:3", () => {
    const props = getDirectoryCardImage(image, "grid", "minimal");

    expect(props.imageAspectRatio).toBeCloseTo(4 / 3);
    expect(props.imageFit).toBe("cover");
  });

  it("preserves the existing image path and sizing for Editorial", () => {
    expect(getDirectoryCardImage(image, "gallery", "box")).toEqual({
      image: "https://api.storyden.test/api/assets/portrait?variant=m",
      imageAspectRatio: 0.25,
    });
  });

  it("does not invent an image or ratio when a cover is missing or hidden", () => {
    const props = getDirectoryCardImage(undefined, "grid", "minimal");

    expect(props.image).toBeUndefined();
    expect(props.imageAspectRatio).toBeUndefined();
    expect(props.imageProps?.width).toBeUndefined();
    expect(props.imageProps?.srcSet).toBeUndefined();
  });
});
