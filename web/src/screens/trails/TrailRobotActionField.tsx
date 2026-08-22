import type { UseFieldArrayReturn, UseFormReturn } from "react-hook-form";
import { Controller } from "react-hook-form";

import { Button } from "@/components/ui/button";
import { CardBox } from "@/components/ui/card-box";
import { ComboboxField } from "@/components/ui/combobox";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { FormLabel } from "@/components/ui/form-label";
import { Text } from "@/components/ui/text";
import { Textarea } from "@/components/ui/textarea";
import { LStack, WStack } from "@/styled-system/jsx";

import { TrailFormValues } from "./trailForm";

type TrailRobotActionFieldProps = {
  actions: UseFieldArrayReturn<TrailFormValues, "actions", "id">;
  form: UseFormReturn<TrailFormValues>;
  index: number;
  robotItems: { label: string; value: string }[];
  robotsAvailable: boolean;
  robotsError?: string;
  robotsLoading: boolean;
};

export function TrailRobotActionField({
  actions,
  form,
  index,
  robotItems,
  robotsAvailable,
  robotsError,
  robotsLoading,
}: TrailRobotActionFieldProps) {
  return (
    <CardBox as="section">
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
            placeholder={robotsLoading ? "Loading Robots..." : "Search Robots"}
            ariaLabel={`Robot action ${index + 1}`}
            disabled={robotsLoading || !robotsAvailable}
          />
          {robotsError ? (
            <FormErrorText>{robotsError}</FormErrorText>
          ) : !robotsAvailable && !robotsLoading ? (
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
              <Textarea
                {...field}
                aria-label={`Unattended instruction ${index + 1}`}
                placeholder="Describe the result the Robot must produce on every run."
                rows={6}
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
  );
}
