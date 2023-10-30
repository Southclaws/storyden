"use client";

import { groupBy, keys, partition } from "lodash";

import { Button } from "src/theme/components/Button";
import { Link } from "src/theme/components/Link";
import { getColourVariants } from "src/utils/colour";

import { Box, HStack, VStack, styled } from "@/styled-system/jsx";

export function Palette({ accent_colour, theme }: any) {
  const colours = getColourVariants(accent_colour);

  const variables = keys(colours);

  const { flat, dark, other } = groupBy(variables, (v) => {
    if (v.includes("flat")) return "flat";
    if (v.includes("dark")) return "dark";
    return "other";
  });

  const [flatText, flatFill] = partition(flat, (v) => v.includes("text"));
  const [darkText, darkFill] = partition(dark, (v) => v.includes("text"));

  return (
    <>
      {other?.map((v) => (
        <styled.p key={v} style={{ backgroundColor: `var(${v})` }}>
          {v}: {colours[v]}
        </styled.p>
      ))}

      <HStack w="full" gap="0">
        <VStack w="full" gap="0">
          {flatFill?.map((v, i) => (
            <styled.p
              key={v}
              w="full"
              style={{
                backgroundColor: `var(${v})`,
                color: `var(${flatText[i]})`,
              }}
            >
              {v}
            </styled.p>
          ))}
        </VStack>

        <VStack w="full" gap="0">
          {darkFill?.map((v, i) => (
            <styled.p
              key={v}
              w="full"
              style={{
                backgroundColor: `var(${v})`,
                color: `var(${darkText[i]})`,
              }}
            >
              {v}
            </styled.p>
          ))}
        </VStack>
      </HStack>

      <HStack>
        <VStack>
          <Button kind="neutral" size="xs">
            xs
          </Button>
          <Button kind="neutral" size="sm">
            sm
          </Button>
          <Button kind="neutral" size="md">
            md
          </Button>
          <Button kind="neutral" size="lg">
            lg
          </Button>
          <Button kind="neutral" size="xl">
            xl
          </Button>
          <Button kind="neutral" size="2xl">
            2xl
          </Button>
        </VStack>

        <VStack>
          <Button kind="neutral">Neutral</Button>
          <Button kind="primary">Primary</Button>
          <Button kind="secondary">Secondary</Button>
          <Button kind="destructive">Destructive</Button>
        </VStack>

        <VStack>
          <Button kind="neutral" disabled>
            Neutral
          </Button>
          <Button kind="primary" disabled>
            Primary
          </Button>
          <Button kind="secondary" disabled>
            Secondary
          </Button>
          <Button kind="destructive" disabled>
            Destructive
          </Button>
        </VStack>

        <VStack>
          <Link href="/xdev" kind="neutral">
            Neutral internal
          </Link>
          <Link href="#" kind="primary">
            Primary
          </Link>
          <Link href="#" kind="secondary">
            Secondary
          </Link>
        </VStack>
      </HStack>

      <Box className="typography">
        <pre>{theme}</pre>
      </Box>
    </>
  );
}
