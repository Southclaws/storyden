import { Node } from "src/api/openapi/schemas";
import { CardGrid, CardRows } from "src/theme/components/Card";

import { RichCardVariantProps } from "@/styled-system/recipes";

import { NodeCard, NodeCardContext } from "./NodeCard";

type Props = {
  directoryPath: string[];
  nodes: Node[];
  size?: RichCardVariantProps["size"];
  context: NodeCardContext;
};

export function NodeCardRows({ directoryPath, nodes, size, context }: Props) {
  return (
    <CardRows>
      {nodes.map((c) => (
        <NodeCard
          key={c.id}
          shape="row"
          size={size}
          context={context}
          directoryPath={directoryPath}
          node={c}
        />
      ))}
    </CardRows>
  );
}

export function NodeCardGrid({ directoryPath, nodes, context }: Props) {
  return (
    <CardGrid>
      {nodes.map((c) => (
        <NodeCard
          key={c.id}
          shape="box"
          context={context}
          directoryPath={directoryPath}
          node={c}
        />
      ))}
    </CardGrid>
  );
}
