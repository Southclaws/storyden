import Link from "next/link";
import {
  Children,
  DOMAttributes,
  PropsWithChildren,
  ReactNode,
  useEffect,
  useRef,
  useState,
} from "react";

import { css } from "@/styled-system/css";
import { Box, LStack, styled } from "@/styled-system/jsx";
import { linkOverlay } from "@/styled-system/patterns";
import {
  RichCardVariantProps,
  cardGrid,
  cardRows,
  richCard,
} from "@/styled-system/recipes";
import { isExternalURL } from "@/utils/url";

import { ContentComposer } from "../../content/ContentComposer/ContentComposer";

export type CardItem = {
  id: string;
  title?: string;
  url: string;
  text?: string;
  titleIcon?: React.ReactNode;
  content?: string;
  image?: string;
  header?: React.ReactNode;
  menu?: React.ReactNode;
  controls?: React.ReactNode;
  disableAnchors?: boolean;
};

export type Props = CardItem & RichCardVariantProps;

export function Card({
  id,
  title,
  url,
  text,
  titleIcon,
  content,
  image,
  header,
  menu,
  controls,
  shape,
  backgroundColor,
  disableAnchors = false,
  children,
}: PropsWithChildren<Props>) {
  const hasImage = Boolean(image);
  const textContainerRef = useRef<HTMLDivElement>(null);
  const [showingMore, setShowingMore] = useState(false);
  const [showMore, setShowMore] = useState(false);

  useEffect(() => {
    if (!textContainerRef.current) return;

    const rect = textContainerRef.current.getBoundingClientRect();

    // 112 = "spacing.28" token * 4
    if (rect.height >= 112) {
      setShowMore(true);
    } else {
      setShowMore(false);
    }
  }, [showingMore, textContainerRef]);

  function handleShowMore() {
    setShowingMore(!showingMore);
  }

  const styles = richCard({
    shape,
    backgroundColor,
  });

  const longContentStyles = css({
    maxHeight: showingMore ? "full" : "28",
    overflow: "hidden",
  });

  const externalURL = isExternalURL(url);

  return (
    <styled.article id={id} className={styles.container}>
      <div className={styles.root}>
        {image && (
          <>
            <div className={styles.mediaBackdropContainer}>
              <styled.img className={styles.mediaBackdrop} src={image} />
            </div>
            <div className={styles.mediaContainer}>
              <styled.img
                className={styles.media}
                src={image}
                maxHeight={showingMore && shape !== "fill" ? "28" : "full"}
              />
            </div>
          </>
        )}

        {header && <div className={styles.headerContainer}>{header}</div>}
        {menu && (
          <div
            className={styles.menuContainer}
            data-header={header ? "present" : "absent"}
          >
            {menu}
          </div>
        )}

        {title && (
          <styled.div className={styles.titleContainer}>
            {titleIcon && (
              <Box className={styles.titleIcon} aria-hidden="true">
                {titleIcon}
              </Box>
            )}

            <styled.h1 className={styles.title}>
              <Link
                className={linkOverlay()}
                href={url}
                onClick={(e) => disableAnchors && e.preventDefault()}
              >
                {title}
              </Link>
            </styled.h1>
          </styled.div>
        )}

        <div className={styles.contentContainer}>
          <div className={styles.textArea}>
            <div ref={textContainerRef} className={longContentStyles}>
              <Link
                href={url}
                className={linkOverlay()}
                onClick={(e) => disableAnchors && e.preventDefault()}
              >
                {text && <p className={styles.text}>{text}</p>}
                {content && (
                  <>
                    <ContentComposer
                      placeholder=""
                      disabled
                      initialValue={content}
                    />
                  </>
                )}
              </Link>
            </div>
            {showMore && (
              <ShowMore showingMore={showingMore} onClick={handleShowMore} />
            )}
          </div>
        </div>

        <div className={styles.footerContainer}>
          {children} {controls}
        </div>
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
      children: ReactNode;
    };

export function CardRows(props: CardGroupProps) {
  return (
    <div className={cardRows()}>
      {props.items
        ? props.items.map((i) => <Card key={i.id} shape="row" {...i} />)
        : props.children}
    </div>
  );
}

export function CardGrid(props: CardGroupProps) {
  const cardCount = props.items
    ? props.items.length
    : Children.toArray(props.children).length;

  const balance =
    cardCount === 1
      ? "one"
      : cardCount % 3 === 1
        ? "last-four"
        : cardCount % 3 === 2
          ? "last-two"
          : "none";

  const styles = cardGrid();

  return (
    <div className={styles.container}>
      <div className={styles.grid} data-balance={balance}>
        {props.items
          ? props.items.map((i) => <Card key={i.id} shape="box" {...i} />)
          : props.children}
      </div>
    </div>
  );
}

function ShowMore({
  showingMore,
  ...props
}: { showingMore: boolean } & DOMAttributes<HTMLAnchorElement>) {
  return (
    <styled.p display="flex" justifyContent="space-between">
      <styled.span color="text.subtle">{showingMore || "..."}</styled.span>
      <styled.a
        fontSize="sm"
        cursor="pointer"
        color="accent.text"
        _hover={{
          textDecoration: "underline",
        }}
        onClick={(e) => e.preventDefault()}
        {...props}
      >
        {showingMore ? "hide" : "show more"}
      </styled.a>
    </styled.p>
  );
}
