"use client";

import { type DateValue, parseDate } from "@internationalized/date";
import { parseAsString, useQueryState } from "nuqs";

import { CancelAction } from "@/components/site/Action/Cancel";
import { Button } from "@/components/ui/button";
import { DateRangePicker } from "@/components/ui/date-picker";
import { IconButton } from "@/components/ui/icon-button";
import { css } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";

export function JoinedDateFilter() {
  const [joined, setJoined] = useQueryState("joined", parseAsString);

  const handleValueChange = async (details: { value: DateValue[] }) => {
    const [start, end] = details.value;

    if (!start && !end) {
      await setJoined(null);
      return;
    }

    const startISO = start ? start.toString() : "";
    const endISO = end ? end.toString() : "";

    const rangeString = `${startISO}/${endISO}`;
    await setJoined(rangeString);
  };

  const handleResetDateRange = async () => {
    await setJoined(null);
  };

  const parseInitialValue = () => {
    if (!joined) return undefined;

    const parts = joined.split("/").filter((p) => p);
    try {
      return parts.map((p) => parseDate(p));
    } catch {
      return undefined;
    }
  };

  return (
    <HStack gap="0">
      <DateRangePicker
        defaultValue={parseInitialValue()}
        onValueChange={handleValueChange}
        active={!!joined}
        hideInputs={true}
        triggerClassName={css({
          borderRightRadius: "none",
        })}
      />
      <CancelAction
        variant="subtle"
        size="sm"
        borderLeftRadius="none"
        onClick={handleResetDateRange}
      />
    </HStack>
  );
}
