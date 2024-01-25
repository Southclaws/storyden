import {
  Cluster,
  ClusterInitialProps,
  ClusterWithItems,
  Item,
  ItemInitialProps,
  ItemWithParents,
} from "src/api/openapi/schemas";

export type DatagraphNode =
  | ({
      type: "cluster";
    } & Cluster)
  | ({
      type: "item";
    } & Item);

export type DatagraphNodeWithRelations =
  | ({
      type: "cluster";
    } & ClusterWithItems)
  | ({
      type: "item";
    } & ItemWithParents);

export type DatagraphNodeInitialProps =
  | ({
      type: "cluster";
    } & ClusterInitialProps)
  | ({
      type: "item";
    } & ItemInitialProps);
