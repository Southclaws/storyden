import { Portal } from "@ark-ui/react";
import { PlusCircleIcon } from "@heroicons/react/24/outline";
import Link from "next/link";

import { joinDirectoryPath } from "src/screens/directory/datagraph/directory-path";
import { useDirectoryPath } from "src/screens/directory/datagraph/useDirectoryPath";
import { Button } from "src/theme/components/Button";
import {
  Menu,
  MenuContent,
  MenuItem,
  MenuItemGroup,
  MenuItemGroupLabel,
  MenuPositioner,
  MenuSeparator,
  MenuTrigger,
} from "src/theme/components/Menu";

export function DatagraphCreateMenu() {
  const directoryPath = useDirectoryPath();
  const jointNew = joinDirectoryPath(directoryPath, "new");

  return (
    <Menu size="sm">
      <MenuTrigger asChild>
        <Button size="xs" variant="outline">
          <PlusCircleIcon /> Create
        </Button>
      </MenuTrigger>
      <Portal>
        <MenuPositioner>
          <MenuContent minW="36">
            <MenuItemGroup id="user">
              <MenuItemGroupLabel
                htmlFor="user"
                display="flex"
                flexDir="column"
                userSelect="none"
              >
                <p>Create a knowledgebase page</p>
              </MenuItemGroupLabel>

              <MenuSeparator />

              <Link href={`/directory/${jointNew}`}>
                <MenuItem id="create-one">
                  <>Create one</>
                </MenuItem>
              </Link>
              <Link href={`/directory/${jointNew}?bulk`}>
                <MenuItem id="create-many">
                  <>Create many</>
                </MenuItem>
              </Link>
            </MenuItemGroup>
          </MenuContent>
        </MenuPositioner>
      </Portal>
    </Menu>
  );
}
