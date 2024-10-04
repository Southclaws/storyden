"use client";

import Link from "next/link";

import { useNodeList } from "@/api/openapi-client/nodes";
import { useSession } from "@/auth";
import { LibraryPageTree } from "@/components/library/LibraryPageTree/LibraryPageTree";
import { LStack } from "@/styled-system/jsx";
import { hasPermission } from "@/utils/permissions";

import { AddAction } from "../../Action/Add";
import { Unready } from "../../Unready";
import { NavigationHeader } from "../ContentNavigationList/NavigationHeader";

type Props = {
  currentNode: string | undefined;
};

export function useDatagraphNavTree() {
  const session = useSession();
  const { data, error } = useNodeList();
  if (!data) {
    return {
      ready: false as const,
      error,
    };
  }

  const canManageLibrary = hasPermission(session, "MANAGE_LIBRARY");

  return {
    ready: true as const,
    data,
    canManageLibrary,
  };
}

export function DatagraphNavTree({ currentNode }: Props) {
  const { ready, error, data, canManageLibrary } = useDatagraphNavTree();
  if (!ready) {
    return <Unready {...error} />;
  }

  return (
    <LStack gap="1">
      <NavigationHeader
        href="/l"
        controls={
          canManageLibrary && (
            <Link href="/l/new">
              <AddAction size="xs" color="fg.subtle" title="Add a node" />
            </Link>
          )
        }
      >
        Library
      </NavigationHeader>

      <LibraryPageTree
        currentNode={currentNode}
        data={{
          label: "Library",
          children: data.nodes,
        }}
      />
    </LStack>
  );
}
