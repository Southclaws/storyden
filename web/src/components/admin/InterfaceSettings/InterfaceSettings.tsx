import { Controller } from "react-hook-form";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { FormControl } from "@/components/ui/form-control";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { FormLabel } from "@/components/ui/form-label";
import { Heading } from "@/components/ui/heading";
import { FormNumberInputField } from "@/components/ui/number-input";
import { FormRadioGroupField } from "@/components/ui/radio-group";
import { WStack, styled } from "@/styled-system/jsx";

import { Props, useInterfaceSettings } from "./useInterfaceSettings";

export function InterfaceSettingsForm(props: Props) {
  const { control, signaturesEnabled, formState, onSubmit } =
    useInterfaceSettings(props);

  return (
    <styled.form
      width="full"
      display="flex"
      flexDirection="column"
      gap="4"
      onSubmit={onSubmit}
    >
      <WStack>
        <Heading size="md">Interface settings</Heading>
        <Button type="submit" loading={formState.isSubmitting}>
          Save
        </Button>
      </WStack>

      <FormControl>
        <FormLabel>Default editor</FormLabel>
        <FormRadioGroupField
          control={control}
          name="editorMode"
          items={[
            { label: "Rich text", value: "richtext" },
            { label: "Markdown", value: "markdown" },
          ]}
        />
        <FormHelperText>
          Choose the default editor for composing threads, replies and pages.
        </FormHelperText>
      </FormControl>

      <FormControl>
        <FormLabel>Signatures</FormLabel>
        <Controller
          control={control}
          name="signaturesEnabled"
          render={({ field }) => (
            <Checkbox
              size="sm"
              checked={!!field.value}
              onCheckedChange={({ checked }) => {
                field.onChange(checked === true);
              }}
            >
              Enable member signatures
            </Checkbox>
          )}
        />
        <FormHelperText>
          When disabled, signatures are hidden under posts and on profiles.
        </FormHelperText>
      </FormControl>

      <FormControl>
        <FormLabel>Signature max height (px)</FormLabel>
        <FormNumberInputField
          control={control}
          name="signatureMaxHeight"
          min={32}
          max={2000}
          disabled={!signaturesEnabled}
        />
        <FormHelperText>
          Limits how tall member signatures can appear below posts.
        </FormHelperText>
      </FormControl>

      <FormControl>
        <FormLabel>Signature max characters</FormLabel>
        <FormNumberInputField
          control={control}
          name="signatureMaxChars"
          min={1}
          max={10000}
          disabled={!signaturesEnabled}
        />
        <FormHelperText>
          Visible characters, not including HTML tags.
        </FormHelperText>
      </FormControl>
    </styled.form>
  );
}
