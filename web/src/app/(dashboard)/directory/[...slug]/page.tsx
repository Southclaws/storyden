import { last } from "lodash";
import { notFound } from "next/navigation";

import {
  ClusterGetOKResponse,
  ItemGetOKResponse,
} from "src/api/openapi/schemas";
import { server } from "src/api/server";
import { ClusterScreen } from "src/screens/directory/datagraph/ClusterScreen/ClusterScreen";
import { ItemScreen } from "src/screens/directory/datagraph/ItemScreen/ItemScreen";
import {
  Params,
  ParamsSchema,
} from "src/screens/directory/datagraph/useDirectoryPath";

type Props = {
  params: Params;
};

export default async function Page(props: Props) {
  const { slug } = ParamsSchema.parse(props.params);

  const top = last(slug) ?? "";

  // TODO: here we're firing two requests to the server, one for a cluster and
  // one for the item at the same slug. We should probably have a single request
  // that returns either or a 404. We're also not handling other errors either.
  const [cluster, item] = await Promise.all([
    server<ClusterGetOKResponse>({ url: `/v1/clusters/${top}` }).catch(() => {
      // ignore any errors
    }),
    server<ItemGetOKResponse>({ url: `/v1/items/${top}` }).catch(() => {
      // ignore any errors
    }),
  ]);

  if (cluster) {
    return <ClusterScreen slug={top} cluster={cluster} />;
  }

  if (item) {
    return <ItemScreen slug={top} item={item} />;
  }

  notFound();
}
