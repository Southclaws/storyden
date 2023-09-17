import {
  Button,
  FormControl,
  FormHelperText,
  FormLabel,
  HStack,
  Heading,
  Input,
  Stack,
  VStack,
} from "@chakra-ui/react";

import { useGetInfo } from "src/api/openapi/misc";
import { ColourField } from "src/components/form/ColourInput/ColourInput";
import { Unready } from "src/components/site/Unready";

import { SettingsSection } from "../SettingsSection/SettingsSection";

import { IconEditor } from "./IconEditor/IconEditor";
import { Props, useBrandSettings } from "./useBrandSettings";

function BrandSettingsForm(props: Props) {
  const {
    register,
    control,
    onSubmit,
    currentIcon,
    onSaveIcon,
    onColourChangePreview,
  } = useBrandSettings(props);

  return (
    <SettingsSection>
      <Heading size="md">Brand settings</Heading>

      <VStack as="form" width="full" gap={4} align="start" onSubmit={onSubmit}>
        <Stack
          gap={4}
          direction={{
            base: "column",
            lg: "row",
          }}
        >
          <FormControl display="flex" flexDirection="column">
            <FormLabel>Community name</FormLabel>
            <Input {...register("title")} />
            <FormHelperText>
              The name of your community. This appears in the sidebar, Google
              indexing and tab titles.
            </FormHelperText>

            <FormHelperText>
              Your icon will be automatically resized and optimised for various
              devices. It is used for the website favicon and a PWA app icon for
              iOS and Android devices.
            </FormHelperText>
          </FormControl>

          <FormControl display="flex" flexDirection="column">
            <FormLabel>Icon</FormLabel>

            <IconEditor initialValue={currentIcon} onSave={onSaveIcon} />
          </FormControl>
        </Stack>

        <FormControl>
          <FormLabel>Description</FormLabel>
          <Input {...register("description")} />
          <FormHelperText>
            Describe your community with a few words here. This will be used for
            Google indexing, social previews and the PWA manifest.
          </FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>Colour</FormLabel>
          <ColourField
            defaultValue={props.accent_colour}
            control={control}
            onUpdate={onColourChangePreview}
            {...register("accentColour")}
          />
          <FormHelperText>
            Pick a colour that best represents your community or brand. It will
            be used throughout the site for accenting certain elements such as
            buttons, mobile browser borders, PWA theme, etc.
          </FormHelperText>
        </FormControl>

        <HStack justify="end">
          <Button type="submit">Save</Button>
        </HStack>
      </VStack>
    </SettingsSection>
  );
}

export function BrandSettings() {
  const { data, error } = useGetInfo();
  if (!data) return <Unready {...error} />;

  return <BrandSettingsForm {...data} />;
}
