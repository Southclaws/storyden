import { ChevronRightIcon, PlusCircleIcon } from "@heroicons/react/24/outline";
import { pull } from "lodash";
import { FormEventHandler, ForwardedRef, Fragment, forwardRef } from "react";

import { useSession } from "src/auth";
import {
  DirectoryPath,
  joinDirectoryPath,
} from "src/screens/directory/datagraph/useDirectoryPath";
import { Input } from "src/theme/components/Input";
import { Link } from "src/theme/components/Link";

import { Box, HStack } from "@/styled-system/jsx";

type Props = {
  directoryPath: DirectoryPath;
  create: "hide" | "show" | "edit";
  value?: string;
  defaultValue?: string;
  onChange?: FormEventHandler<HTMLInputElement>;
};

export const _Breadcrumbs = (
  { directoryPath, create, value, defaultValue, onChange, ...rest }: Props,
  ref: ForwardedRef<HTMLInputElement>,
) => {
  const session = useSession();
  const isEditing = session && create == "edit" && onChange !== undefined;
  const paths = pull(directoryPath, "new");
  const jointNew = joinDirectoryPath(directoryPath, "new");

  return (
    <HStack w="full" color="fg.subtle" overflowX="scroll" py="2">
      <Link minW="min" href="/directory" size="xs">
        Directory
      </Link>
      {paths.map((p) => (
        <Fragment key={p}>
          <Box flexShrink="0">
            <ChevronRightIcon width="1rem" />
          </Box>
          <Link
            flexShrink="0"
            key={p}
            href={`/directory/${joinDirectoryPath(paths, p)}`}
            size="xs"
          >
            {p}
          </Link>
        </Fragment>
      ))}
      {session && create == "show" && (
        <>
          <Box flexShrink="0">
            <ChevronRightIcon width="1rem" />
          </Box>
          <Link
            flexShrink="0"
            kind="primary"
            href={`/directory/${jointNew}`}
            size="xs"
          >
            <PlusCircleIcon /> Create
          </Link>
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
