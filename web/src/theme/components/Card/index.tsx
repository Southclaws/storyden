import { Empty } from "src/components/site/Empty";

import { Heading3 } from "../Heading/Index";

import { cx } from "@/styled-system/css";
import { Center, Grid, LStack, styled } from "@/styled-system/jsx";
import { CardVariantProps, card } from "@/styled-system/recipes";

export type CardItem = {
  title: string;
  url: string;
  text?: string;
  image?: string;
};

export type Props = CardItem & CardVariantProps;

export function Card({ title, url, text, image, shape }: Props) {
  const hasImage = Boolean(image);

  const styles = card({
    shape,
    mediaDisplay: hasImage ? "with" : "without",
  });

  return (
    <styled.article className={styles.root}>
      {image && <styled.img className={styles.mediaBackdrop} src={image} />}

      {image ? (
        <div className={styles.mediaContainer}>
          <styled.img className={styles.media} src={image} />
        </div>
      ) : (
        <div className={styles.mediaContainer}>
          <Center h="full" display={{ base: "none", sm: "flex" }}>
            <Empty></Empty>
          </Center>
        </div>
      )}

      <div className={styles.textArea}>
        <Heading3 className={cx("fluid-font-size")} lineClamp={2}>
          <a href={url} className={styles.title}>
            {title || "(untitled)"}
          </a>
        </Heading3>

        {text && <p className={styles.text}>{text}</p>}
      </div>
    </styled.article>
  );
}

export function CardRows({ items }: { items: CardItem[] }) {
  return (
    <LStack maxH="min">
      {items.map((i) => (
        <Card key={i.title} shape="row" {...i} />
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
