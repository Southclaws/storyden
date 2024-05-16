import { PropsWithChildren, ReactNode } from "react";

import { Empty } from "src/components/site/Empty";

import { Heading3 } from "../Heading/Index";

import { cx } from "@/styled-system/css";
import { Center, Grid, LStack, styled } from "@/styled-system/jsx";
import { RichCardVariantProps, richCard } from "@/styled-system/recipes";

export type CardItem = {
  id: string;
  title: string;
  url: string;
  text?: string;
  image?: string;
  controls?: React.ReactNode;
};

export type Props = CardItem & RichCardVariantProps;

export function Card({
  children,
  title,
  url,
  text,
  image,
  controls,
  shape,
  size,
}: PropsWithChildren<Props>) {
  const hasImage = Boolean(image);

  const styles = richCard({
    shape,
    size,
    mediaDisplay: hasImage ? "with" : "without",
  });

  return (
    <styled.article className={styles.root}>
      <div className={styles.controlsOverlayContainer}>
        <div className={styles.controls}>{controls}</div>
      </div>

      {image && (
        <div className={styles.mediaBackdropContainer}>
          <styled.img className={styles.mediaBackdrop} src={image} />
        </div>
      )}

      <div className={styles.mediaContainer}>
        {image ? (
          <styled.img className={styles.media} src={image} />
        ) : (
          <div className={styles.mediaMissing}>
            <Center h="full">
              <Empty>no image</Empty>
            </Center>
          </div>
        )}
      </div>

      <div className={styles.contentContainer}>
        <div className={styles.textArea}>
          <Heading3 className={cx("fluid-font-size")}>
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

export type CardGroupProps =
  | {
      items: CardItem[];
      children?: undefined;
    }
  | {
      items?: undefined;
      children: ReactNode[];
    };

export function CardRows(props: CardGroupProps) {
  return (
    <LStack maxH="min">
      {props.children
        ? props.children
        : props.items.map((i) => <Card key={i.id} shape="row" {...i} />)}
    </LStack>
  );
}

export function CardGrid(props: CardGroupProps) {
  const items = props.items?.length ?? props.children?.length ?? 0;

  return (
    <Grid
      w="full"
      gridTemplateColumns={{
        base: "2",
        // Dynamically change the columns based on number of items.
        md: items === 3 ? "3" : items === 4 ? "4" : "2",
      }}
    >
      {props.children
        ? props.children
        : props.items.map((i) => <Card key={i.id} shape="box" {...i} />)}
    </Grid>
  );
}
