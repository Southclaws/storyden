import Image, { getImageProps } from "next/image";
import Link from "next/link";

import { cx } from "@/styled-system/css";
import { Box, Grid, GridItem, HStack, VStack } from "@/styled-system/jsx";
import { linkButton } from "@/styled-system/patterns";

const { props: mobileLake } = getImageProps({
  src: "/square-nice-lake.webp",
  width: 1024,
  height: 1024,
  sizes: "100vw",
  alt: "",
});
const { props: desktopLake } = getImageProps({
  src: "/wide-nice-lake.webp",
  width: 3456,
  height: 1728,
  sizes: "100vw",
  alt: "",
  loading: "eager",
  fetchPriority: "high",
});

export function HeroSection() {
  return (
    <Grid>
      <GridItem gridRow="1/2" gridColumn="1/2">
        <picture>
          <source
            media="(max-width: 767px)"
            srcSet={mobileLake.srcSet}
            sizes="100vw"
            width={1024}
            height={1024}
          />
          <source
            media="(min-width: 768px)"
            srcSet={desktopLake.srcSet}
            sizes="100vw"
            width={3456}
            height={1728}
          />
          <img {...desktopLake} alt="" />
        </picture>
      </GridItem>

      <GridItem
        gridRow="1/2"
        gridColumn="1/2"
        zIndex={2}
        display="flex"
        justifyContent="center"
        alignItems="center"
        background="linear-gradient(180deg, rgba(59, 83, 111, 0) 49.48%, #000000 97.92%)"
      >
        <VStack gap="6" pt="16">
          <Box as="h1" width={[36, 40, 40, 48]}>
            <Image
              className="hero-logo"
              src="/brand/fullmark_newspaper_vertical_large.png"
              width="1790"
              height="1170"
              alt="Storyden — A home for your people and what they know."
              sizes="(min-width: 1024px) 192px, (min-width: 640px) 160px, 144px"
              loading="eager"
            />
          </Box>

          <HStack gap={4}>
            <Link
              className={linkButton({
                backgroundColor: "white",
                boxShadow: "xl",
              })}
              href="/docs/introduction"
              prefetch={false}
            >
              Get Started
            </Link>
            <Link
              target="_blank"
              className={cx(
                "story__text-overlay",
                linkButton({
                  backdropBlur: "lg",
                  backdropFilter: "auto",
                  backgroundColor: "rgba(98, 98, 98, 0.5)",
                  boxShadow: "xl",
                  color: "white",
                  _hover: {
                    color: "black",
                    background: "white",
                  },
                }),
              )}
              href="https://makeroom.club"
            >
              Live Demo
            </Link>
          </HStack>
          <p className="hero-description">
            An open-source forum and knowledgebase for communities that want
            their conversations, discoveries, and wisdom to last.
          </p>
        </VStack>
      </GridItem>
    </Grid>
  );
}
