import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  Children,
  PropsWithChildren,
  ReactNode,
  useEffect,
  useId,
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
import { Button } from "../button";

export type CardItem = {
  id: string;
  title?: string;
  url: string;
  text?: string;
  titleIcon?: React.ReactNode;
  titleStatus?: string;
  headingAs?: "h2" | "h3" | "h4";
  content?: string;
  image?: string;
  imageAlt?: string;
  header?: React.ReactNode;
  menu?: React.ReactNode;
  controls?: React.ReactNode;
  disableAnchors?: boolean;
  prefetchOnIntent?: boolean;
};

export type Props = CardItem & RichCardVariantProps;

export function Card({
  id,
  title,
  url,
  text,
  titleIcon,
  titleStatus,
  headingAs = "h2",
  content,
  image,
  imageAlt = "",
  header,
  menu,
  controls,
  shape,
  backgroundColor,
  disableAnchors = false,
  prefetchOnIntent = false,
  children,
}: PropsWithChildren<Props>) {
  const router = useRouter();
  const titleId = useId();
  const prefetchedURL = useRef<string | null>(null);
  const [showingMore, setShowingMore] = useState(false);

  const styles = richCard({
    shape,
    backgroundColor,
  });

  const externalURL = isExternalURL(url);

  function handlePrefetch() {
    if (
      !prefetchOnIntent ||
      disableAnchors ||
      externalURL ||
      prefetchedURL.current === url
    ) {
      return;
    }

    prefetchedURL.current = url;
    router.prefetch(url);
  }

  const preview = (
    <>
      {text && (
        <p className={styles.text} data-search-highlight>
          {text}
        </p>
      )}
      {content && (
        <ContentComposer placeholder="" disabled initialValue={content} />
      )}
    </>
  );

  return (
    <styled.article
      id={id}
      aria-labelledby={title ? titleId : undefined}
      className={styles.container}
    >
      <div className={styles.root}>
        {image && (
          <>
            <div className={styles.mediaBackdropContainer}>
              <styled.img
                className={styles.mediaBackdrop}
                src={image}
                alt=""
                aria-hidden="true"
              />
            </div>
            <div className={styles.mediaContainer}>
              <styled.img
                className={styles.media}
                src={image}
                alt={imageAlt}
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

            <styled.h2 as={headingAs} id={titleId} className={styles.title}>
              <Link
                data-search-highlight
                className={linkOverlay()}
                href={url}
                prefetch={false}
                onMouseEnter={handlePrefetch}
                onFocus={handlePrefetch}
                onClick={(e) => disableAnchors && e.preventDefault()}
              >
                {title}
              </Link>
              {titleStatus && (
                <>
                  {" "}
                  <styled.span srOnly>{titleStatus}</styled.span>
                </>
              )}
            </styled.h2>
          </styled.div>
        )}

        <div className={styles.contentContainer}>
          <div className={styles.textArea}>
            <ShowMore
              showingMore={showingMore}
              title={title}
              onClick={() => setShowingMore(!showingMore)}
            >
              {title && shape !== "fill" ? (
                preview
              ) : (
                <Link
                  aria-label={title}
                  href={url}
                  prefetch={false}
                  onMouseEnter={handlePrefetch}
                  onFocus={handlePrefetch}
                  className={linkOverlay()}
                  onClick={(e) => disableAnchors && e.preventDefault()}
                >
                  {preview}
                </Link>
              )}
            </ShowMore>
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
  title,
  onClick,
  children,
}: PropsWithChildren<{
  showingMore: boolean;
  title?: string;
  onClick: () => void;
}>) {
  const contentId = useId();
  const textContainerRef = useRef<HTMLDivElement>(null);
  const [showMore, setShowMore] = useState(false);

  useEffect(() => {
    const element = textContainerRef.current;
    if (!element) return;
    const measure = () => setShowMore(element.scrollHeight > 112);
    measure();
    const observer = new ResizeObserver(measure);
    observer.observe(element);
    return () => observer.disconnect();
  }, [children]);

  const label = showingMore ? "Show less" : "Show more";
  return (
    <>
      <div
        id={contentId}
        ref={textContainerRef}
        className={css({
          maxHeight: showingMore ? "full" : "28",
          overflow: "hidden",
        })}
      >
        {children}
      </div>
      {showMore && (
        <Button
          type="button"
          variant="subtle"
          size="sm"
          position="relative"
          zIndex="var(--card-content-index)"
          aria-expanded={showingMore}
          aria-controls={contentId}
          aria-label={title ? `${label} of ${title}` : label}
          onClick={onClick}
        >
          {label}
        </Button>
      )}
    </>
  );
}
