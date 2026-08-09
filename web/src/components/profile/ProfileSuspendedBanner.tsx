import { formatDistanceToNow } from "date-fns";

import { CardBox } from "@/components/ui/card-box";
import { Text } from "@/components/ui/text";
import { Box, Flex, HStack, styled } from "@/styled-system/jsx";

import { BanIcon } from "../ui/icons/BanIcon";

type Props = {
  date: Date;
};

export function ProfileSuspendedBanner({ date }: Props) {
  return (
    <CardBox
      p="0"
      borderColor="status.danger.border"
      borderWidth="thin"
      borderStyle="dashed"
    >
      <Box
        bgColor="status.danger.surface"
        borderTopRadius="md"
        pl="3"
        pr="2"
        py="1"
      >
        <HStack gap="1" color="status.danger.content" fontSize="xs">
          <BanIcon w="4" />
          <p>Suspended</p>
        </HStack>
      </Box>

      <Flex
        p="3"
        gap="4"
        direction={{ base: "column", md: "row" }}
        alignItems="start"
      >
        <Text
          variant="supporting"
          color="status.danger.content"
          wordBreak="keep-all"
        >
          This member was suspended&nbsp;
          <styled.time textWrap="nowrap">
            {formatDistanceToNow(date, {
              addSuffix: true,
            })}
          </styled.time>
        </Text>
      </Flex>
    </CardBox>
  );
}
