import { render } from "@testing-library/react";
import { describe, expect, test, vi } from "vitest";

vi.mock("@/config", () => ({ API_ADDRESS: "https://api.storyden.test" }));

import { AssetImage } from "./AssetImage";

const asset = {
  path: "/api/assets/landscape.png",
  width: 2400,
  height: 1200,
};

function url(variant: string) {
  return `https://api.storyden.test/api/assets/landscape.png?variant=${variant}`;
}

describe("AssetImage", () => {
  test("renders responsive candidates lazily with intrinsic dimensions", () => {
    const { getByRole } = render(
      <AssetImage asset={asset} variant="s" sizes="8rem" alt="Landscape" />,
    );

    const image = getByRole("img");
    expect(image).toHaveAttribute("src", url("s"));
    expect(image).toHaveAttribute(
      "srcset",
      `${url("s")} 320w, ${url("m")} 960w, ${url("l")} 1920w`,
    );
    expect(image).toHaveAttribute("sizes", "8rem");
    expect(image).toHaveAttribute("width", "320");
    expect(image).toHaveAttribute("height", "160");
    expect(image).toHaveAttribute("loading", "lazy");
    expect(image).toHaveAttribute("decoding", "async");
    expect(image).not.toHaveAttribute("fetchpriority");
  });

  test("falls back to a single variant when dimensions are unknown", () => {
    const { getByRole } = render(
      <AssetImage
        asset={{ ...asset, width: 0, height: 0 }}
        sizes="8rem"
        alt=""
        role="img"
      />,
    );

    const image = getByRole("img");
    expect(image).toHaveAttribute("src", url("m"));
    expect(image).not.toHaveAttribute("srcset");
    expect(image).not.toHaveAttribute("sizes");
    expect(image).not.toHaveAttribute("width");
  });

  test("loads priority images eagerly and preloads the candidate set", () => {
    const { getByRole } = render(
      <AssetImage
        asset={asset}
        variant="l"
        sizes="100vw"
        alt="Cover"
        priority
      />,
    );

    const image = getByRole("img");
    expect(image).toHaveAttribute("loading", "eager");
    expect(image).toHaveAttribute("fetchpriority", "high");

    const link = document.head.querySelector('link[rel="preload"][as="image"]');
    expect(link).toHaveAttribute("imagesrcset", image.getAttribute("srcset"));
    expect(link).toHaveAttribute("imagesizes", "100vw");
    expect(link).toHaveAttribute("fetchpriority", "high");
  });
});
