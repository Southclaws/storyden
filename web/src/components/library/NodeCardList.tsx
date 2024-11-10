import { Node } from "src/api/openapi-schema";

import { CardGrid, CardRows } from "@/components/ui/rich-card";
import { RichCardVariantProps } from "@/styled-system/recipes";

import { NodeCard, NodeCardContext } from "./NodeCard";

type Props = {
  libraryPath: string[];
  nodes: Node[];
  context: NodeCardContext;
};

export function NodeCardRows({ libraryPath, nodes, context }: Props) {
  return (
    <CardRows>
      {nodes.map((c) => (
        <NodeCard key={c.id} shape="row" libraryPath={libraryPath} node={c} />
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
          shape="responsive"
          libraryPath={libraryPath}
          node={c}
        />
      ))}
    </CardGrid>
  );
}
