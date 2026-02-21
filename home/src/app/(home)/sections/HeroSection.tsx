import { Box, Grid, GridItem, HStack, VStack } from "@/styled-system/jsx";
import Image from "next/image";
import Link from "next/link";

import { cx } from "@/styled-system/css";
import { linkButton } from "@/styled-system/patterns";

export function HeroSection() {
  return (
    <Grid>
      <GridItem gridRow="1/2" gridColumn="1/2">
        <picture>
          <source
            media="(max-width: 768px)"
            srcSet="square-nice-lake.webp"
            width={1024}
            height={1024}
          />
          <source
            media="(min-width: 768px)"
            srcSet="wide-nice-lake.webp"
            width={3456}
            height={1728}
          />
          <img
            src="wide-nice-lake.webp"
            role="presentation"
            alt="A sun-lit lake sitting before tall snow-covered mountains in the distance."
            width={3456}
            height={1728}
          />
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
          <Box width={[36, 40, 40, 48]}>
            <Image
              src="/brand/fullmark_newspaper_vertical_large.png"
              width="1790"
              height="1170"
              alt="The Storyden logo"
            />
          </Box>

          <HStack gap={4}>
            <Link
              className={linkButton({
                backgroundColor: "white",
                boxShadow: "xl",
              })}
              href="/docs/introduction"
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
                })
              )}
              href="https://makeroom.club"
            >
              Live Demo
            </Link>
          </HStack>
        </VStack>
      </GridItem>
    </Grid>
  );
}
