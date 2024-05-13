import { notFound, redirect } from "next/navigation";

import { ClusterGetOKResponse } from "src/api/openapi/schemas";
import { server } from "src/api/server";
import { getServerSession } from "src/auth/server-session";
import { getTargetSlug } from "src/components/directory/datagraph/utils";
import { ClusterCreateManyScreen } from "src/screens/directory/datagraph/ClusterCreateManyScreen/ClusterCreateManyScreen";
import { ClusterCreateScreen } from "src/screens/directory/datagraph/ClusterCreateScreen/ClusterCreateScreen";
import { ClusterViewerScreen } from "src/screens/directory/datagraph/ClusterViewerScreen/ClusterViewerScreen";
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

  const cluster = await server<ClusterGetOKResponse>({
    url: `/v1/clusters/${targetSlug}`,
  });

  if (cluster) {
    if (isNew) {
      if (!session) {
        redirect(`/login?return=${fallback}`);
      }

      if (bulk) {
        return <ClusterCreateManyScreen cluster={cluster} />;
      }

      return <ClusterCreateScreen session={session} />;
    }

    return <ClusterViewerScreen slug={targetSlug} cluster={cluster} />;
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
