import { ChevronRightIcon } from "@heroicons/react/24/outline";
import { last } from "lodash";
import { FormEventHandler, ForwardedRef, Fragment, forwardRef } from "react";

import { Visibility } from "src/api/openapi-schema";
import { useSession } from "src/auth";

import { Input } from "@/components/ui/input";
import { LinkButton } from "@/components/ui/link-button";
import { LibraryPath, joinLibraryPath } from "@/screens/library/library-path";
import { Box, HStack } from "@/styled-system/jsx";

import { LibraryPageCreateTrigger } from "./LibraryPageCreateTrigger/LibraryPageCreateTrigger";

type Props = {
  libraryPath: LibraryPath;
  visibility?: Visibility;
  create: "hide" | "show" | "edit";
  value?: string;
  defaultValue?: string;
  onChange?: FormEventHandler<HTMLInputElement>;
};

export const Breadcrumbs_ = (
  {
    libraryPath,
    visibility,
    create,
    value,
    defaultValue,
    onChange,
    ...rest
  }: Props,
  ref: ForwardedRef<HTMLInputElement>,
) => {
  const session = useSession();
  const isEditing = session && create == "edit" && onChange !== undefined;
  const paths = libraryPath.filter((p) => p !== "new");
  const current = last(paths);

  return (
    <HStack w="full" color="fg.subtle" overflowX="scroll" py="2">
      <LinkButton minW="min" href="/l" size="xs">
        Library
      </LinkButton>
      {paths.map((p) => {
        const isCurrent = p === current && create === "show";

        return (
          <Fragment key={p}>
            <Box flexShrink="0">
              <ChevronRightIcon width="1rem" />
            </Box>
            <LinkButton
              flexShrink="0"
              borderColor={
                isCurrent && visibility === "published"
                  ? "white"
                  : visibility === "draft"
                    ? "accent"
                    : "blue.8"
              }
              borderStyle={
                isCurrent && visibility !== "published" ? "dashed" : "none"
              }
              borderWidth={
                isCurrent && visibility !== "published" ? "thin" : "none"
              }
              key={p}
              href={`/l/${joinLibraryPath(paths, p)}`}
              size="xs"
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
            <ChevronRightIcon width="1rem" />
          </Box>
          <LibraryPageCreateTrigger />
        </>
      )}
      {isEditing && (
        <>
          <Box flexShrink="0">
            <ChevronRightIcon width="1rem" />
          </Box>
          <Input
            ref={ref}
            w="full"
            minW="32"
            size="xs"
            placeholder="URL slug"
            defaultValue={defaultValue}
            value={value}
            onChange={onChange}
            {...rest}
          />
        </>
      )}
    </HStack>
  );
};

export const Breadcrumbs = forwardRef(Breadcrumbs_);
