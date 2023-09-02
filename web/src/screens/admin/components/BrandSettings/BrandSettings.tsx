import {
  Button,
  FormControl,
  FormHelperText,
  FormLabel,
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
  const { register, control, onSubmit } = useBrandSettings(props);

  return (
    <SettingsSection>
      <Heading size="sm">Brand settings</Heading>

      <VStack as="form" gap={2} align="start" onSubmit={onSubmit}>
        <FormControl gap={2}>
          <FormLabel>Title</FormLabel>
          <Input {...register("title")} />
          <FormHelperText>The name of your forum</FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>Description</FormLabel>
          <Input {...register("description")} />
          <FormHelperText>The description of your forum</FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>Colour</FormLabel>
          <ColourField
            defaultValue={props.accent_colour}
            control={control}
            {...register("accentColour")}
          />
          <FormHelperText>Your brand&apos;s colour</FormHelperText>
        </FormControl>

        <Button type="submit">Save</Button>
      </VStack>
    </SettingsSection>
  );
}

export function BrandSettings() {
  const { data, error } = useGetInfo();
  if (!data) return <Unready {...error} />;

  return <BrandSettingsForm {...data} />;
}
