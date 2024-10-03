import { PlusCircleIcon } from "@heroicons/react/24/outline";
import Link from "next/link";

import { joinLibraryPath } from "@/screens/library/library-path";
import { useLibraryPath } from "@/screens/library/useLibraryPath";
import { button } from "@/styled-system/recipes";

export function LibraryPageCreateTrigger() {
  const libraryPath = useLibraryPath();
  const jointNew = joinLibraryPath(libraryPath, "new");

  return (
    <Link
      className={button({
        size: "xs",
        variant: "outline",
      })}
      href={`/l/${jointNew}`}
    >
      <PlusCircleIcon /> Create
    </Link>
  );
}
