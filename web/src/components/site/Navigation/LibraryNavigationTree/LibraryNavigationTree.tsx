"use client";

import { CreatePageAction } from "@/components/library/CreatePage";
import { LibraryPageTree } from "@/components/library/LibraryPageTree/LibraryPageTree";
import { LibraryIcon } from "@/components/ui/icons/Library";
import { HStack, LStack } from "@/styled-system/jsx";

import { LibraryLabel, LibraryRoute } from "../Anchors/Library";
import { NavigationHeader } from "../ContentNavigationList/NavigationHeader";

import { Props, useLibraryNavigationTree } from "./useLibraryNavigationTree";

export function LibraryNavigationTree(props: Props) {
  const { ready, data, canManageLibrary } = useLibraryNavigationTree(props);
  if (!ready) {
    // TODO: Render a small version of <Unready /> that's more suitable for this
    return null;
  }

  const { currentNode } = props;

  return (
    <LStack gap="1">
      <NavigationHeader
        href={LibraryRoute}
        controls={
          canManageLibrary && <CreatePageAction variant="ghost" hideLabel />
        }
      >
        <HStack gap="1">
          <LibraryIcon />
          {LibraryLabel}
        </HStack>
      </NavigationHeader>

      <LibraryPageTree
        currentNode={currentNode}
        nodes={data.nodes}
        canManageLibrary={canManageLibrary}
      />
    </LStack>
  );
}
