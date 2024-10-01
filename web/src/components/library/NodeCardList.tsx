import { Node } from "src/api/openapi-schema";

import { CardGrid, CardRows } from "@/components/ui/rich-card";
import { RichCardVariantProps } from "@/styled-system/recipes";

import { NodeCard, NodeCardContext } from "./NodeCard";

type Props = {
  libraryPath: string[];
  nodes: Node[];
  size?: RichCardVariantProps["size"];
  context: NodeCardContext;
};

export function NodeCardRows({ libraryPath, nodes, size, context }: Props) {
  return (
    <CardRows>
      {nodes.map((c) => (
        <NodeCard
          key={c.id}
          shape="row"
          size={size}
          context={context}
          libraryPath={libraryPath}
          node={c}
        />
      ))}
    </CardRows>
  );
}

export function NodeCardGrid({ libraryPath, nodes, context }: Props) {
  return (
    <CardGrid>
      {nodes.map((c) => (
        <NodeCard
          key={c.id}
          shape="box"
          context={context}
          libraryPath={libraryPath}
          node={c}
        />
      ))}
    </CardGrid>
  );
}
