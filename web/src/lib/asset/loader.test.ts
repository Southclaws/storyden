import { getImgProps } from "next/dist/shared/lib/get-img-props";
import { imageConfigDefault } from "next/dist/shared/lib/image-config";
import { describe, expect, it, vi } from "vitest";

vi.mock("@/config", () => ({ API_ADDRESS: "https://api.storyden.test" }));

import nextConfig from "../../../next.config.js";

import storydenImageLoader from "./loader";

const asset = "https://api.storyden.test/api/assets/photo.png";

function imageProps(props: Parameters<typeof getImgProps>[0]) {
  return getImgProps(props, {
    defaultLoader: storydenImageLoader,
    imgConf: { ...imageConfigDefault, ...nextConfig.images },
  }).props;
}

describe("storydenImageLoader", () => {
  it.each([
    [1, "s"],
    [320, "s"],
    [321, "m"],
    [960, "m"],
    [961, "l"],
    [3840, "l"],
  ])("maps width %i to the %s preset", (width, variant) => {
    expect(storydenImageLoader({ src: asset, width })).toBe(
      `${asset}?variant=${variant}`,
    );
  });

  it("replaces an existing variant and leaves other sources untouched", () => {
    expect(storydenImageLoader({ src: `${asset}?variant=l`, width: 320 })).toBe(
      `${asset}?variant=s`,
    );
    expect(
      storydenImageLoader({
        src: "https://api.storyden.test/api/accounts/a/avatar",
        width: 64,
      }),
    ).toBe("https://api.storyden.test/api/accounts/a/avatar");
  });
});

describe("next/image with the Storyden loader", () => {
  it("emits one candidate per preset for responsive images", () => {
    const props = imageProps({
      src: asset,
      alt: "",
      fill: true,
      sizes: "8rem",
    });

    expect(props.srcSet).toBe(
      `${asset}?variant=s 320w, ${asset}?variant=m 960w, ${asset}?variant=l 1920w`,
    );
    expect(props.sizes).toBe("8rem");
  });

  it("rounds fixed-width images up to the nearest preset", () => {
    expect(
      imageProps({ src: asset, alt: "", width: 600, height: 300 }).srcSet,
    ).toBe(`${asset}?variant=m 1x, ${asset}?variant=l 2x`);
    expect(
      imageProps({ src: asset, alt: "", width: 400, height: 300 }).srcSet,
    ).toBe(`${asset}?variant=m 1x`);
  });

  it("keeps the source as-is when unoptimized", () => {
    const props = imageProps({
      src: asset,
      alt: "",
      width: 400,
      height: 300,
      unoptimized: true,
    });

    expect(props.src).toBe(asset);
    expect(props.srcSet).toBeUndefined();
  });
});
