"use client";

import { groupBy, keys, partition } from "lodash";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { LinkButton } from "@/components/ui/link-button";
import { Box, HStack, VStack, styled } from "@/styled-system/jsx";
import { getColourVariants } from "@/utils/colour";

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
          <Button variant="ghost" size="xs">
            xs
          </Button>
          <Button variant="ghost" size="sm">
            sm
          </Button>
          <Button variant="ghost" size="md">
            md
          </Button>
          <Button variant="ghost" size="lg">
            lg
          </Button>
          <Button variant="ghost" size="xl">
            xl
          </Button>
          <Button variant="ghost" size="2xl">
            2xl
          </Button>
        </VStack>

        <VStack>
          <Button variant="ghost">Neutral</Button>
          <Button>Primary</Button>
          <Button colorPalette="red">Destructive</Button>
        </VStack>

        <VStack>
          <Button variant="ghost" disabled>
            Neutral
          </Button>
          <Button disabled>Primary</Button>
          <Button colorPalette="red" disabled>
            Destructive
          </Button>
        </VStack>

        <VStack>
          <LinkButton href="/xdev" variant="ghost">
            Neutral internal
          </LinkButton>
          <LinkButton href="#">Primary</LinkButton>
        </VStack>
      </HStack>

      <HStack>
        <VStack alignItems="start">
          <Checkbox checked={false}>Unchecked</Checkbox>

          <Checkbox size="sm" checked={true}>
            Checked
          </Checkbox>

          <Checkbox checked="indeterminate">Indeterminate</Checkbox>

          <Checkbox size="lg">Uncontrolled</Checkbox>
        </VStack>

        <VStack alignItems="start">
          <Checkbox defaultChecked={false}>Unchecked</Checkbox>

          <Checkbox size="sm" defaultChecked={true}>
            Checked
          </Checkbox>

          <Checkbox size="md" defaultChecked="indeterminate">
            Indeterminate
          </Checkbox>

          <Checkbox size="lg">Uncontrolled</Checkbox>
        </VStack>
      </HStack>

      <Box className="typography">
        <pre>{theme}</pre>
      </Box>
    </>
  );
}
