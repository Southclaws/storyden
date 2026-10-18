import { Admonition } from "@/components/ui/admonition";
import * as Alert from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { FormLabel } from "@/components/ui/form-label";
import { NumberInputField } from "@/components/ui/number-input";
import { PageHeading } from "@/components/ui/page-heading";
import { Text } from "@/components/ui/text";
import { Flex, LStack, WStack, styled } from "@/styled-system/jsx";

import { type Props, useMediaSettings } from "./useMediaSettings";

export function MediaSettingsForm(props: Props) {
  const { clearErrors, control, formState, onSubmit } = useMediaSettings(props);
  const processing = props.settings.services?.media;
  const estimatedMemory = processing?.image_estimated_memory;
  const exceedsBudget =
    estimatedMemory !== undefined &&
    processing?.image_memory_budget !== undefined &&
    estimatedMemory > processing.image_memory_budget;

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
          <PageHeading>Media settings</PageHeading>
          <Button type="submit" loading={formState.isSubmitting}>
            Save
          </Button>
        </WStack>
        <Text variant="supporting">
          Set image upload limits for attachments, avatars, and branding.
        </Text>
      </LStack>

      <Flex flexDir={{ base: "column", md: "row" }} gap="2">
        <FormControl
          data-invalid={formState.errors.imageMaxSizeMiB ? "" : undefined}
        >
          <FormLabel id="image-max-size-label">
            Maximum image size (MiB)
          </FormLabel>
          <NumberInputField
            control={control}
            name="imageMaxSizeMiB"
            ariaLabel="Maximum image size (MiB)"
            min={1 / 1024 ** 2}
            max={Math.floor(Number.MAX_SAFE_INTEGER / 1024 ** 2)}
            step={1}
            inputProps={{
              id: "image-max-size",
              "aria-labelledby": "image-max-size-label",
              "aria-describedby": "image-max-size-help image-max-size-error",
            }}
          />
          <FormHelperText id="image-max-size-help">
            Maximum encoded file size for an uploaded image. Other file types
            are unaffected.
          </FormHelperText>
          <FormErrorText id="image-max-size-error">
            {formState.errors.imageMaxSizeMiB?.message}
          </FormErrorText>
        </FormControl>

        <FormControl
          data-invalid={formState.errors.imageMaxMegapixels ? "" : undefined}
        >
          <FormLabel id="image-max-resolution-label">
            Maximum image resolution (megapixels)
          </FormLabel>
          <NumberInputField
            control={control}
            name="imageMaxMegapixels"
            ariaLabel="Maximum image resolution (megapixels)"
            min={0.000001}
            max={Math.floor(Number.MAX_SAFE_INTEGER / 1_000_000)}
            step={1}
            inputProps={{
              id: "image-max-resolution",
              "aria-labelledby": "image-max-resolution-label",
              "aria-describedby":
                "image-max-resolution-help image-max-resolution-error",
            }}
          />
          <FormHelperText id="image-max-resolution-help">
            Width multiplied by height, in millions of pixels. Vector images
            have no pixel limit.
          </FormHelperText>
          <FormErrorText id="image-max-resolution-error">
            {formState.errors.imageMaxMegapixels?.message}
          </FormErrorText>
        </FormControl>
      </Flex>

      {formState.errors.root && (
        <Admonition
          kind="failure"
          value={true}
          onChange={() => clearErrors("root")}
        >
          {formState.errors.root.message}
        </Admonition>
      )}

      <Flex flexDir={{ base: "column", md: "row" }} gap="4">
        <LStack flex="1" gap="1">
          <Text variant="supporting">
            Uploaded originals are preserved. Optimized versions are generated
            in the background within the server's memory budget.
          </Text>
          {processing?.image_memory_budget !== undefined && (
            <Text variant="supporting">
              Worker memory budget:{" "}
              {Math.round(processing.image_memory_budget / 1024 ** 2)} MiB per
              server process. Configure this with IMAGE_MEMORY_BUDGET on the
              server.
            </Text>
          )}
        </LStack>
        <LStack flex="1" gap="1">
          <Text variant="supporting">Full HD: 1920 × 1080 → 2.1 MP</Text>
          <Text variant="supporting">4K: 3840 × 2160 → 8.3 MP</Text>
          <Text variant="supporting">
            Phone photograph: 4000 × 3000 → 12 MP
          </Text>
          <Text variant="supporting">
            Professional camera photograph: 6000 × 4000 → 24 MP
          </Text>
          <Text variant="supporting">
            Photo sizes vary by device and camera mode.
          </Text>
        </LStack>
      </Flex>

      {exceedsBudget && processing && (
        <Alert.Root tone="warning" role="status">
          <Alert.Content>
            <Alert.Title>
              Some images may exceed the processing budget
            </Alert.Title>
            <Alert.Description>
              A square image at the saved upload limits is estimated to need{" "}
              {Math.ceil(estimatedMemory / 1024 ** 2)} MiB to generate the
              largest fitted version. Actual memory use varies by image
              dimensions and format. Images that exceed the budget remain
              available as originals. Lower the upload limits or increase the
              server's memory budget to allow more images to be optimized.
            </Alert.Description>
          </Alert.Content>
        </Alert.Root>
      )}
    </styled.form>
  );
}
