import Image from "next/image";

import { css } from "@/styled-system/css";
import { Box, VStack, styled } from "@/styled-system/jsx";

function Permanence() {
  return (
    <VStack
      gap="8"
      maxW="breakpoint-sm"
      h="full"
      py="48"
      px="4"
      textAlign="center"
    >
      <styled.h2 color="Mono.ink/60">Communities deserve permanence.</styled.h2>
      <styled.p textWrap="balance">
        In a tangled web of fleeting <strong>feeds</strong> and forgotten
        threads, Storyden makes <strong>what matters</strong> stay discoverable,
        readable, and shareable. <strong>On your terms</strong>.
      </styled.p>
      <styled.p textWrap="balance">
        Own your <strong>data</strong>. Run your own Reddit, your own Pinterest,
        your own Hacker News,
        <wbr />
        &nbsp;your own <strong>corner</strong> of the web.
      </styled.p>
      <styled.p textWrap="balance" color="Mono.ink/60">
        <em>Collective memory. Without the noise.</em>
      </styled.p>
    </VStack>
  );
}

function Organise() {
  return (
    <VStack gap="8" h="full" py="48" bgColor="Shades.newspaper" w="full">
      <VStack p="4" maxW="prose" gap="4" textAlign="center">
        <styled.h2 color="Mono.ink/60">
          Conversations flow. Ideas grow.
        </styled.h2>

        <styled.p textWrap="balance">
          Allow your collective knowledge to flourish without losing great ideas
          to the void of time, banished to the archive.
        </styled.p>

        <Box maxW="lg">
          <Image
            className={css({
              borderRadius: "sm",
              boxShadow: "sm",
            })}
            src="/organise.png"
            alt="A diagram of the link fetching flow going from fetching, assigning tags and creating a page to demonstrate the AI assisted curation capabilities."
            width={1320 / 2}
            height={864 / 2}
          />
        </Box>
      </VStack>
    </VStack>
  );
}

function Gardens() {
  return (
    <VStack w="full" gap="8" py="48" bgColor="Mono.ink">
      <VStack p="4">
        <Box p="4">
          <styled.h2 color="Mono.slush/80">
            Knowledge&nbsp;gardens,
            <wbr /> not&nbsp;content&nbsp;farms.
          </styled.h2>
        </Box>

        <Box
          aspectRatio="1"
          maxW="sm"
          position="relative"
          bgColor="#F6F6F6"
          borderRadius="xl"
          boxShadow="xs"
          backgroundImage="linear-gradient(180deg, rgba(0, 0, 0, 0) 0%, rgba(245, 245, 245, 0.62) 63%, rgba(245, 245, 245, 1) 86%), url('/square-tree-smol.png')"
          backgroundSize="cover"
          backgroundPosition="top"
          overflowClipMargin="unset"
          display="flex"
          flexDir="column"
          justifyContent="space-between"
          overflow="hidden"
        >
          <Box height={{ base: "12" }} />

          <styled.p
            backgroundColor="white/20"
            backdropBlur="sm"
            backdropFilter="auto"
            // background="linear-gradient(180deg, rgba(1, 1, 1, 0.1) 25%, rgba(240, 240, 240, 1) 100%)"
            borderBottomRadius="md"
            p="6"
            textWrap="pretty"
          >
            From fan <strong>clubs</strong> to research groups, digital zines to
            esports <strong>teams</strong>. For anyone who cares about ideas and
            the <strong>people</strong> behind them.
          </styled.p>
        </Box>
      </VStack>
    </VStack>
  );
}

function Accessible() {
  return (
    <VStack w="full" gap="8" py="48" bgColor="Shades.iron">
      <VStack
        alignItems="center"
        justifyContent="center"
        minW={{ base: "0", lg: "160px" }}
      >
        <Image
          src="/accessibility.png"
          width="120"
          height="120"
          alt="The accessibility icon"
        />
      </VStack>

      <VStack maxW="prose" textAlign="center" color="Shades.newspaper">
        <h2>Optimised for humans.</h2>
        <styled.p textWrap="balance">
          Optimised for <strong>humans</strong>, ready for the web{" "}
          <strong>renaissance</strong>. A stable foundation for the future
          decades of internet citizens and the <strong>networks</strong> they
          build.
        </styled.p>
      </VStack>
    </VStack>
  );
}

function Curate() {
  return (
    <VStack
      w="full"
      gap="8"
      py="48"
      bgColor="Primary.moonlit"
      color="Shades.newspaper"
    >
      <VStack maxW="prose" textAlign="center" color="Shades.newspaper">
        <h2>Curate, effortlessly.</h2>
        <styled.p textWrap="balance">
          Let your community's <strong>Library</strong> grow with structured and
          organised pages. Alexandria-scale corpus? No worries, let Storyden's
          intelligence sort it out.
        </styled.p>
      </VStack>

      <Box maxW="sm">
        <Image
          className={css({
            borderRadius: "sm",
            boxShadow: "sm",
          })}
          src="/curate.png"
          alt=""
          width={1320 / 2}
          height={864 / 2}
        />
      </Box>
    </VStack>
  );
}

export function CollectiveMemorySection() {
  return (
    <styled.section w="full" bgColor="Mono.slush">
      <VStack position="relative" zIndex="1" w="full">
        <Permanence />
        <Organise />
        <Gardens />
        {/* <Accessible /> */}
        <Curate />
      </VStack>
    </styled.section>
  );
}
