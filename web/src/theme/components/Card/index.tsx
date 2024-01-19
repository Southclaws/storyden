import { Empty } from "src/components/site/Empty";

import { Heading3 } from "../Heading/Index";

import { Center, styled } from "@/styled-system/jsx";
import { CardVariantProps, card } from "@/styled-system/recipes";

export type Props = {
  title: string;
  url: string;
  text?: string;
  image?: string;
} & CardVariantProps;

export function Card(allProps: Props) {
  const [styleProps, props] = card.splitVariantProps(allProps);
  const styles = card(styleProps);

  return (
    <styled.article className={styles.root}>
      {props.image && (
        <div className={styles.mediaBackdropContainer}>
          <styled.img className={styles.mediaBackdrop} src={props.image} />
        </div>
      )}

      {props.image ? (
        <div className={styles.mediaContainer}>
          <styled.img className={styles.media} src={props.image} />
        </div>
      ) : (
        <Center display={{ base: "none", md: "flex" }}>
          <Empty>no image</Empty>
        </Center>
      )}

      <div className={styles.textArea}>
        <Heading3 className="fluid-font-size" lineClamp={2}>
          <a href={props.url}>{props.title}</a>
        </Heading3>

        {props.text && <p className={styles.text}>{props.text}</p>}
      </div>
    </styled.article>
  );
}
