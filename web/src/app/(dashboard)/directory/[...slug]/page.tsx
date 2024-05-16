import { notFound, redirect } from "next/navigation";

import { NodeGetOKResponse } from "src/api/openapi/schemas";
import { server } from "src/api/server";
import { getServerSession } from "src/auth/server-session";
import { getTargetSlug } from "src/components/directory/datagraph/utils";
import { NodeCreateManyScreen } from "src/screens/directory/datagraph/NodeCreateManyScreen/NodeCreateManyScreen";
import { NodeCreateScreen } from "src/screens/directory/datagraph/NodeCreateScreen/NodeCreateScreen";
import { NodeViewerScreen } from "src/screens/directory/datagraph/NodeViewerScreen/NodeViewerScreen";
import {
  Params,
  ParamsSchema,
  Query,
  QuerySchema,
} from "src/screens/directory/datagraph/directory-path";

type Props = {
  params: Params;
  searchParams: Query;
};

export default async function Page(props: Props) {
  const { bulk } = QuerySchema.parse(props.searchParams);
  const { slug } = ParamsSchema.parse(props.params);
  const session = await getServerSession();

  const [targetSlug, fallback, isNew] = getTargetSlug(slug);

  const node = await server<NodeGetOKResponse>({
    url: `/v1/nodes/${targetSlug}`,
  });

  if (node) {
    if (isNew) {
      if (!session) {
        redirect(`/login?return=${fallback}`);
      }

      if (bulk) {
        return <NodeCreateManyScreen node={node} />;
      }

      return <NodeCreateScreen session={session} />;
    }

    return <NodeViewerScreen slug={targetSlug} node={node} />;
  }

  // Creating a new item or node from the root: "/directory/new"
  if (isNew) {
    if (!session) {
      redirect(`/login`); // TODO: ?return= back to this path.
    }

    if (bulk) {
      return <NodeCreateManyScreen />;
    }

    return <NodeCreateScreen session={session} />;
  }

  notFound();
}
