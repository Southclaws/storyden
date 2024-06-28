import { last } from "lodash";
import { useRouter } from "next/navigation";

import { nodeCreate } from "src/api/openapi/nodes";
import {
  Account,
  NodeInitialProps,
  NodeWithChildren,
} from "src/api/openapi/schemas";

import { joinDirectoryPath } from "../directory-path";
import { useDirectoryPath } from "../useDirectoryPath";

export type Props = {
  session: Account;
};

export function useNodeCreateScreen(props: Props) {
  const router = useRouter();
  const directoryPath = useDirectoryPath();

  const initial: NodeWithChildren = {
    id: "",
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    name: "",
    slug: undefined as any, // TODO: Fix the types for this whole screen
    description: "",
    owner: props.session,
    properties: {},
    children: [],
    assets: [],
    visibility: "draft",
    recomentations: [],
  };

  async function handleCreate(node: NodeInitialProps) {
    const parentSlug = last(directoryPath as string[]);
    const created = await nodeCreate({
      name: node.name,
      slug: node.slug,
      url: node.url,
      content: node.content,
      asset_ids: node.asset_ids,
      parent: parentSlug,
    });

    const newPath = joinDirectoryPath(directoryPath, created.slug);

    router.push(`/directory/${newPath}`);
  }

  return {
    initial,
    handlers: {
      handleCreate,
    },
  };
}
