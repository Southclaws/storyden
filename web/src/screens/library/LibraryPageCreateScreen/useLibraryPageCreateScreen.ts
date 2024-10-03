import { last } from "lodash";
import { useRouter } from "next/navigation";

import { nodeCreate } from "src/api/openapi-client/nodes";
import {
  Account,
  NodeInitialProps,
  NodeWithChildren,
} from "src/api/openapi-schema";

import { joinLibraryPath } from "../library-path";
import { useLibraryPath } from "../useLibraryPath";

export type Props = {
  session: Account;
};

export function useLibraryPageCreateScreen(props: Props) {
  const router = useRouter();
  const libraryPath = useLibraryPath();

  const initial: NodeWithChildren = {
    id: "",
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    name: "",
    slug: undefined as any, // TODO: Fix the types for this whole screen
    description: "",
    owner: props.session,
    meta: {},
    children: [],
    assets: [],
    visibility: "draft",
    recomentations: [],
  };

  async function handleCreate(node: NodeInitialProps) {
    const parentSlug = last(libraryPath as string[]);
    const created = await nodeCreate({
      name: node.name,
      slug: node.slug,
      url: node.url,
      content: node.content,
      asset_ids: node.asset_ids,
      parent: parentSlug,
    });

    const newPath = joinLibraryPath(libraryPath, created.slug);

    router.push(`/l/${newPath}`);
  }

  return {
    initial,
    handlers: {
      handleCreate,
    },
  };
}
