import { last } from "lodash";
import { uniq } from "lodash/fp";
import { FormEventHandler, ForwardedRef, Fragment, forwardRef } from "react";

import { Visibility } from "src/api/openapi-schema";
import { useSession } from "src/auth";

import { Input } from "@/components/ui/input";
import { LinkButton } from "@/components/ui/link-button";
import { LibraryPath, joinLibraryPath } from "@/screens/library/library-path";
import { Box, HStack } from "@/styled-system/jsx";

import { BreadcrumbIcon } from "../ui/icons/Breadcrumb";

import { CreatePageAction } from "./CreatePage";

type Props = {
  libraryPath: LibraryPath;
  visibility?: Visibility;
  create: "hide" | "show" | "edit";
  value?: string;
  invalid?: boolean;
  defaultValue?: string;
  onChange?: FormEventHandler<HTMLInputElement>;
};

export const Breadcrumbs_ = (
  {
    libraryPath,
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

  // Sometimes, due to bugs, the path can contain duplicate slug entries.
  const uniquePaths = uniq(libraryPath);

  // When editing, the slug edit input takes the place of the last breadcrumb.
  const paths = isEditing ? uniquePaths.slice(0, -1) : uniquePaths;
  const current = last(paths);

  return (
    <HStack
      w="full"
      color="fg.subtle"
      overflowX="scroll"
      pt="scrollGutter"
      mt="-scrollGutter"
    >
      <LinkButton
        size="xs"
        variant="subtle"
        flexShrink="0"
        minW="min"
        href="/l"
      >
        Library
      </LinkButton>
      {paths.map((p) => {
        const isCurrent = p === current && create === "show";

        return (
          <Fragment key={p}>
            <Box flexShrink="0">
              <BreadcrumbIcon />
            </Box>
            <LinkButton
              size="xs"
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
              key={p}
              href={`/l/${joinLibraryPath(paths, p)}`}
            >
              {p}{" "}
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
          <CreatePageAction parentSlug={current} />
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
            size="xs"
            height="6" // TODO: Make this default for size="xs"
            borderRadius="sm"
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
