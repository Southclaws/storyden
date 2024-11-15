import { Box, HStack } from "@/styled-system/jsx";

export function ColourPreview() {
  return (
    <>
      <Box p="1">
        <HStack maxWidth="full" justify="space-between" gap="1" height="16">
          {[3, 2, 1].map((value) => (
            <Shade
              key={value}
              fill={`var(--accent-colour-flat-fill-${value})`}
              text={`var(--accent-colour-flat-text-${value})`}
            />
          ))}

          <Shade fill={`var(--accent-colour)`} text={`var(--text-colour)`} />

          {[1, 2, 3].map((value) => (
            <Shade
              key={value}
              fill={`var(--accent-colour-dark-fill-${value})`}
              text={`var(--accent-colour-dark-text-${value})`}
            />
          ))}
        </HStack>
        <Box w="full" textAlign="center" backgroundColor="white" color="black">
          Pure White
        </Box>
        <Box w="full" textAlign="center" backgroundColor="black" color="white">
          Pure Black
        </Box>
      </Box>
    </>
  );
}

function Shade({ fill, text }: { fill: string; text: string }) {
  return (
    <Box width="full" height="full" borderRadius="md" overflow="hidden">
      <svg
        width="100%"
        height="100%"
        viewBox="0 0 100 100"
        preserveAspectRatio="xMidYMax slice"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <rect width="100" height="100" rx="12" fill={fill} />
        <text
          x={50}
          y={50}
          textAnchor="middle"
          alignmentBaseline="central"
          fill={text}
        >
          Tt
        </text>
      </svg>
    </Box>
  );
}
