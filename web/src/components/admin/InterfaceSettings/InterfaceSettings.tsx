import { Button } from "@/components/ui/button";
import { CheckboxField } from "@/components/ui/checkbox";
import { FormControl } from "@/components/ui/form-control";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { FormLabel } from "@/components/ui/form-label";
import { NumberInputField } from "@/components/ui/number-input";
import { PageHeading } from "@/components/ui/page-heading";
import { RadioGroupField } from "@/components/ui/radio-group";
import { Text } from "@/components/ui/text";
import { LStack, WStack, styled } from "@/styled-system/jsx";

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
      <LStack gap="1">
        <WStack>
          <PageHeading>Interface settings</PageHeading>
          <Button type="submit" loading={formState.isSubmitting}>
            Save
          </Button>
        </WStack>
        <Text variant="supporting">
          Configure the default editing experience and how member content is
          displayed.
        </Text>
      </LStack>

      <FormControl>
        <FormLabel>Default editor</FormLabel>
        <RadioGroupField
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
        <CheckboxField control={control} name="signaturesEnabled" size="sm">
          Enable member signatures
        </CheckboxField>
        <FormHelperText>
          When disabled, signatures are hidden under posts and on profiles.
        </FormHelperText>
      </FormControl>

      <FormControl>
        <FormLabel>Signature max height (px)</FormLabel>
        <NumberInputField
          control={control}
          name="signatureMaxHeight"
          ariaLabel="Signature max height in pixels"
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
        <NumberInputField
          control={control}
          name="signatureMaxChars"
          ariaLabel="Signature max characters"
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
