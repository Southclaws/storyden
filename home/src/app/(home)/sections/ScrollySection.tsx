import Image from "next/image";

import { Box, Grid, GridItem, VStack, styled } from "@/styled-system/jsx";

const featureRows = [
  {
    id: "smart",
    title: "Built-in smarts.",
    title2: "opt-in brains.",
    paragraphs: [
      "Language model integration, at the core.",
      "But only if you want it.",
    ],
    image: "/sort-search.png",
    imageAlt:
      "A screenshot of Storyden's search and ask feature surfacing a member's vague question looking for a specific thread about design books.",
    surface: "Primary.moonlit",
    color: "Shades.newspaper",
  },
  {
    id: "sort",
    title: "Sort. Search. Ask.",
    paragraphs: [
      "Allow your collective knowledge to grow without losing great ideas to the void of time, banished to the archive.",
      "Looking for something?",
      "Just ask.",
    ],
    image: "/organise.png",
    imageAlt:
      "A graphic showing the flow of a question leading to a directory leading to a specific link to a website.",
    surface: "Primary.forest",
    color: "Shades.newspaper",
  },
  {
    id: "curate",
    title: "Curate, effortlessly.",
    paragraphs: [
      "It doesn't stop at threads. Let your community's Library grow with structured and organised pages.",
      "Alexandria-scale corpus? No worries, let Storyden's intelligence sort it out.",
    ],
    image: "/curate.png",
    imageAlt:
      "A diagram of the link fetching flow going from fetching, assigning tags and creating a page to demonstrate the AI assisted curation capabilities.",
    surface: "Primary.campfire",
    color: "Primary.moonlit",
  },
];

export function ScrollySection() {
  return (
    <styled.section
      bgColor="Mono.ink"
      position="relative"
      overflow="hidden"
      py={{ base: "10", md: "16", lg: "20" }}
      _before={{
        content: '""',
        position: "absolute",
        inset: 0,
        backgroundImage:
          "radial-gradient(circle at 1px 1px, rgba(255, 255, 255, 0.12) 1px, transparent 0)",
        backgroundSize: "20px 20px",
        opacity: 0.22,
        pointerEvents: "none",
      }}
    >
      <VStack
        position="relative"
        zIndex="1"
        gap="5"
        w="full"
        maxW="breakpoint-xl"
        mx="auto"
        px={{ base: "4", sm: "8", md: "12", lg: "16" }}
      >
        {featureRows.map((feature, index) => (
          <Grid
            key={feature.id}
            gap="4"
            w="full"
            alignItems="stretch"
            gridTemplateColumns={{ base: "1fr", lg: "1fr 1fr" }}
          >
            <GridItem
              order={{
                base: 1,
                lg: index % 2 === 0 ? 1 : 2,
              }}
              bgColor={feature.surface}
              color={feature.color}
              borderRadius="2xl"
              borderWidth="thin"
              borderStyle="solid"
              borderColor="black/20"
              p={{ base: "5", md: "8" }}
              display="flex"
              alignItems="center"
            >
              <VStack alignItems="start" gap="3">
                <styled.h1
                  fontSize={{ base: "2xl", md: "4xl", lg: "5xl" }}
                  lineHeight="1.05"
                  textWrap="balance"
                >
                  {feature.title}
                  {feature.title2 ? (
                    <>
                      <br />
                      {feature.title2}
                    </>
                  ) : null}
                </styled.h1>

                {feature.paragraphs.map((paragraph) => (
                  <styled.p
                    key={paragraph}
                    fontSize={{ base: "md", md: "xl" }}
                    lineHeight="relaxed"
                    textWrap="pretty"
                  >
                    {paragraph}
                  </styled.p>
                ))}
              </VStack>
            </GridItem>

            <GridItem
              order={{
                base: 2,
                lg: index % 2 === 0 ? 2 : 1,
              }}
              bgColor="white/95"
              borderRadius="2xl"
              borderWidth="thin"
              borderStyle="solid"
              borderColor="Shades.slate/40"
              p={{ base: "4", md: "6" }}
              minH={{ base: "220px", md: "320px" }}
              display="flex"
              alignItems="center"
              justifyContent="center"
              boxShadow="0 16px 32px -24px rgba(0, 0, 0, 0.55)"
            >
              <Image
                src={feature.image}
                alt={feature.imageAlt}
                width={1320}
                height={864}
              />
            </GridItem>
          </Grid>
        ))}
      </VStack>
    </styled.section>
  );
}
