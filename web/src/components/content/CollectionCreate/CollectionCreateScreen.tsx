import { UseDisclosureProps } from "src/utils/useDisclosure";

import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormHelperText } from "@/components/ui/form/FormHelperText";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { Input } from "@/components/ui/input";
import { useI18n } from "@/i18n/provider";
import { VStack, WStack, styled } from "@/styled-system/jsx";

import { Props, useCollectionCreate } from "./useCollectionCreate";

export function CollectionCreateScreen(props: Props) {
  const { register, onSubmit } = useCollectionCreate(props);
  const { t } = useI18n();

  return (
    <VStack alignItems="start" gap="4">
      <styled.p>
        {t(
          "Use collections to curate content from the community. Collections can include threads, pages and other items from the community knowledgebase.",
        )}
      </styled.p>
      <styled.form
        display="flex"
        flexDir="column"
        gap="4"
        w="full"
        onSubmit={onSubmit}
      >
        <FormControl>
          <FormLabel>{t("Name")}*</FormLabel>
          <Input {...register("name")} type="text" />
          <FormHelperText>{t("The name for your collection")}</FormHelperText>
        </FormControl>
        <FormControl>
          <FormLabel>{t("Description")}</FormLabel>

          {/* TODO: Make a larger textarea component for this. */}
          <Input {...register("description")} type="text" />
          <FormHelperText>
            {t("Optional description for your collection.")}
          </FormHelperText>
        </FormControl>

        <WStack>
          <Button
            flexGrow="1"
            type="button"
            variant="outline"
            onClick={props.onClose}
          >
            {t("Cancel")}
          </Button>
          <Button flexGrow="1" type="submit">
            {t("Create")}
          </Button>
        </WStack>
      </styled.form>
    </VStack>
  );
}
