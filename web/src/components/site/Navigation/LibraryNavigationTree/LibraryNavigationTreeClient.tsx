"use client";

import { NodeListResult } from "@/api/openapi-schema";

import { LibraryNavigationTree } from "../LibraryNavigationTree/LibraryNavigationTree";
import { useNavigation } from "../useNavigation";

type Props = {
  initialNodeList?: NodeListResult;
};

export function LibraryNavigationTreeClient({ initialNodeList }: Props) {
  const { nodeSlug } = useNavigation();

  return (
    <LibraryNavigationTree
      initialNodeList={initialNodeList}
      currentNode={nodeSlug}
      visibility={["draft", "review", "unlisted", "published"]}
    />
  );
}
