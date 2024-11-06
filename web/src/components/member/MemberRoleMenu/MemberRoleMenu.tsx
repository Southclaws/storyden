"use client";

import { Portal } from "@ark-ui/react";
import chroma from "chroma-js";

import { RoleBadge } from "@/components/role/RoleBadge/RoleBadge";
import { badgeColourCSS } from "@/components/role/colours";
import { Unready } from "@/components/site/Unready";
import { CheckCircleIcon } from "@/components/ui/icons/CheckCircle";
import { RemoveCircleIcon } from "@/components/ui/icons/Remove";
import { SubmenuIcon } from "@/components/ui/icons/Submenu";
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
        <SubmenuIcon />
      </Menu.TriggerItem>

      <Portal>
        <Menu.Positioner>
          <Menu.Content>
            {roles.map((r) => {
              const styles = badgeColourCSS(r.colour);

              return (
                <Menu.Item key={r.id} value={r.id} gap="2" style={styles}>
                  {r.selected ? (
                    <CheckCircleIcon
                      width="4"
                      style={{
                        fill: "var(--colors-color-palette)",
                        stroke: "var(--colors-color-palette-text)",
                      }}
                    />
                  ) : (
                    <RemoveCircleIcon width="4" />
                  )}

                  <RoleBadge role={r} />
                </Menu.Item>
              );
            })}
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
