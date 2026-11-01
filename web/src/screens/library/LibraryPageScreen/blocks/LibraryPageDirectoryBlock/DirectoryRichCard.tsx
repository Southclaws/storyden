import { ReactNode } from "react";

import { Card } from "@/components/ui/surface";

import { DirectoryFieldEditor } from "./DirectoryCardFields";
import { DirectoryBadgeFooter } from "./DirectoryPropertyBadges";
import {
  DirectoryRichCardBody,
  DirectoryRichCardHeader,
  DirectoryRichCardImageFallback,
} from "./DirectoryRichCardContent";
import {
  DirectoryRichCardProps,
  getDirectoryBadgeColumns,
  getRichCardShape,
  isRichCardFieldVisible,
} from "./directoryCard";
import { getDirectoryCardImage } from "./media";

export function DirectoryRichCard(props: DirectoryRichCardProps) {
  const { node, config, layout, columns, values, editing, onFieldValueChange } =
    props;
  const showCover = isRichCardFieldVisible(config, "primary_image");

  let image = node.primary_image;

  if (!showCover) {
    image = undefined;
  }

  const shape = getRichCardShape(layout, config.card_style, Boolean(image));
  const showTitle = isRichCardFieldVisible(config, "name", shape);
  const showDescription = isRichCardFieldVisible(config, "description", shape);
  const showLink = isRichCardFieldVisible(config, "link", shape);
  const imageProps = getDirectoryCardImage(image, layout, shape);
  const url = `/l/${node.slug}`;
  let imageFallback: ReactNode;

  if (shape === "minimal" && showCover) {
    imageFallback = (
      <DirectoryRichCardImageFallback node={node} showTitle={showTitle} />
    );
  }

  const headerBadges = getDirectoryBadgeColumns(
    columns,
    values,
    editing,
    "header",
  );
  const footerStart = getDirectoryBadgeColumns(
    columns,
    values,
    editing,
    "footer-start",
  );
  const footerEnd = getDirectoryBadgeColumns(
    columns,
    values,
    editing,
    "footer-end",
  );
  const hasFooter = footerStart.length + footerEnd.length > 0;
  const hasDetails = columns.some(
    (column) => (column.placement ?? "details") === "details",
  );

  let hasHeader = headerBadges.length > 0;

  if (shape === "fill" && showTitle) {
    hasHeader = true;
  }

  let hasDescription = false;

  if (showDescription) {
    hasDescription = editing || Boolean(node.description);
  }

  let hasLink = false;

  if (showLink) {
    hasLink = editing || Boolean(node.link);
  }

  let hasBody = hasDetails;

  if (hasDescription) {
    hasBody = true;
  }

  if (hasLink) {
    hasBody = true;
  }

  let content: string | undefined;

  if (config.show_content) {
    content = node.content;
  }

  let title: string | undefined;

  if (showTitle) {
    title = node.name;
  }

  let titleEditor: ReactNode;

  if (editing) {
    if (showTitle && shape !== "fill") {
      titleEditor = (
        <DirectoryFieldEditor
          node={node}
          onFieldValueChange={onFieldValueChange}
          fid="fixed:name"
          label="Name"
          value={node.name}
          elevated
        />
      );
    }
  }

  return (
    <Card
      id={`directory-card-${node.id}`}
      title={title}
      label={node.name || "Untitled page"}
      titleContent={titleEditor}
      headingAs="h3"
      url={url}
      disableAnchors={editing}
      shape={shape}
      content={content}
      {...imageProps}
      imageFallback={imageFallback}
      header={
        hasHeader && (
          <DirectoryRichCardHeader
            {...props}
            shape={shape}
            showTitle={showTitle}
            url={url}
          />
        )
      }
      body={
        hasBody && (
          <DirectoryRichCardBody
            {...props}
            shape={shape}
            showDescription={hasDescription}
            showLink={hasLink}
          />
        )
      }
    >
      {hasFooter && <DirectoryBadgeFooter {...props} />}
    </Card>
  );
}
