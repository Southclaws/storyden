import { ChevronRightIcon } from "@heroicons/react/24/outline";
import { last, pull } from "lodash";
import { FormEventHandler, ForwardedRef, Fragment, forwardRef } from "react";

import { Visibility } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import {
  DirectoryPath,
  joinDirectoryPath,
} from "src/screens/directory/datagraph/useDirectoryPath";
import { Input } from "src/theme/components/Input";
import { Link } from "src/theme/components/Link";

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

export const _Breadcrumbs = (
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
  const paths = pull(directoryPath, "new");
  // const jointNew = joinDirectoryPath(directoryPath, "new");
  const current = last(paths);

  return (
    <HStack w="full" color="fg.subtle" overflowX="scroll" py="2">
      <Link minW="min" href="/directory" size="xs">
        Directory
      </Link>
      {paths.map((p) => {
        const isCurrent = p === current && create === "show";

        return (
          <Fragment key={p}>
            <Box flexShrink="0">
              <ChevronRightIcon width="1rem" />
            </Box>
            <Link
              flexShrink="0"
              borderColor={
                isCurrent && visibility === "published"
                  ? "white"
                  : visibility === "draft"
                    ? "accent"
                    : "blue.500"
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
            </Link>
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

export const Breadcrumbs = forwardRef(_Breadcrumbs);
