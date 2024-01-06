import { ChevronRightIcon } from "@heroicons/react/24/outline";

import {
  DirectoryPath,
  joinDirectoryPath,
} from "src/screens/directory/datagraph/useDirectoryPath";
import { Link } from "src/theme/components/Link";

import { HStack } from "@/styled-system/jsx";

type Props = {
  directoryPath: DirectoryPath;
};

export function Breadcrumbs(props: Props) {
  return (
    <HStack color="fg.subtle">
      <Link href="/directory" size="xs">
        Directory
      </Link>
      {props.directoryPath.map((p) => (
        <>
          <ChevronRightIcon width="1rem" />
          <Link
            key={p}
            href={`/directory/${joinDirectoryPath(props.directoryPath, p)}`}
            size="xs"
          >
            {p}
          </Link>
        </>
      ))}
    </HStack>
  );
}
