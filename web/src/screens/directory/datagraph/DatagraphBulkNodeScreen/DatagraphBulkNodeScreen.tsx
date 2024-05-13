import { Breadcrumbs } from "src/components/directory/datagraph/Breadcrumbs";
import { DatagraphBulkImport } from "src/components/directory/datagraph/DatagraphBulkImport/DatagraphBulkImport";
import { Props } from "src/components/directory/datagraph/DatagraphBulkImport/useDatagraphBulkImport";
import { NodeCardRows } from "src/components/directory/datagraph/NodeCardList";

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

      <NodeCardRows
        directoryPath={directoryPath}
        context="directory"
        nodes={props.node?.children ?? []}
      />
    </LStack>
  );
}
