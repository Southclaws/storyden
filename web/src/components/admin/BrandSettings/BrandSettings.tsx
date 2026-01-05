import { ColourField } from "src/components/form/ColourInput/ColourInput";

import { ContentFormField } from "@/components/content/ContentComposer/ContentField";
import { FormErrorText } from "@/components/ui/FormErrorText";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormHelperText } from "@/components/ui/form/FormHelperText";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { Heading } from "@/components/ui/heading";
import { Input } from "@/components/ui/input";
import {
  Box,
  CardBox,
  HStack,
  Stack,
  WStack,
  styled,
} from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { BannerEditor } from "./BannerEditor/BannerEditor";
import { IconEditor } from "./IconEditor/IconEditor";
import { Props, useBrandSettings } from "./useBrandSettings";

export function BrandSettingsForm(props: Props) {
  const {
    register,
    control,
    formState,
    onSubmit,
    currentIcon,
    onSaveIcon,
    onColourChangePreview,
  } = useBrandSettings(props);

  return (
    <CardBox className={lstack()}>
      <styled.form
        width="full"
        display="flex"
        flexDirection="column"
        gap="4"
        alignItems="start"
        onSubmit={onSubmit}
      >
        <WStack>
          <Heading size="md">Brand settings</Heading>
          <Button type="submit">Save</Button>
        </WStack>

        <Stack
          gap="4"
          direction={{
            base: "column",
            lg: "row",
          }}
        >
          <FormControl>
            <FormLabel>Community name</FormLabel>
            <Input {...register("title")} />
            <FormHelperText>
              The name of your community. This appears in the sidebar, Google
              indexing and tab titles.
            </FormHelperText>
          </FormControl>
        </Stack>

        <FormControl display="flex" flexDirection="column">
          <FormLabel>Icon</FormLabel>

          <IconEditor initialValue={currentIcon} onSave={onSaveIcon} />

          <FormHelperText>
            Your icon will be automatically resized and optimised for various
            devices. It is used for the website favicon and a PWA app icon for
            iOS and Android devices.
          </FormHelperText>
        </FormControl>

        <FormControl display="flex" flexDirection="column">
          <FormLabel>Banner</FormLabel>

          <BannerEditor />
          <FormHelperText>
            Your banner will be used for link previews on other platforms.
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
          <FormLabel>About</FormLabel>
          <CardBox>
            <ContentFormField
              control={control}
              name="content"
              // NOTE: Does not update if sidebar is changed. Doesn't matter...
              initialValue={props.settings.content}
              placeholder="About your community..."
            />
            <FormErrorText>{formState.errors.content?.message}</FormErrorText>
          </CardBox>
          <FormHelperText>
            You can write a longer description about your community here. You
            can use rich text formatting and include links and images.
          </FormHelperText>
        </FormControl>

        <FormControl>
          <FormLabel>Colour</FormLabel>
          <HStack>
            <Box>
              <ColourField
                defaultValue={props.settings.accent_colour}
                control={control}
                onUpdate={onColourChangePreview}
                {...register("accentColour")}
              />
            </Box>
          </HStack>

          <FormHelperText>
            Pick a colour that best represents your community or brand. It will
            be used throughout the site for accenting certain elements such as
            buttons, mobile browser borders, PWA theme, etc.
          </FormHelperText>
        </FormControl>

        <WStack justifyContent="end">
          <Button type="submit">Save</Button>
        </WStack>
      </styled.form>
    </CardBox>
  );
}
