"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useState } from "react";
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
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { FormLabel } from "@/components/ui/form-label";
import { Input } from "@/components/ui/input";
import { SectionHeading } from "@/components/ui/section-heading";
import { Text } from "@/components/ui/text";
import { Box, HStack, LStack } from "@/styled-system/jsx";

import { TrailRobotActionField } from "./TrailRobotActionField";
import { TrailTriggerFields } from "./TrailTriggerFields";
import {
  TrailFormSchema,
  TrailFormValues,
  initialTrailFormValues,
  trailFormPayload,
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
  const [browserTimezone] = useState(getBrowserTimezone);
  const [defaultStartDate] = useState(getTomorrow);
  const form = useForm<TrailFormValues>({
    defaultValues: initialTrailFormValues(
      initialValue,
      browserTimezone,
      defaultStartDate,
    ),
    resolver: zodResolver(TrailFormSchema),
  });
  const actions = useFieldArray({ control: form.control, name: "actions" });
  const [requestError, setRequestError] = useState<string>();

  const robotItems = robots.map((robot) => ({
    label: robot.name,
    value: robot.id,
  }));

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

      <FormSection title="Trigger" description="Choose what starts this Trail.">
        <TrailTriggerFields form={form} onPreview={onPreview} />
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

function errorMessage(cause: unknown, fallback: string): string {
  return cause instanceof Error ? cause.message : fallback;
}
