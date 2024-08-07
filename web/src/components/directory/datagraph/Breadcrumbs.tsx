import { ChevronRightIcon } from "@heroicons/react/24/outline";
import { last } from "lodash";
import { FormEventHandler, ForwardedRef, Fragment, forwardRef } from "react";

import { Visibility } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import {
  DirectoryPath,
  joinDirectoryPath,
} from "src/screens/directory/datagraph/directory-path";

import { Input } from "@/components/ui/input";
import { LinkButton } from "@/components/ui/link-button";
import { Box, HStack } from "@/styled-system/jsx";

import { DatagraphCreateMenu } from "./DatagraphCreateMenu/DatagraphCreateMenu";

type Props = {
  directoryPath: DirectoryPath;
  visibility?: Visibility;
  create: "hide" | "show" | "edit";
  value?: string;
  defaultValue?: string;
  onChange?: FormEventHandler<HTMLInputElement>;
};

export const Breadcrumbs_ = (
  {
    directoryPath,
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
  const paths = directoryPath.filter((p) => p !== "new");
  // const jointNew = joinDirectoryPath(directoryPath, "new");
  const current = last(paths);

  return (
    <HStack w="full" color="fg.subtle" overflowX="scroll" py="2">
      <LinkButton minW="min" href="/directory" size="xs">
        Directory
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
              href={`/directory/${joinDirectoryPath(paths, p)}`}
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
          <DatagraphCreateMenu />
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
