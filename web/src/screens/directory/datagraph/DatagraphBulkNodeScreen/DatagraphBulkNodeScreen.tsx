import { Breadcrumbs } from "src/components/directory/datagraph/Breadcrumbs";
import { ClusterCardRows } from "src/components/directory/datagraph/ClusterCardList";
import { DatagraphBulkImport } from "src/components/directory/datagraph/DatagraphBulkImport/DatagraphBulkImport";
import { Props } from "src/components/directory/datagraph/DatagraphBulkImport/useDatagraphBulkImport";

import { useDirectoryPath } from "../useDirectoryPath";

import { LStack } from "@/styled-system/jsx";

export function DatagraphBulkNodeScreen(props: Props) {
  const directoryPath = useDirectoryPath();

  return (
    <LStack>
      <Breadcrumbs directoryPath={directoryPath} create="edit" />

      <DatagraphBulkImport
        node={props.node}
        onCreateNodeFromLink={props.onCreateNodeFromLink}
      />

      <ClusterCardRows
        directoryPath={directoryPath}
        context="directory"
        clusters={props.node?.clusters ?? []}
      />
    </LStack>
  );
}
