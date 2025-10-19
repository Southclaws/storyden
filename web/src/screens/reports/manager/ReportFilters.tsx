import { SelectValueChangeDetails, createListCollection } from "@ark-ui/react";
import { useMemo } from "react";

import { DatagraphItemKind, ReportStatus } from "@/api/openapi-schema";
import { CheckIcon } from "@/components/ui/icons/Check";
import { SelectIcon } from "@/components/ui/icons/Select";
import * as Select from "@/components/ui/select";
import { DatagraphKindTable } from "@/lib/datagraph/schema";
import { WStack } from "@/styled-system/jsx";

import { useReportsScreenFilters } from "./useReportsScreenFilters";

const REPORT_STATUS_VALUES = Object.values(ReportStatus) as ReportStatus[];
const REPORT_STATUS_LABEL: Record<ReportStatus, string> = {
  [ReportStatus.submitted]: "Submitted",
  [ReportStatus.acknowledged]: "Acknowledged",
  [ReportStatus.resolved]: "Resolved",
};

// TODO: Don't show "post", only thread and reply.
const REPORT_KIND_OPTIONS = Object.entries(DatagraphKindTable).map(
  ([value, label]) => ({
    value: value as DatagraphItemKind,
    label,
  }),
);

export function ReportFilters() {
  const { status, setStatus, kind, setKind } = useReportsScreenFilters();

  const statusCollection = useMemo(
    () =>
      createListCollection({
        items: [
          { label: "All statuses", value: "__all" },
          ...REPORT_STATUS_VALUES.map((value) => ({
            label: REPORT_STATUS_LABEL[value],
            value,
          })),
        ],
      }),
    [],
  );

  const kindCollection = useMemo(
    () =>
      createListCollection({
        items: [{ label: "All items", value: "__all" }, ...REPORT_KIND_OPTIONS],
      }),
    [],
  );

  const statusValue = status;
  const kindValue = kind;

  function handleStatus({ value }: SelectValueChangeDetails) {
    const [selected] = value;
    if (!selected || selected === "__all") {
      setStatus(null);
      return;
    }
    setStatus(selected as ReportStatus);
  }

  function handleKind({ value }: SelectValueChangeDetails) {
    const [selected] = value;
    if (!selected || selected === "__all") {
      setKind(null);
      return;
    }
    setKind(selected as DatagraphItemKind);
  }

  return (
    <WStack gap="2">
      <Select.Root
        size="sm"
        collection={statusCollection}
        value={[statusValue]}
        positioning={{ sameWidth: false }}
        onValueChange={handleStatus}
      >
        <Select.Control>
          <Select.Trigger>
            <Select.ValueText placeholder="Filter by status" />
            <SelectIcon />
          </Select.Trigger>
        </Select.Control>
        <Select.Positioner>
          <Select.Content>
            {statusCollection.items.map((item) => (
              <Select.Item key={item.value} item={item}>
                <Select.ItemText>{item.label}</Select.ItemText>
                <Select.ItemIndicator>
                  <CheckIcon />
                </Select.ItemIndicator>
              </Select.Item>
            ))}
          </Select.Content>
        </Select.Positioner>
      </Select.Root>

      <Select.Root
        size="sm"
        collection={kindCollection}
        value={[kindValue]}
        positioning={{ sameWidth: false }}
        onValueChange={handleKind}
      >
        <Select.Control>
          <Select.Trigger>
            <Select.ValueText placeholder="Filter by item type" />
            <SelectIcon />
          </Select.Trigger>
        </Select.Control>
        <Select.Positioner>
          <Select.Content>
            {kindCollection.items.map((item) => (
              <Select.Item key={item.value} item={item}>
                <Select.ItemText>{item.label}</Select.ItemText>
                <Select.ItemIndicator>
                  <CheckIcon />
                </Select.ItemIndicator>
              </Select.Item>
            ))}
          </Select.Content>
        </Select.Positioner>
      </Select.Root>
    </WStack>
  );
}
