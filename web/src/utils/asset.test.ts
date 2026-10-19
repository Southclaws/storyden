import { describe, expect, it, vi } from "vitest";

vi.mock("@/config", () => ({ API_ADDRESS: "https://api.storyden.test" }));

import { getAssetImageProps, getAssetURL, getImageVariantURL } from "./asset";

describe("asset image delivery", () => {
  it("keeps original download URLs and adds presets only when requested", () => {
    expect(getAssetURL("/api/assets/cover.png")).toBe(
      "https://api.storyden.test/api/assets/cover.png",
    );
    expect(getAssetURL("/api/assets/cover.png?token=a", "m")).toBe(
      "https://api.storyden.test/api/assets/cover.png?token=a&variant=m",
    );
  });

  it.each([
    "https://other.test/api/assets/image.png",
    "blob:https://api.storyden.test/preview",
    "data:image/png;base64,abc",
    "http://[invalid",
    "/api/info/banner",
  ])("does not transform other sources: %s", (src) => {
    expect(getImageVariantURL(src, "l")).toBe(src);
  });

  it("uses actual portrait widths for srcset descriptors", () => {
    expect(
      getAssetImageProps({
        path: "/api/assets/portrait.png",
        width: 800,
        height: 1600,
      }).srcSet,
    ).toBe(
      "https://api.storyden.test/api/assets/portrait.png?variant=s 160w, https://api.storyden.test/api/assets/portrait.png?variant=m 480w, https://api.storyden.test/api/assets/portrait.png?variant=l 800w",
    );
  });

  it("deduplicates undersized candidates and tolerates unknown dimensions", () => {
    expect(
      getAssetImageProps({
        path: "/api/assets/small.png",
        width: 100,
        height: 50,
      }).srcSet,
    ).toBe("https://api.storyden.test/api/assets/small.png?variant=s 100w");
    expect(
      getAssetImageProps({ path: "/api/assets/new.png", width: 0, height: 0 }),
    ).toEqual({
      src: "https://api.storyden.test/api/assets/new.png?variant=m",
    });
  });

  it("preserves fragments and resolves protocol-relative API URLs", () => {
    expect(getImageVariantURL("/api/assets/image.png?token=a#image", "m")).toBe(
      "/api/assets/image.png?token=a&variant=m#image",
    );
    expect(getAssetURL("//api.storyden.test/api/assets/image.png", "s")).toBe(
      "https://api.storyden.test/api/assets/image.png?variant=s",
    );
  });

  it.each([
    "https://other.test/image.png",
    "//other.test/image.png",
    "blob:https://api.storyden.test/preview",
    "data:image/png;base64,abc",
  ])("does not create responsive candidates for other sources: %s", (path) => {
    expect(getAssetImageProps({ path, width: 1800, height: 900 })).toEqual({
      src: path,
    });
  });

  it("uses the crop's dimensions instead of its parent's dimensions", () => {
    const original = {
      path: "/api/assets/original.png",
      width: 1800,
      height: 900,
    };
    const crop = {
      path: "/api/assets/crop.png",
      width: 1536,
      height: 384,
      parent: original,
    };

    expect(getAssetImageProps(crop, "l")).toEqual({
      src: "https://api.storyden.test/api/assets/crop.png?variant=l",
      srcSet:
        "https://api.storyden.test/api/assets/crop.png?variant=s 320w, https://api.storyden.test/api/assets/crop.png?variant=m 960w, https://api.storyden.test/api/assets/crop.png?variant=l 1536w",
      width: 1536,
      height: 384,
    });
  });

  it.each([
    [1536, 384, 960, 240],
    [800, 1600, 480, 960],
    [100, 50, 100, 50],
  ])(
    "fits %i x %i into the medium preset without upscaling",
    (width, height, expectedWidth, expectedHeight) => {
      const props = getAssetImageProps(
        { path: "/api/assets/image.png", width, height },
        "m",
      );
      expect(props.width).toBe(expectedWidth);
      expect(props.height).toBe(expectedHeight);
    },
  );
});
