import { useI18n } from "@/i18n/provider";
import { Box, CardBox, Flex, HStack, styled } from "@/styled-system/jsx";
import { relativeTimestamp } from "@/utils/date";

import { BanIcon } from "../ui/icons/BanIcon";

type Props = {
  date: Date;
};

export function ProfileSuspendedBanner({ date }: Props) {
  const { locale, t } = useI18n();

  return (
    <CardBox
      p="0"
      borderColor="border.error"
      borderWidth="thin"
      borderStyle="dashed"
    >
      <Box bgColor="bg.error" borderTopRadius="md" pl="3" pr="2" py="1">
        <HStack gap="1" color="fg.error" fontSize="xs">
          <BanIcon w="4" />
          <p>{t("Suspended")}</p>
        </HStack>
      </Box>

      <Flex
        p="3"
        gap="4"
        direction={{ base: "column", md: "row" }}
        alignItems="start"
      >
        <styled.p color="fg.destructive" wordBreak="keep-all">
          {t("This member was suspended")}{" "}
          <styled.time textWrap="nowrap">
            {relativeTimestamp(date, locale)}
          </styled.time>
        </styled.p>
      </Flex>
    </CardBox>
  );
}
