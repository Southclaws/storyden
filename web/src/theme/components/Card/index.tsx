import { PropsWithChildren } from "react";

import { Heading3 } from "../Heading/Index";

import { cx } from "@/styled-system/css";
import { Grid, LStack, styled } from "@/styled-system/jsx";
import { CardVariantProps, card } from "@/styled-system/recipes";

export type CardItem = {
  id: string;
  title: string;
  url: string;
  text?: string;
  image?: string;
  controls?: React.ReactNode;
};

export type Props = CardItem & CardVariantProps;

export function Card({
  children,
  title,
  url,
  text,
  image,
  controls,
  shape,
}: PropsWithChildren<Props>) {
  const hasImage = Boolean(image);

  const styles = card({
    shape,
    mediaDisplay: hasImage ? "with" : "without",
  });

  return (
    <styled.article className={styles.root}>
      <div className={styles.childrenOverlay}>{controls}</div>

      {image && <styled.img className={styles.mediaBackdrop} src={image} />}

      {image && (
        <div className={styles.mediaContainer}>
          <styled.img className={styles.media} src={image} />
        </div>
      )}

      <div className={styles.contentContainer}>
        <div className={styles.textArea}>
          <Heading3 className={cx("fluid-font-size")} lineClamp={2}>
            <a href={url} className={styles.title}>
              {title || "(untitled)"}
            </a>
          </Heading3>

          {text && <p className={styles.text}>{text}</p>}
        </div>

        <div className={styles.footer}>{children}</div>
      </div>
    </styled.article>
  );
}

export function CardRows({ items }: { items: CardItem[] }) {
  return (
    <LStack maxH="min">
      {items.map((i) => (
        <Card key={i.id} shape="row" {...i} />
      ))}
    </LStack>
  );
}

export function CardGrid({ items }: { items: CardItem[] }) {
  return (
    <Grid
      w="full"
      gridTemplateColumns={{
        base: "2",
        sm: "4",
        lg: "6",
      }}
    >
      {items.map((i) => (
        <Card key={i.title} shape="row" {...i} />
      ))}
    </Grid>
  );
}
