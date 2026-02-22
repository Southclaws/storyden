"use client";

import {
  closestCenter,
  DndContext,
  DragEndEvent,
  PointerSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import {
  SortableContext,
  arrayMove,
  rectSortingStrategy,
  useSortable,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import { mutateTransaction } from "@/api/mutate";
import {
  getRoleListKey,
  roleUpdateOrder,
  useRoleList,
} from "@/api/openapi-client/roles";
import {
  Account,
  Permission,
  Role,
  RoleListOKResponse,
} from "@/api/openapi-schema";
import { RoleCard } from "@/components/role/RoleCard";
import { RoleCreateModalTrigger } from "@/components/role/RoleEdit/RoleEditModal";
import { InfoTip } from "@/components/site/InfoTip";
import { UnreadyBanner } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import { IconButton } from "@/components/ui/icon-button";
import { DragHandleIcon } from "@/components/ui/icons/DragHandle";
import { CardGrid } from "@/components/ui/rich-card";
import { isDefaultRole } from "@/lib/role/defaults";
import { HStack, LStack, WStack } from "@/styled-system/jsx";
import { hasPermission } from "@/utils/permissions";

type Props = {
  session?: Account;
  initialRoles: RoleListOKResponse;
};

export function RoleScreen(props: Props) {
  const { data, error } = useRoleList({
    swr: { fallbackData: props.initialRoles },
  });
  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  const canEdit = hasPermission(props.session, Permission.MANAGE_ROLES);

  return (
    <LStack>
      <WStack>
        <Heading>Roles</Heading>

        {canEdit && <RoleCreateModalTrigger />}
      </WStack>

      <HStack gap="1">
        <p>
          Roles provide granular permission control and profile customisation
          for members.
        </p>
        <InfoTip title="Aesthetic roles and badges">
          You can also use Roles as a purely aesthetic tool for providing
          members with ways to express themselves on their profile. Members can
          choose one role as a &ldquo;Badge&rdquo; which is displayed next to
          their name around the site.
        </InfoTip>
      </HStack>

      <SortableRoleGrid roles={data.roles} canEdit={canEdit} />
    </LStack>
  );
}

type SortableRoleGridProps = {
  roles: Role[];
  canEdit: boolean;
};

function SortableRoleGrid({ roles, canEdit }: SortableRoleGridProps) {
  const { mutate } = useSWRConfig();
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: { distance: 3 },
    }),
  );

  const allRoleIDs = roles.map((r) => r.id);
  const customRoleIDs = roles.filter((r) => !isDefaultRole(r)).map((r) => r.id);

  async function handleDragEnd(event: DragEndEvent) {
    if (!canEdit) {
      return;
    }

    const overId = event.over?.id ? String(event.over.id) : null;
    if (!overId) {
      return;
    }

    const activeId = String(event.active.id);
    if (activeId === overId) {
      return;
    }

    const oldIndex = customRoleIDs.indexOf(activeId);
    const newIndex = customRoleIDs.indexOf(overId);
    if (oldIndex === -1 || newIndex === -1) {
      return;
    }

    const customRoles = roles.filter((r) => !isDefaultRole(r));
    const reorderedCustomRoles = arrayMove(customRoles, oldIndex, newIndex);
    const roleIDs = reorderedCustomRoles.map((r) => r.id);
    let customRoleIndex = 0;
    const optimisticRoles = roles.map((r) => {
      if (isDefaultRole(r)) {
        return r;
      }

      return reorderedCustomRoles[customRoleIndex++] ?? r;
    });

    await handle(
      async () => {
        await mutateTransaction(
          mutate,
          [
            {
              key: getRoleListKey(),
              optimistic: (current: RoleListOKResponse | undefined) => {
                if (!current) {
                  return current;
                }

                return {
                  ...current,
                  roles: optimisticRoles,
                };
              },
            },
          ],
          async () => {
            return await roleUpdateOrder({ role_ids: roleIDs });
          },
          { revalidate: false },
        );
      },
      {
        promiseToast: {
          loading: "Updating role order...",
          success: "Role order updated",
        },
      },
    );
  }

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragEnd={handleDragEnd}
    >
      <SortableContext items={allRoleIDs} strategy={rectSortingStrategy}>
        <CardGrid>
          {roles.map((r) => (
            <SortableRoleCard
              key={r.id}
              role={r}
              editable={canEdit}
              draggable={canEdit && !isDefaultRole(r)}
            />
          ))}
        </CardGrid>
      </SortableContext>
    </DndContext>
  );
}

type SortableRoleCardProps = {
  role: Role;
  editable: boolean;
  draggable: boolean;
};

function SortableRoleCard({
  role,
  editable,
  draggable,
}: SortableRoleCardProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: role.id,
    disabled: !draggable,
  });

  const dragStyle = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  return (
    <div ref={setNodeRef} style={dragStyle}>
      <RoleCard
        role={role}
        editable={editable}
        dragHandle={
          draggable ? (
            <IconButton
              variant={{ base: "subtle", md: "ghost" }}
              size="xs"
              minWidth="5"
              width="5"
              height="5"
              padding="0"
              color="fg.muted"
              cursor={isDragging ? "grabbing" : "grab"}
              aria-label={`Reorder role ${role.name}`}
              {...attributes}
              {...listeners}
            >
              <DragHandleIcon width="4" />
            </IconButton>
          ) : undefined
        }
      />
    </div>
  );
}
