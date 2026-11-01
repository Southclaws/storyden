import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  Children,
  Fragment,
  PropsWithChildren,
  ReactNode,
  isValidElement,
  useEffect,
  useId,
  useRef,
  useState,
} from "react";

import { isContentEmpty } from "@/lib/content/content";
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

import { CardMedia, CardMediaOptions } from "./RichCardMedia";

export type CardItem = CardMediaOptions & {
  id: string;
  title?: string;
  label?: string;
  titleContent?: ReactNode;
  url: string;
  text?: string;
  titleIcon?: React.ReactNode;
  titleStatus?: string;
  headingAs?: "h2" | "h3" | "h4";
  content?: string;
  body?: ReactNode;
  header?: React.ReactNode;
  menu?: React.ReactNode;
  controls?: React.ReactNode;
  disableAnchors?: boolean;
  prefetchOnIntent?: boolean;
};

export type Props = CardItem & Omit<RichCardVariantProps, "hasMedia">;

export function Card({
  id,
  title,
  label,
  titleContent,
  url,
  text,
  titleIcon,
  titleStatus,
  headingAs = "h2",
  content,
  image,
  imageAlt = "",
  imageAspectRatio,
  imageFit,
  imageProps,
  imageFallback,
  body,
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

  const hasMedia = Boolean(image) || hasSlotContent(imageFallback);
  const styles = richCard({
    shape,
    backgroundColor,
    hasMedia,
  });

  const externalURL = isExternalURL(url);
  const hasRichContent = !isContentEmpty(content);
  const hasPreview = Boolean(text) || hasRichContent;
  const hasBody = hasSlotContent(body);
  const hasHeader = hasSlotContent(header);
  const hasMenu = hasSlotContent(menu);
  const hasTitleContent = hasSlotContent(titleContent);
  const hasFooter = hasSlotContent(children) || hasSlotContent(controls);
  let navigationClass: string | undefined;
  let showBackdrop = Boolean(image);

  if (!disableAnchors && shape !== "minimal") {
    navigationClass = linkOverlay();
  }

  if (shape === "minimal") {
    showBackdrop = false;
  }

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
      {hasRichContent && (
        <ContentComposer placeholder="" disabled initialValue={content} />
      )}
    </>
  );
  const contentLink = (
    <Link
      aria-label={title || label}
      href={url}
      prefetch={false}
      onMouseEnter={handlePrefetch}
      onFocus={handlePrefetch}
      className={navigationClass}
      onClick={(e) => disableAnchors && e.preventDefault()}
    >
      {preview}
    </Link>
  );
  let standaloneContentLink: ReactNode;
  let needsStandaloneLink = !title && !hasTitleContent;

  if (shape === "fill") {
    needsStandaloneLink = true;
  }

  if (hasPreview || shape === "minimal") {
    needsStandaloneLink = false;
  }

  if (needsStandaloneLink) {
    standaloneContentLink = contentLink;
  }

  return (
    <styled.article
      id={id}
      aria-labelledby={title && !hasTitleContent ? titleId : undefined}
      aria-label={label ?? title}
      className={styles.container}
    >
      <div className={styles.root}>
        {showBackdrop && (
          <div className={styles.mediaBackdropContainer}>
            <styled.img
              className={styles.mediaBackdrop}
              src={image}
              alt=""
              aria-hidden="true"
            />
          </div>
        )}
        {hasMedia && (
          <CardMedia
            image={image}
            imageAlt={imageAlt}
            imageAspectRatio={imageAspectRatio}
            imageFit={imageFit}
            imageProps={imageProps}
            imageFallback={imageFallback}
            styles={styles}
            url={url}
            label={label ?? title}
            linked={shape === "minimal" && !disableAnchors}
            collapseImage={showingMore && shape !== "fill"}
          />
        )}

        {hasHeader && <div className={styles.headerContainer}>{header}</div>}
        {hasMenu && (
          <div
            className={styles.menuContainer}
            data-header={hasHeader ? "present" : "absent"}
          >
            {menu}
          </div>
        )}

        {(title || hasTitleContent) && (
          <styled.div className={styles.titleContainer}>
            {hasSlotContent(titleIcon) && (
              <Box className={styles.titleIcon} aria-hidden="true">
                {titleIcon}
              </Box>
            )}

            <styled.h2
              as={headingAs}
              id={titleId}
              className={styles.title}
              flex={hasTitleContent ? "1" : undefined}
              minW="0"
            >
              {hasTitleContent ? (
                titleContent
              ) : (
                <Link
                  data-search-highlight
                  className={navigationClass}
                  href={url}
                  prefetch={false}
                  onMouseEnter={handlePrefetch}
                  onFocus={handlePrefetch}
                  onClick={(e) => disableAnchors && e.preventDefault()}
                >
                  {title}
                </Link>
              )}
              {titleStatus && (
                <>
                  {" "}
                  <styled.span srOnly>{titleStatus}</styled.span>
                </>
              )}
            </styled.h2>
          </styled.div>
        )}

        {(hasPreview || hasBody) && (
          <div className={styles.contentContainer}>
            {hasPreview && (
              <div className={styles.textArea}>
                <ShowMore
                  showingMore={showingMore}
                  title={title}
                  onClick={() => setShowingMore(!showingMore)}
                >
                  {title && shape !== "fill" ? preview : contentLink}
                </ShowMore>
              </div>
            )}
            {hasBody && body}
          </div>
        )}

        {standaloneContentLink}

        {hasFooter && (
          <div className={styles.footerContainer}>
            {children}
            {controls}
          </div>
        )}
      </div>
    </styled.article>
  );
}

function hasSlotContent(value: ReactNode): boolean {
  return Children.toArray(value).some((child) => {
    if (
      isValidElement<{ children?: ReactNode }>(child) &&
      child.type === Fragment
    ) {
      return hasSlotContent(child.props.children);
    }
    return typeof child === "string" ? child.trim().length > 0 : Boolean(child);
  });
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

    if (!element) {
      return;
    }

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
          variant="plain"
          size="sm"
          data-text-link=""
          fontSize="inherit"
          fontWeight="normal"
          lineHeight="inherit"
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
