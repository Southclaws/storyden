"use client";

import {
  type TabsValueChangeDetails,
  createListCollection,
} from "@ark-ui/react";
import { useEffect, useState } from "react";
import { type UseFormReturn, useWatch } from "react-hook-form";

import { events as eventDefinitions } from "@/api/events";
import { RecurrenceSchedule } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import { ComboboxField } from "@/components/ui/combobox";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { FormLabel } from "@/components/ui/form-label";
import {
  MultiSelectPicker,
  type MultiSelectPickerItem,
} from "@/components/ui/multi-select-picker";
import { SelectField } from "@/components/ui/select";
import * as Tabs from "@/components/ui/tabs";
import { Text } from "@/components/ui/text";
import { ToggleGroupField } from "@/components/ui/toggle-group";
import { describeTrailSchedule, formatOccurrence } from "@/lib/trails";
import { Box, Grid, LStack, styled } from "@/styled-system/jsx";
import { capitalise, humanise } from "@/utils/text";

import { TrailScheduleFields } from "./TrailScheduleFields";
import { TrailFormValues, trailFormSchedule, weekdays } from "./trailForm";

const cadenceCollection = createListCollection({
  items: [
    { label: "One time", value: "once" },
    { label: "Every few hours or days", value: "interval" },
    { label: "Selected weekdays", value: "weekly" },
    { label: "Monthly", value: "monthly" },
    { label: "Yearly", value: "yearly" },
  ],
});

const weekdayItems = weekdays.map((day) => ({
  label: day.slice(0, 3),
  value: day,
  description: capitalise(day),
}));

export function TrailTriggerFields({
  form,
  onPreview,
}: {
  form: UseFormReturn<TrailFormValues>;
  onPreview: (schedule: RecurrenceSchedule) => Promise<string[]>;
}) {
  const triggerType = useWatch({
    control: form.control,
    name: "triggerType",
  });

  function changeTriggerType({ value }: TabsValueChangeDetails) {
    if (value !== "schedule" && value !== "event") return;

    form.setValue("triggerType", value, {
      shouldDirty: true,
      shouldValidate: true,
    });
  }

  return (
    <Tabs.Root
      value={triggerType}
      onValueChange={changeTriggerType}
      variant="line"
      width="full"
    >
      <Tabs.List aria-label="Trigger type">
        <Tabs.Trigger value="schedule">Schedule</Tabs.Trigger>
        <Tabs.Trigger value="event">Event</Tabs.Trigger>
        <Tabs.Indicator />
      </Tabs.List>

      <Tabs.Content value="schedule">
        <ScheduleTriggerFields form={form} onPreview={onPreview} />
      </Tabs.Content>

      <Tabs.Content value="event">
        <EventTriggerFields form={form} />
      </Tabs.Content>
    </Tabs.Root>
  );
}

function ScheduleTriggerFields({
  form,
  onPreview,
}: {
  form: UseFormReturn<TrailFormValues>;
  onPreview: (schedule: RecurrenceSchedule) => Promise<string[]>;
}) {
  const scheduleKind = useWatch({
    control: form.control,
    name: "scheduleKind",
  });
  const [preview, setPreview] = useState<{
    schedule: RecurrenceSchedule;
    occurrences: string[];
  }>();
  const [previewing, setPreviewing] = useState(false);
  const [previewError, setPreviewError] = useState<string>();
  const [timezoneItems] = useState(() =>
    getTimezoneItems(form.getValues("timezone")),
  );

  useEffect(() => {
    const subscription = form.watch((_, { name }) => {
      if (name && isScheduleField(name)) {
        setPreview(undefined);
        setPreviewError(undefined);
      }
    });

    return () => subscription.unsubscribe();
  }, [form]);

  function previewSchedule() {
    setPreviewing(true);
    setPreviewError(undefined);

    void loadSchedulePreview(form.getValues(), onPreview)
      .then(setPreview, (cause: unknown) => {
        setPreviewError(
          errorMessage(cause, "Could not preview this schedule."),
        );
      })
      .finally(() => setPreviewing(false));
  }

  return (
    <LStack gap="4" alignItems="stretch" paddingTop="4">
      <Text variant="supporting">
        Run on a recurring schedule using the selected timezone.
      </Text>

      <Grid
        gap="4"
        gridTemplateColumns={{
          base: "1fr",
          md: "minmax(0, 1fr) minmax(0, 1fr)",
        }}
      >
        <FormControl>
          <FormLabel>Cadence</FormLabel>
          <SelectField
            control={form.control}
            name="scheduleKind"
            collection={cadenceCollection}
            placeholder="Choose a cadence"
            positioning={{ sameWidth: true }}
          />
        </FormControl>

        <FormControl>
          <FormLabel>Timezone</FormLabel>
          <ComboboxField
            control={form.control}
            name="timezone"
            items={timezoneItems}
            placeholder="Search timezones"
            ariaLabel="Timezone"
          />
          <FormHelperText>
            Daylight saving gaps are skipped. Overlaps run once.
          </FormHelperText>
          <FormErrorText>
            {form.formState.errors.timezone?.message}
          </FormErrorText>
        </FormControl>
      </Grid>

      {scheduleKind === "weekly" && (
        <FormControl>
          <FormLabel>Weekdays</FormLabel>
          <ToggleGroupField
            control={form.control}
            name="selectedDays"
            items={weekdayItems}
          />
          <FormErrorText>
            {form.formState.errors.selectedDays?.message}
          </FormErrorText>
        </FormControl>
      )}

      <TrailScheduleFields form={form} scheduleKind={scheduleKind} />

      <Box>
        <Button
          type="button"
          variant="outline"
          loading={previewing}
          loadingText="Checking schedule..."
          onClick={previewSchedule}
        >
          Preview next five
        </Button>
      </Box>

      {previewError && (
        <FormErrorText role="alert">{previewError}</FormErrorText>
      )}

      {preview && (
        <Box backgroundColor="background.inset" borderRadius="md" padding="4">
          <Text fontWeight="medium">
            {describeTrailSchedule(preview.schedule)}
          </Text>
          <styled.ol
            display="grid"
            gap="1"
            listStyle="decimal"
            marginTop="2"
            paddingLeft="5"
          >
            {preview.occurrences.map((value) => (
              <styled.li key={value} color="text.muted" fontSize="sm">
                <styled.time dateTime={value}>
                  {formatOccurrence(value)}
                </styled.time>
              </styled.li>
            ))}
          </styled.ol>
        </Box>
      )}
    </LStack>
  );
}

async function loadSchedulePreview(
  values: TrailFormValues,
  onPreview: (schedule: RecurrenceSchedule) => Promise<string[]>,
) {
  const schedule = trailFormSchedule(values);
  const occurrences = await onPreview(schedule);
  return { schedule, occurrences };
}

function EventTriggerFields({
  form,
}: {
  form: UseFormReturn<TrailFormValues>;
}) {
  const [query, setQuery] = useState("");
  const eventNames = useWatch({
    control: form.control,
    name: "events",
  });
  const selectedEvents = eventNames.map((name) => {
    const event = eventDefinitions.find((candidate) => candidate.name === name);
    return { label: event?.label ?? name, value: name };
  });
  const normalizedQuery = query.trim().toLowerCase();
  const queryResults = eventDefinitions.flatMap((event) => {
    if (
      normalizedQuery &&
      !event.label.toLowerCase().includes(normalizedQuery) &&
      !event.name.toLowerCase().includes(normalizedQuery) &&
      !event.description.toLowerCase().includes(normalizedQuery)
    ) {
      return [];
    }

    return [{ label: event.label, value: event.name }];
  });

  async function changeEvents(items: MultiSelectPickerItem[]) {
    form.setValue(
      "events",
      items.map((item) => item.value),
      { shouldDirty: true, shouldValidate: true },
    );
  }

  return (
    <LStack gap="4" alignItems="stretch" paddingTop="4">
      <Text variant="supporting">When an event is emitted.</Text>

      <FormControl>
        <FormLabel>Events</FormLabel>
        <MultiSelectPicker
          value={selectedEvents}
          onChange={changeEvents}
          onQuery={setQuery}
          queryResults={queryResults}
          inputPlaceholder="Select events..."
          ariaLabel="Select events"
          size="md"
        />
        <FormErrorText>{form.formState.errors.events?.message}</FormErrorText>
      </FormControl>
    </LStack>
  );
}

function getTimezoneItems(current: string) {
  const supported =
    typeof Intl.supportedValuesOf === "function"
      ? Intl.supportedValuesOf("timeZone")
      : [];
  const timezones = new Set([current, "UTC", ...supported]);

  return Array.from(timezones)
    .sort((a, b) => a.localeCompare(b))
    .map((timezone) => ({
      label: humanise(timezone),
      value: timezone,
    }));
}

function isScheduleField(name: string): boolean {
  return [
    "scheduleKind",
    "timezone",
    "startDate",
    "localTime",
    "interval",
    "intervalUnit",
    "selectedDays",
    "month",
    "monthDay",
  ].some((field) => name === field || name.startsWith(`${field}.`));
}

function errorMessage(cause: unknown, fallback: string): string {
  return cause instanceof Error ? cause.message : fallback;
}
