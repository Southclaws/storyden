import { notFound, redirect } from "next/navigation";

import {
  ClusterGetOKResponse,
  ItemGetOKResponse,
} from "src/api/openapi/schemas";
import { server } from "src/api/server";
import { useServerSession } from "src/auth/server-session";
import { getTargetSlug } from "src/components/directory/datagraph/utils";
import { ClusterCreateManyScreen } from "src/screens/directory/datagraph/ClusterCreateManyScreen/ClusterCreateManyScreen";
import { ClusterCreateScreen } from "src/screens/directory/datagraph/ClusterCreateScreen/ClusterCreateScreen";
import { ClusterViewerScreen } from "src/screens/directory/datagraph/ClusterViewerScreen/ClusterViewerScreen";
import { ItemViewerScreen } from "src/screens/directory/datagraph/ItemViewerScreen/ItemViewerScreen";
import {
  Params,
  ParamsSchema,
  Query,
  QuerySchema,
} from "src/screens/directory/datagraph/useDirectoryPath";

type Props = {
  params: Params;
  searchParams: Query;
};

export default async function Page(props: Props) {
  const { bulk } = QuerySchema.parse(props.searchParams);
  const { slug } = ParamsSchema.parse(props.params);
  const session = await useServerSession();

  const [targetSlug, fallback, isNew] = getTargetSlug(slug);

  // TODO: here we're firing two requests to the server, one for a cluster and
  // one for the item at the same slug. We should probably have a single request
  // that returns either or a 404. We're also not handling other errors either.
  const [cluster, item] = await Promise.all([
    server<ClusterGetOKResponse>({ url: `/v1/clusters/${targetSlug}` }).catch(
      () => {
        // ignore any errors
      },
    ),
    server<ItemGetOKResponse>({ url: `/v1/items/${targetSlug}` }).catch(() => {
      // ignore any errors
    }),
  ]);

  if (cluster) {
    if (isNew) {
      if (!session) {
        redirect(`/login`); // TODO: ?return= back to this path.
      }

      if (bulk) {
        return <ClusterCreateManyScreen cluster={cluster} />;
      }

      return <ClusterCreateScreen session={session} />;
    }

    return <ClusterViewerScreen slug={targetSlug} cluster={cluster} />;
  }

  if (item) {
    if (isNew) {
      redirect(`/directory/${fallback}`);
    }
    return <ItemViewerScreen slug={targetSlug} item={item} />;
  }

  // Creating a new item or cluster from the root: "/directory/new"
  if (isNew) {
    if (!session) {
      redirect(`/login`); // TODO: ?return= back to this path.
    }

    if (bulk) {
      return <ClusterCreateManyScreen />;
    }

    return <ClusterCreateScreen session={session} />;
  }

  notFound();
}
