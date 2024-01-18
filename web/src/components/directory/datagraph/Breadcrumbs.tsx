import { ChevronRightIcon, PlusCircleIcon } from "@heroicons/react/24/outline";
import { pull } from "lodash";
import { FormEventHandler, ForwardedRef, Fragment, forwardRef } from "react";

import {
  DirectoryPath,
  joinDirectoryPath,
} from "src/screens/directory/datagraph/useDirectoryPath";
import { Input } from "src/theme/components/Input";
import { Link } from "src/theme/components/Link";

import { HStack } from "@/styled-system/jsx";

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
  const isEditing = create == "edit" && onChange !== undefined;
  const editingPath = isEditing ? directoryPath.slice(0, -1) : directoryPath;
  const paths = pull(editingPath, "new");
  const jointNew = joinDirectoryPath(directoryPath, "new");

  return (
    <HStack w="full" color="fg.subtle">
      <Link minW="min" href="/directory" size="xs">
        Directory
      </Link>
      {paths.map((p) => (
        <Fragment key={p}>
          <ChevronRightIcon width="1rem" />
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
      {create == "show" && (
        <>
          <ChevronRightIcon width="1rem" />
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
          <ChevronRightIcon width="1rem" />
          <Input
            ref={ref}
            w="full"
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
