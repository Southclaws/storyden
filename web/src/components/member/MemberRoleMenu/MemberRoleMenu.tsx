"use client";

import { ChevronRightIcon, MinusCircleIcon } from "@heroicons/react/24/outline";
import { CheckCircleIcon } from "@heroicons/react/24/solid";
import chroma from "chroma-js";

import { RoleBadge } from "@/components/role/RoleBadge/RoleBadge";
import { Unready } from "@/components/site/Unready";
import * as Menu from "@/components/ui/menu";
import { HStack } from "@/styled-system/jsx";

import { Props, useMemberRoleMenu } from "./useMemberRoleMenu";

export function MemberRoleMenu(props: Props) {
  const { ready, error, data, handlers } = useMemberRoleMenu(props);
  if (!ready) {
    return <Unready error={error} />;
  }

  const { roles } = data;
  const { handleSelect } = handlers;

  return (
    <Menu.Root
      size="xs"
      lazyMount
      positioning={{ placement: "right-start", gutter: -2 }}
      closeOnSelect={false}
      onSelect={handleSelect}
    >
      <Menu.TriggerItem justifyContent="space-between">
        <HStack gap="2">Roles</HStack>
        <ChevronRightIcon />
      </Menu.TriggerItem>

      <Menu.Positioner>
        <Menu.Content>
          {roles.map((r) => {
            const colour = chroma(r.colour);

            const bgColour = colour.brighten(2).desaturate(1).css();

            return (
              <Menu.Item key={r.id} value={r.id} gap="2">
                {r.selected ? (
                  <CheckCircleIcon width="1rem" fill={bgColour} />
                ) : (
                  <MinusCircleIcon width="1rem" />
                )}

                <RoleBadge role={r} />
              </Menu.Item>
            );
          })}
        </Menu.Content>
      </Menu.Positioner>
    </Menu.Root>
  );
}
