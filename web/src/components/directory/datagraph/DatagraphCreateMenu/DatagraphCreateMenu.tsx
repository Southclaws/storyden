import { Portal } from "@ark-ui/react";
import { PlusCircleIcon } from "@heroicons/react/24/outline";
import Link from "next/link";

import { joinDirectoryPath } from "src/screens/directory/datagraph/directory-path";
import { useDirectoryPath } from "src/screens/directory/datagraph/useDirectoryPath";

import { Button } from "@/components/ui/button";
import * as Menu from "@/components/ui/menu";

export function DatagraphCreateMenu() {
  const directoryPath = useDirectoryPath();
  const jointNew = joinDirectoryPath(directoryPath, "new");

  return (
    <Menu.Root size="sm">
      <Menu.Trigger asChild>
        <Button size="xs" variant="outline">
          <PlusCircleIcon /> Create
        </Button>
      </Menu.Trigger>
      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            <Menu.ItemGroup id="user">
              <Menu.ItemGroupLabel
                htmlFor="user"
                display="flex"
                flexDir="column"
                userSelect="none"
              >
                <p>Create a knowledgebase page</p>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              <Link href={`/directory/${jointNew}`}>
                <Menu.Item id="create-one">
                  <>Create one</>
                </Menu.Item>
              </Link>
              <Link href={`/directory/${jointNew}?bulk`}>
                <Menu.Item id="create-many">
                  <>Create many</>
                </Menu.Item>
              </Link>
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
