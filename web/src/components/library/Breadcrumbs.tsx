import { FormEventHandler, ForwardedRef, Fragment, forwardRef } from "react";

import { NodeReference, Visibility } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { Input } from "@/components/ui/input";
import { LinkButton } from "@/components/ui/link-button";
import { Box, HStack } from "@/styled-system/jsx";

import { BreadcrumbIcon } from "../ui/icons/Breadcrumb";

import { CreatePageAction } from "./CreatePage";

type Props = {
  ancestors: NodeReference[];
  current: NodeReference;
  visibility?: Visibility;
  create: "hide" | "show" | "edit";
  value?: string;
  invalid?: boolean;
  defaultValue?: string;
  onChange?: FormEventHandler<HTMLInputElement>;
};

export const Breadcrumbs_ = (
  {
    ancestors,
    current,
    visibility,
    create,
    value,
    invalid,
    defaultValue,
    onChange,
    ...rest
  }: Props,
  ref: ForwardedRef<HTMLInputElement>,
) => {
  const session = useSession();
  const isEditing = session && create == "edit" && onChange !== undefined;

  // When editing, the slug edit input takes the place of the last breadcrumb.
  const nodes = isEditing ? ancestors : [...ancestors, current];

  return (
    <HStack
      w="full"
      color="text.muted"
      overflowX="scroll"
      pt="scrollGutter"
      mt="-scrollGutter"
    >
      <LinkButton variant="subtle" flexShrink="0" minW="min" href="/l">
        Library
      </LinkButton>
      {nodes.map((node, index) => {
        const isCurrent = node.id === current.id && create === "show";
        const path = nodes
          .slice(0, index + 1)
          .map((part) => part.slug)
          .join("/");

        return (
          <Fragment key={node.id}>
            <Box flexShrink="0">
              <BreadcrumbIcon />
            </Box>
            <LinkButton
              variant="subtle"
              flexShrink="0"
              minW="min"
              colorPalette={
                visibility === "draft"
                  ? "visibility.draft"
                  : visibility === "review"
                    ? "visibility.review"
                    : visibility === "unlisted"
                      ? "visibility.unlisted"
                      : "visibility.published"
              }
              borderColor={
                isCurrent && visibility === "published"
                  ? "white"
                  : isCurrent
                    ? "colorPalette.6"
                    : "border.default"
              }
              borderStyle={
                isCurrent && visibility !== "published" ? "dashed" : "none"
              }
              borderWidth={
                isCurrent && visibility !== "published" ? "thin" : "none"
              }
              href={`/l/${path}`}
            >
              {node.name}{" "}
              {isCurrent && visibility && visibility !== "published" && (
                <span>({visibility})</span>
              )}
            </LinkButton>
          </Fragment>
        );
      })}
      {session && create == "show" && (
        <>
          <Box flexShrink="0">
            <BreadcrumbIcon />
          </Box>
          <CreatePageAction parentSlug={current.slug} />
        </>
      )}
      {isEditing && (
        <>
          <Box flexShrink="0">
            <BreadcrumbIcon />
          </Box>
          <Input
            ref={ref}
            w="full"
            minW="32"
            placeholder="URL slug"
            defaultValue={defaultValue}
            value={value}
            {...(invalid ? { "aria-invalid": "true" } : {})}
            onChange={onChange}
            {...rest}
          />
        </>
      )}
    </HStack>
  );
};

export const Breadcrumbs = forwardRef(Breadcrumbs_);
