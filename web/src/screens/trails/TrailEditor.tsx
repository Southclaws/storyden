"use client";

import { createListCollection } from "@ark-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import { useFieldArray, useForm } from "react-hook-form";
import { useSWRConfig } from "swr";

import { useRobotsList } from "@/api/openapi-client/robots";
import {
  getTrailListKey,
  trailCreate,
  trailSchedulePreview,
  trailUpdate,
} from "@/api/openapi-client/trails";
import {
  RecurrenceSchedule,
  Trail,
  TrailMutableProps,
} from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import { ComboboxField } from "@/components/ui/combobox";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { FormLabel } from "@/components/ui/form-label";
import { Input } from "@/components/ui/input";
import { SectionHeading } from "@/components/ui/section-heading";
import { SelectField } from "@/components/ui/select";
import { Text } from "@/components/ui/text";
import { ToggleGroupField } from "@/components/ui/toggle-group";
import { describeTrailSchedule, formatOccurrence } from "@/lib/trails";
import { Box, Grid, HStack, LStack, styled } from "@/styled-system/jsx";
import { capitalise, humanise } from "@/utils/text";

import { TrailRobotActionField } from "./TrailRobotActionField";
import { TrailScheduleFields } from "./TrailScheduleFields";
import {
  TrailFormSchema,
  TrailFormValues,
  initialTrailFormValues,
  trailFormPayload,
  trailFormSchedule,
  weekdays,
} from "./trailForm";

type Props = {
  trail?: Trail;
  onSaved?: (trail: Trail) => void | Promise<void>;
};

export type TrailEditorRobot = {
  id: string;
  name: string;
};

export type TrailEditorFormProps = {
  initialValue?: TrailMutableProps;
  robots: TrailEditorRobot[];
  robotsError?: string;
  robotsLoading?: boolean;
  onPreview: (schedule: RecurrenceSchedule) => Promise<string[]>;
  onSubmit: (payload: TrailMutableProps) => Promise<void>;
};

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

export function TrailEditor({ trail, onSaved }: Props) {
  const router = useRouter();
  const { mutate } = useSWRConfig();
  const robots = useRobotsList();
  const initialValue = trail ? mutableTrailValue(trail) : undefined;

  async function save(payload: TrailMutableProps) {
    const saved = trail
      ? await trailUpdate(trail.id, payload)
      : await trailCreate(payload);

    await mutate(getTrailListKey());

    if (trail && onSaved) {
      await onSaved(saved);
      return;
    }

    router.push(`/robots/trails/${saved.id}`);
    router.refresh();
  }

  return (
    <TrailEditorForm
      initialValue={initialValue}
      robots={robots.data?.robots ?? []}
      robotsError={
        robots.error ? "Could not load available Robots." : undefined
      }
      robotsLoading={!robots.data && !robots.error}
      onPreview={async (schedule) => {
        const result = await trailSchedulePreview({ schedule });
        return result.occurrences;
      }}
      onSubmit={save}
    />
  );
}

export function TrailEditorForm({
  initialValue,
  robots,
  robotsError,
  robotsLoading,
  onPreview,
  onSubmit,
}: TrailEditorFormProps) {
  const browserTimezone = useMemo(getBrowserTimezone, []);
  const defaultStartDate = useMemo(getTomorrow, []);
  const form = useForm<TrailFormValues>({
    defaultValues: initialTrailFormValues(
      initialValue,
      browserTimezone,
      defaultStartDate,
    ),
    resolver: zodResolver(TrailFormSchema),
  });
  const actions = useFieldArray({ control: form.control, name: "actions" });
  const scheduleKind = form.watch("scheduleKind");
  const [preview, setPreview] = useState<{
    schedule: RecurrenceSchedule;
    occurrences: string[];
  }>();
  const [previewing, setPreviewing] = useState(false);
  const [requestError, setRequestError] = useState<string>();

  const timezoneItems = useMemo(
    () => getTimezoneItems(form.getValues("timezone")),
    [form],
  );
  const robotItems = useMemo(
    () => robots.map((robot) => ({ label: robot.name, value: robot.id })),
    [robots],
  );

  useEffect(() => {
    const subscription = form.watch((_, { name }) => {
      if (name && isScheduleField(name)) setPreview(undefined);
    });

    return () => subscription.unsubscribe();
  }, [form]);

  const save = form.handleSubmit(async (values) => {
    setRequestError(undefined);

    try {
      await onSubmit(
        trailFormPayload(values, initialValue?.status ?? "active"),
      );
    } catch (cause) {
      setRequestError(errorMessage(cause, "Could not save this Trail."));
    }
  });

  async function previewSchedule() {
    setPreviewing(true);
    setRequestError(undefined);

    try {
      const schedule = trailFormSchedule(form.getValues());
      const occurrences = await onPreview(schedule);
      setPreview({ schedule, occurrences });
    } catch (cause) {
      setRequestError(errorMessage(cause, "Could not preview this schedule."));
    } finally {
      setPreviewing(false);
    }
  }

  return (
    <LStack as="form" gap="8" alignItems="stretch" maxW="3xl" onSubmit={save}>
      <LStack gap="4" alignItems="stretch">
        <FormControl>
          <FormLabel>Name</FormLabel>
          <Input
            {...form.register("name")}
            aria-label="Name"
            placeholder="Scheduled moderation check"
          />
          <FormHelperText>
            Used to identify this Trail, its runs, and Robot sessions.
          </FormHelperText>
          <FormErrorText>{form.formState.errors.name?.message}</FormErrorText>
        </FormControl>

        <FormControl>
          <FormLabel>Description</FormLabel>
          <Input
            {...form.register("description")}
            aria-label="Description"
            placeholder="What this Trail does and why it runs"
          />
          <FormHelperText>
            Optional context for people who manage this Trail.
          </FormHelperText>
        </FormControl>
      </LStack>

      <FormSection
        title="Schedule"
        description="Choose when the Trail emits its trigger. Times use the selected timezone."
      >
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
      </FormSection>

      <FormSection
        title="Robot actions"
        description="Each Robot receives the same trigger and runs independently with a fresh session."
      >
        {actions.fields.map((action, index) => (
          <TrailRobotActionField
            key={action.id}
            actions={actions}
            form={form}
            index={index}
            robotItems={robotItems}
            robotsAvailable={robots.length > 0}
            robotsError={robotsError}
            robotsLoading={robotsLoading ?? false}
          />
        ))}

        <Box>
          <Button
            type="button"
            variant="outline"
            disabled={robots.length === 0}
            onClick={() =>
              actions.append({
                type: "robot_run",
                robot_ref: "",
                instruction: "",
              })
            }
          >
            Add Robot action
          </Button>
        </Box>
      </FormSection>

      {requestError && (
        <FormErrorText role="alert">{requestError}</FormErrorText>
      )}

      <HStack justifyContent="end">
        <Button
          type="submit"
          loading={form.formState.isSubmitting}
          loadingText={initialValue ? "Saving changes..." : "Creating Trail..."}
          disabled={robots.length === 0}
        >
          {initialValue ? "Save changes" : "Create Trail"}
        </Button>
      </HStack>
    </LStack>
  );
}

function FormSection({
  title,
  description,
  children,
}: {
  title: string;
  description: string;
  children: React.ReactNode;
}) {
  return (
    <LStack as="section" gap="4" alignItems="stretch">
      <Box maxW="layout.readable">
        <SectionHeading emphasis="strong">{title}</SectionHeading>
        <Text variant="supporting">{description}</Text>
      </Box>
      {children}
    </LStack>
  );
}

function mutableTrailValue(trail: Trail): TrailMutableProps {
  return {
    name: trail.name,
    description: trail.description,
    status:
      trail.status === "archived"
        ? "archived"
        : trail.status === "paused"
          ? "paused"
          : "active",
    trigger: trail.trigger,
    actions: trail.actions.map(({ action }) => action),
  };
}

function getBrowserTimezone(): string {
  return Intl.DateTimeFormat().resolvedOptions().timeZone || "UTC";
}

function getTomorrow(): string {
  const date = new Date(Date.now() + 86_400_000);
  return date.toISOString().slice(0, 10);
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
