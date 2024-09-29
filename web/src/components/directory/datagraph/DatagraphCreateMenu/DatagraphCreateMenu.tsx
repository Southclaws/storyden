import { PlusCircleIcon } from "@heroicons/react/24/outline";
import Link from "next/link";

import { joinDirectoryPath } from "src/screens/directory/datagraph/directory-path";
import { useDirectoryPath } from "src/screens/directory/datagraph/useDirectoryPath";

import { button } from "@/styled-system/recipes";

export function DatagraphCreateMenu() {
  const directoryPath = useDirectoryPath();
  const jointNew = joinDirectoryPath(directoryPath, "new");

  return (
    <Link
      className={button({
        size: "xs",
        variant: "outline",
      })}
      href={`/directory/${jointNew}`}
    >
      <PlusCircleIcon /> Create
    </Link>
  );
}
