import {
  Button,
  FormControl,
  FormHelperText,
  FormLabel,
  HStack,
  Heading,
  Input,
  VStack,
} from "@chakra-ui/react";

import { useGetInfo } from "src/api/openapi/misc";
import { ColourField } from "src/components/ColourInput/ColourInput";
import { Unready } from "src/components/Unready";

import { SettingsSection } from "../SettingsSection/SettingsSection";

import { Props, useBrandSettings } from "./useBrandSettings";

function BrandSettingsForm(props: Props) {
  const { register, control, onSubmit, onColourChangePreview } =
    useBrandSettings(props);

  return (
    <SettingsSection>
      <Heading size="sm">Brand settings</Heading>

      <VStack as="form" width="full" gap={2} align="start" onSubmit={onSubmit}>
        <FormControl gap={2}>
          <FormLabel>Title</FormLabel>
          <Input {...register("title")} maxW="20em" />
          <FormHelperText>
            The name of your community. This appears in the sidebar, Google
            indexing and tab titles.
          </FormHelperText>
        </FormControl>

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
