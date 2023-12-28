import { Link as LinkSchema, ThreadReference } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import { Byline } from "src/components/content/Byline";
import { CollectionMenu } from "src/components/content/CollectionMenu/CollectionMenu";
import { Heading1 } from "src/theme/components/Heading/Index";
import { Link } from "src/theme/components/Link";

import { FeedItemMenu } from "../common/FeedItemMenu/FeedItemMenu";
import { Empty } from "../common/PostRef/Empty";

import {
  Box,
  Flex,
  HStack,
  LinkBox,
  LinkOverlay,
  VStack,
  styled,
} from "@/styled-system/jsx";

type Props = {
  thread: ThreadReference;
  onDelete: () => void;
};

export function LinkPost(props: Props) {
  const session = useSession();

  const permalink = `/t/${props.thread.slug}`;
  const link = props.thread.link as LinkSchema;
  const asset = link.assets?.[0] ?? props.thread.assets?.[0];

  return (
    <LinkBox>
      <styled.article
        display="flex"
        flexDir="column"
        width="full"
        boxShadow="md"
        borderRadius="md"
        backgroundColor="white"
        overflow="hidden"
      >
        <Box display="flex" w="full" height="24">
          <Box flexGrow="1" flexShrink="0" width="32">
            {asset ? (
              <styled.img
                src={asset.url}
                height="full"
                width="full"
                objectPosition="center"
                objectFit="cover"
              />
            ) : (
              <VStack justify="center" w="full" h="full">
                <Empty />
              </VStack>
            )}
          </Box>

          <VStack
            w="full"
            alignItems="start"
            justifyContent="space-evenly"
            gap="0"
            p="2"
          >
            <Flex width="full" justifyContent="space-between">
              <Heading1 size="sm" lineClamp={2}>
                <LinkOverlay href={permalink}>
                  {/* TODO: Next.js Link */}
                  {props.thread.title}
                </LinkOverlay>
              </Heading1>

              <Link flexShrink="0" kind="ghost" size="xs" href={link.url}>
                {link.domain}
              </Link>
            </Flex>

            <styled.p lineClamp={2}>
              <styled.span color="gray.500">{link.description}</styled.span>
            </styled.p>
          </VStack>
        </Box>

        <Flex justifyContent="space-between" p="2">
          <Byline
            href={permalink}
            author={props.thread.author}
            time={new Date(props.thread.createdAt)}
            updated={new Date(props.thread.updatedAt)}
          />

          <HStack>
            {session && <CollectionMenu thread={props.thread} />}
            <FeedItemMenu thread={props.thread} onDelete={props.onDelete} />
          </HStack>
        </Flex>
      </styled.article>
    </LinkBox>
  );
}
