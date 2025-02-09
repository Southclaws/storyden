import { useRef, useState } from "react";
import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import { handle } from "@/api/client";
import { tagList } from "@/api/openapi-client/tags";
import {
  InstanceCapability,
  Node,
  NodeWithChildren,
  TagNameList,
} from "@/api/openapi-schema";
import { IntelligenceAction } from "@/components/site/Action/Intelligence";
import { TagBadgeList } from "@/components/tag/TagBadgeList";
import { Combotags, CombotagsHandle } from "@/components/ui/combotags";
import { useLibraryMutation } from "@/lib/library/library";
import { useCapability } from "@/lib/settings/capabilities";
import { HStack, styled } from "@/styled-system/jsx";

export type Props<T extends FieldValues> = Omit<
  ControllerProps<T>,
  "render"
> & {
  editing: boolean;
  node: NodeWithChildren;
};

export function LibraryPagePropertyTable<T extends FieldValues>({
  editing,
  node,
  ...props
}: Props<T>) {
  return (
    <Controller
      {...props}
      render={({ field, fieldState, formState }) => {
        //

        return (
          <styled.dl display="table" borderCollapse="collapse">
            {node.properties.map((p) => {
              if (!p.value) {
                return null;
              }

              return (
                <HStack key={p.name} display="table-row">
                  <styled.dt
                    display="table-cell"
                    w="32"
                    p="1"
                    borderRadius="sm"
                    textOverflow="ellipsis"
                    overflowX="hidden"
                    color="fg.muted"
                    _hover={{
                      color: "fg.default",
                      background: "bg.muted",
                      cursor: "pointer",
                    }}
                  >
                    {p.name}
                  </styled.dt>
                  <styled.dd
                    display="table-cell"
                    p="1"
                    w="min"
                    borderRadius="sm"
                    _hover={{
                      color: "fg.default",
                      background: "bg.muted",
                      cursor: "pointer",
                    }}
                  >
                    {p.value}
                  </styled.dd>
                </HStack>
              );
            })}
          </styled.dl>
        );
      }}
    />
  );
}
