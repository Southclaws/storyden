import { createListCollection } from "@ark-ui/react";
import type { UseFormReturn } from "react-hook-form";

import { DatePickerField } from "@/components/ui/date-picker";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormLabel } from "@/components/ui/form-label";
import { Input } from "@/components/ui/input";
import { NumberInputField } from "@/components/ui/number-input";
import { SelectField } from "@/components/ui/select";
import { Grid } from "@/styled-system/jsx";

import { TrailFormValues } from "./trailForm";

const intervalUnitCollection = createListCollection({
  items: [
    { label: "Hours", value: "hours" },
    { label: "Days", value: "days" },
  ],
});

const monthCollection = createListCollection({
  items: Array.from({ length: 12 }, (_, index) => ({
    label: new Intl.DateTimeFormat(undefined, { month: "long" }).format(
      new Date(2000, index, 1),
    ),
    value: String(index + 1),
  })),
});

const monthDayCollection = createListCollection({
  items: [
    { label: "Last day", value: "last" },
    ...Array.from({ length: 31 }, (_, index) => ({
      label: String(index + 1),
      value: String(index + 1),
    })),
  ],
});

export function TrailScheduleFields({
  form,
  scheduleKind,
}: {
  form: UseFormReturn<TrailFormValues>;
  scheduleKind: TrailFormValues["scheduleKind"];
}) {
  return (
    <Grid
      gap="4"
      gridTemplateColumns={{
        base: "1fr",
        sm: "repeat(2, minmax(0, 1fr))",
        lg: "repeat(4, minmax(0, 1fr))",
      }}
    >
      {scheduleKind !== "once" && (
        <FormControl>
          <FormLabel>Every</FormLabel>
          <NumberInputField
            control={form.control}
            name="interval"
            ariaLabel="Interval"
            min={1}
          />
          <FormErrorText>
            {form.formState.errors.interval?.message}
          </FormErrorText>
        </FormControl>
      )}

      {scheduleKind === "interval" && (
        <FormControl>
          <FormLabel>Unit</FormLabel>
          <SelectField
            control={form.control}
            name="intervalUnit"
            collection={intervalUnitCollection}
            placeholder="Choose a unit"
            positioning={{ sameWidth: true }}
          />
        </FormControl>
      )}

      {scheduleKind === "yearly" && (
        <FormControl>
          <FormLabel>Month</FormLabel>
          <SelectField
            control={form.control}
            name="month"
            collection={monthCollection}
            placeholder="Choose a month"
            positioning={{ sameWidth: true }}
          />
        </FormControl>
      )}

      {(scheduleKind === "monthly" || scheduleKind === "yearly") && (
        <FormControl>
          <FormLabel>Day</FormLabel>
          <SelectField
            control={form.control}
            name="monthDay"
            collection={monthDayCollection}
            placeholder="Choose a day"
            positioning={{ sameWidth: true }}
          />
        </FormControl>
      )}

      <FormControl>
        <FormLabel>
          {scheduleKind === "once" ? "Date" : "Starting date"}
        </FormLabel>
        <DatePickerField control={form.control} name="startDate" />
        <FormErrorText>
          {form.formState.errors.startDate?.message}
        </FormErrorText>
      </FormControl>

      <FormControl>
        <FormLabel>Local time</FormLabel>
        <Input
          {...form.register("localTime")}
          aria-label="Local time"
          type="time"
        />
        <FormErrorText>
          {form.formState.errors.localTime?.message}
        </FormErrorText>
      </FormControl>
    </Grid>
  );
}
