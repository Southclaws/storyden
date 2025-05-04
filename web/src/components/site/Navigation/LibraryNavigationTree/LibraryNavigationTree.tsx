"use client";

import { CreatePageAction } from "@/components/library/CreatePage";
import { LibraryPageTree } from "@/components/library/LibraryPageTree/LibraryPageTree";
import { Unready } from "@/components/site/Unready";
import { LibraryIcon } from "@/components/ui/icons/Library";
import { HStack, LStack } from "@/styled-system/jsx";

import { LibraryLabel, LibraryRoute } from "../Anchors/Library";
import { NavigationHeader } from "../ContentNavigationList/NavigationHeader";

import { Props, useLibraryNavigationTree } from "./useLibraryNavigationTree";

export function LibraryNavigationTree(props: Props) {
  const { ready, error, data, canManageLibrary } =
    useLibraryNavigationTree(props);
  if (!ready) {
    return <Unready error={error} />;
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
