"use client";

import { LibraryPageTree } from "@/components/library/LibraryPageTree/LibraryPageTree";
import { LStack } from "@/styled-system/jsx";

import { Unready } from "../../Unready";
import { CreatePageAction } from "../Actions/CreatePage";
import { NavigationHeader } from "../ContentNavigationList/NavigationHeader";

import { Props, useLibraryNavigationTree } from "./useLibraryNavigationTree";

export function LibraryNavigationTree(props: Props) {
  const { ready, error, data, canManageLibrary } =
    useLibraryNavigationTree(props);
  if (!ready) {
    return <Unready error={error} />;
  }

  const { label, href, currentNode } = props;

  return (
    <LStack gap="1">
      <NavigationHeader
        href={href}
        controls={
          canManageLibrary && <CreatePageAction variant="ghost" hideLabel />
        }
      >
        {label}
      </NavigationHeader>

      <LibraryPageTree
        currentNode={currentNode}
        data={{
          label: label,
          children: data.nodes,
        }}
      />
    </LStack>
  );
}
