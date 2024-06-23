import {
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import { sortableKeyboardCoordinates } from "@dnd-kit/sortable";
import { mutate } from "swr";

import { nodeDelete } from "@/api/openapi/nodes";

export function useDatagraphNodeTree() {
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        delay: 50,
        tolerance: 5,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  );

  async function handleDragEnd(event: any) {
    console.log(event);
  }

  async function handleDelete(slug: string) {
    await nodeDelete(slug);
    await mutate("/api/v1/nodes");
  }

  return {
    handleDelete,
    sensors,
    handleDragEnd,
  };
}
