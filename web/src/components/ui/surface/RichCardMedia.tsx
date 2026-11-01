import Link from "next/link";
import { CSSProperties, ComponentProps, ReactNode } from "react";

import { styled } from "@/styled-system/jsx";
import { richCard } from "@/styled-system/recipes";

export type CardMediaOptions = {
  image?: string;
  imageAlt?: string;
  imageAspectRatio?: number;
  imageFit?: "cover" | "contain";
  imageProps?: Pick<
    ComponentProps<"img">,
    "srcSet" | "sizes" | "width" | "height" | "loading" | "decoding"
  >;
  imageFallback?: ReactNode;
};

type Props = CardMediaOptions & {
  styles: ReturnType<typeof richCard>;
  url: string;
  label?: string;
  linked: boolean;
  collapseImage: boolean;
};

export function CardMedia({
  image,
  imageAlt = "",
  imageAspectRatio,
  imageFit,
  imageProps,
  imageFallback,
  styles,
  url,
  label,
  linked,
  collapseImage,
}: Props) {
  const { width, height, ...responsiveProps } = imageProps ?? {};
  const containerStyle: CSSProperties = {};
  const imageStyle: CSSProperties = { objectFit: imageFit };
  let media: ReactNode = (
    <div className={styles.mediaMissing}>{imageFallback}</div>
  );

  if (image) {
    if (imageAspectRatio) {
      containerStyle.aspectRatio = imageAspectRatio;
      containerStyle.minHeight = 0;
      containerStyle.contain = "none";
    }

    if (collapseImage) {
      imageStyle.maxHeight = "7rem";
    }

    media = (
      <styled.img
        {...responsiveProps}
        htmlWidth={width}
        htmlHeight={height}
        className={styles.media}
        src={image}
        alt={imageAlt}
        style={imageStyle}
      />
    );
  }

  if (linked) {
    return (
      <Link
        href={url}
        aria-label={label}
        prefetch={false}
        className={styles.mediaContainer}
        data-aspect-ratio={Boolean(imageAspectRatio)}
        style={containerStyle}
      >
        {media}
      </Link>
    );
  }

  return (
    <div
      className={styles.mediaContainer}
      data-aspect-ratio={Boolean(imageAspectRatio)}
      style={containerStyle}
    >
      {media}
    </div>
  );
}
