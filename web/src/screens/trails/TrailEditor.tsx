"use client";

import { createListCollection } from "@ark-ui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import { Controller, useFieldArray, useForm } from "react-hook-form";
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
import { CardBox } from "@/components/ui/card-box";
import { ComboboxField } from "@/components/ui/combobox";
import { DatePickerField } from "@/components/ui/date-picker";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { FormLabel } from "@/components/ui/form-label";
import { Input } from "@/components/ui/input";
import { NumberInputField } from "@/components/ui/number-input";
import { SelectField } from "@/components/ui/select";
import { Text } from "@/components/ui/text";
import { ToggleGroupField } from "@/components/ui/toggle-group";
import { describeTrailSchedule, formatOccurrence } from "@/lib/trails";
import { Box, Grid, HStack, LStack, WStack, styled } from "@/styled-system/jsx";

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

const weekdayItems = weekdays.map((day) => ({
  label: day.slice(0, 3),
  value: day,
  description: capitalise(day),
}));

export function TrailEditor({ trail, onSaved }: Props) {
  const router = useRouter();
  const { mutate } = useSWRConfig();
  const robots = useRobotsList();
  const initialValue = useMemo(
    () => (trail ? mutableTrailValue(trail) : undefined),
    [trail],
  );

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
    <styled.form
      display="flex"
      flexDirection="column"
      gap="8"
      maxW="3xl"
      width="full"
      onSubmit={save}
    >
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

        <ScheduleFields form={form} scheduleKind={scheduleKind} />

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
          <CardBox as="section" key={action.id}>
            <LStack gap="4" alignItems="stretch">
              <WStack alignItems="center">
                <Text fontWeight="semibold">Robot action {index + 1}</Text>
                {actions.fields.length > 1 && (
                  <Button
                    type="button"
                    size="sm"
                    variant="ghost"
                    onClick={() => actions.remove(index)}
                  >
                    Remove
                  </Button>
                )}
              </WStack>

              <FormControl>
                <FormLabel>Robot</FormLabel>
                <ComboboxField
                  control={form.control}
                  name={`actions.${index}.robot_ref`}
                  items={robotItems}
                  placeholder={
                    robotsLoading ? "Loading Robots..." : "Search Robots"
                  }
                  ariaLabel={`Robot action ${index + 1}`}
                  disabled={robotsLoading || robots.length === 0}
                />
                {robotsError ? (
                  <FormErrorText>{robotsError}</FormErrorText>
                ) : robots.length === 0 && !robotsLoading ? (
                  <FormHelperText>
                    Create a Robot before adding it to this Trail.
                  </FormHelperText>
                ) : (
                  <FormHelperText>
                    This Robot runs unattended for every occurrence.
                  </FormHelperText>
                )}
                <FormErrorText>
                  {form.formState.errors.actions?.[index]?.robot_ref?.message}
                </FormErrorText>
              </FormControl>

              <FormControl>
                <FormLabel>Instruction</FormLabel>
                <Controller
                  control={form.control}
                  name={`actions.${index}.instruction`}
                  render={({ field }) => (
                    <styled.textarea
                      {...field}
                      aria-label={`Unattended instruction ${index + 1}`}
                      placeholder="Describe the result the Robot must produce on every run."
                      rows={6}
                      width="full"
                      padding="3"
                      fontSize="sm"
                      lineHeight="relaxed"
                      borderWidth="thin"
                      borderStyle="solid"
                      borderColor="border.default"
                      borderRadius="sm"
                      backgroundColor="background.control"
                      resize="vertical"
                      _focus={{
                        outline: "none",
                        borderColor: "accent.solid",
                      }}
                    />
                  )}
                />
                <FormHelperText>
                  This fixed instruction starts a new unattended Robot session.
                </FormHelperText>
                <FormErrorText>
                  {form.formState.errors.actions?.[index]?.instruction?.message}
                </FormErrorText>
              </FormControl>
            </LStack>
          </CardBox>
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
    </styled.form>
  );
}

function ScheduleFields({
  form,
  scheduleKind,
}: {
  form: ReturnType<typeof useForm<TrailFormValues>>;
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
    <styled.section display="flex" flexDirection="column" gap="4">
      <Box maxW="layout.readable">
        <styled.h2 fontSize="lg" fontWeight="semibold">
          {title}
        </styled.h2>
        <Text variant="supporting">{description}</Text>
      </Box>
      {children}
    </styled.section>
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
      label: timezone.replaceAll("_", " "),
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

function capitalise(value: string): string {
  return value.charAt(0).toUpperCase() + value.slice(1);
}

function errorMessage(cause: unknown, fallback: string): string {
  return cause instanceof Error ? cause.message : fallback;
}
