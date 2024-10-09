"use client";

import Link from "next/link";

import { useNodeList } from "@/api/openapi-client/nodes";
import { Visibility } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { LibraryPageTree } from "@/components/library/LibraryPageTree/LibraryPageTree";
import { LStack } from "@/styled-system/jsx";
import { hasPermission } from "@/utils/permissions";

import { AddAction } from "../../Action/Add";
import { Unready } from "../../Unready";
import { NavigationHeader } from "../ContentNavigationList/NavigationHeader";

type Props = {
  label: string;
  href: string;
  currentNode: string | undefined;
  visibility: Visibility[];
};

export function useDatagraphNavTree({ visibility }: Props) {
  const session = useSession();
  const { data, error } = useNodeList({
    visibility,
  });
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

export function DatagraphNavTree(props: Props) {
  const { ready, error, data, canManageLibrary } = useDatagraphNavTree(props);
  if (!ready) {
    return <Unready error={error} />;
  }

  const { label, href, currentNode } = props;

  return (
    <LStack gap="1">
      <NavigationHeader
        href={href}
        controls={
          canManageLibrary && (
            <Link href="/l/new">
              <AddAction size="xs" color="fg.subtle" title="Add a node" />
            </Link>
          )
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
