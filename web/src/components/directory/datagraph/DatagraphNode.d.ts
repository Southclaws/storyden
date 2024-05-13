import {
  Cluster,
  ClusterInitialProps,
  ClusterWithItems,
} from "src/api/openapi/schemas";

export type DatagraphNode = Cluster;

export type DatagraphNodeWithRelations = ClusterWithItems;

export type DatagraphNodeInitialProps = ClusterInitialProps;
