import { Box, Button, Flex, keyframes } from "@chakra-ui/react";

function Circle({ c, x1, y1, x2, y2 }) {
  const move = keyframes`
    33% { transform: translate(${x1}px, ${y1}px) translateZ(0); }
    66% { transform: translate(${x2}px, ${y2}px) translateZ(0); }
  `;

  return (
    <Box
      pos="absolute"
      zIndex={1}
      top={0}
      bottom={0}
      left={x1}
      right={x2}
      borderRadius="50%"
      bgColor={c}
      width="100px"
      height="100px"
      filter="blur(20px)"
      animation={`${move} infinite ${Math.random() * 5 + 5}s ease-in-out`}
    />
  );
}

const colours = [
  //
  "hsla(260, 93.3%, 36.7%, 33%)",
  "hsla(260, 73.3%, 96.7%, 93%)",
  "hsla(258, 92.2%, 58.2%, 33%)",
  "hsla(258, 42.2%, 98.2%, 33%)",
  "hsla(254, 93.8%, 71.6%, 33%)",
  "hsla(254, 33.8%, 91.6%, 33%)",
  "hsla(249, 93.3%, 95.3%, 33%)",
  "hsla(249, 23.3%, 90.3%, 33%)",
  "hsla(256, 73.3%, 95.3%, 33%)",
  "hsla(256, 93.3%, 45.3%, 23%)",
];

function c() {
  return colours[(Math.random() * colours.length).toFixed(0)];
}

function Flamegrilled() {
  return (
    <Button
      height="60px"
      width="200px"
      fontSize={"2xl"}
      p={0}
      borderRadius="xl"
      bgColor="hsl(254, 53.8%, 11%)"
      boxShadow="inset 0 0 10px hsla(256, 63.3%, 56.3%, 25%)"
    >
      <Box
        borderRadius="xl"
        py="13px"
        overflow="hidden"
        height="100%"
        width="100%"
        css="mask-image: radial-gradient(black, white)"
      >
        <Box as="span" color="white" zIndex={9} textShadow="4px 2px 10px black">
          Join waitlist
        </Box>
        <Circle c={colours[0]} x1={0} x2={0} y1={-50} y2={-20} />
        <Circle c={colours[1]} x1={2} x2={20} y1={50} y2={-90} />
        <Circle c={colours[2]} x1={4} x2={-15} y1={-50} y2={-50} />
        <Circle c={colours[3]} x1={6} x2={2} y1={50} y2={100} />
        <Circle c={colours[4]} x1={8} x2={-94} y1={-50} y2={-50} />
        <Circle c={colours[5]} x1={10} x2={9} y1={50} y2={-60} />
        <Circle c={colours[6]} x1={12} x2={1} y1={-50} y2={-89} />
        <Circle c={colours[7]} x1={14} x2={26} y1={50} y2={-20} />
        <Circle c={colours[8]} x1={16} x2={-99} y1={50} y2={-20} />
        <Circle c={colours[9]} x1={16} x2={-99} y1={50} y2={-20} />
      </Box>
    </Button>
  );
}

export default function Page() {
  return (
    <Flex
      bgColor={"blackAlpha.900"}
      width="100vw"
      height="100vh"
      alignItems="center"
      justifyContent="center"
    >
      <Flamegrilled />
    </Flex>
  );
}
