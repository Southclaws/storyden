"use client";

import {
  Box,
  Center,
  Grid,
  GridItem,
  styled,
  VStack,
} from "@/styled-system/jsx";
import { token } from "@/styled-system/tokens";
import Image from "next/image";
import {
  forwardRef,
  PropsWithChildren,
  ReactElement,
  useLayoutEffect,
  useRef,
  useState,
} from "react";

type label = "l1" | "l2" | "l3" | "l4" | "l5" | "l6";

const labelCopy: Record<label, { n: number; h: string; p: ReactElement }> = {
  l1: {
    n: 1,
    h: "SQLITE OR POSTGRESQL",
    p: (
      <>
        <p>Production-ready whichever you choose.</p>
      </>
    ),
  },
  l2: {
    n: 2,
    h: "FILESYSTEM OR S3",
    p: (
      <>
        <p>Keep it simple or make it scalable.</p>
      </>
    ),
  },
  l3: {
    n: 3,
    h: "STATE OF THE ART SECURITY",
    p: (
      <>
        <p>
          This ain't no PHP spaghetti.
          <br />
          <strong
            aria-hidden="true" /* don't announce this joke, it probably wouldn't make sense or would sound weird to a blind user */
          >
            '); DROP TABLE accounts;
          </strong>
          <br />
          We do things properly.
        </p>
      </>
    ),
  },
  l4: {
    n: 4,
    h: "CONTAINERISED",
    p: (
      <>
        <p>
          No thousand line bash install scripts. Throw on your VPS, Fly.io or if
          you're feeling over-engineery, Kubernetes!
        </p>
      </>
    ),
  },
  l5: {
    n: 5,
    h: "HEADLESS OPTION",
    p: (
      <>
        <p>Despise our design? Build your own frontend WordPress style.</p>
      </>
    ),
  },
  l6: {
    n: 6,
    h: "OPENAPI SPEC",
    p: (
      <>
        <p>
          Fully documented API with a hand-crafted (with love) OpenAPI
          specification.
        </p>
      </>
    ),
  },
};

function getSvgPoint(svgEl: SVGSVGElement, clientX: number, clientY: number) {
  const pt = svgEl.createSVGPoint();
  pt.x = clientX;
  pt.y = clientY;
  const svgPt = pt.matrixTransform(svgEl.getScreenCTM()?.inverse());
  return svgPt;
}

type labelCoords = Record<label, { x: number; y: number }>;

const initial = { x: 0, y: 0 };

export function StorydenComputer() {
  const [ready, setReady] = useState(false);
  const [modal, setModal] = useState<label | null>(null);

  const svgRef = useRef<SVGSVGElement>(null);
  const l1ref = useRef<HTMLDivElement>(null);
  const l2ref = useRef<HTMLDivElement>(null);
  const l3ref = useRef<HTMLDivElement>(null);
  const l4ref = useRef<HTMLDivElement>(null);
  const l5ref = useRef<HTMLDivElement>(null);
  const l6ref = useRef<HTMLDivElement>(null);
  const [test, setTest] = useState<labelCoords>({
    l1: initial,
    l2: initial,
    l3: initial,
    l4: initial,
    l5: initial,
    l6: initial,
  });

  useLayoutEffect(() => {
    if (
      svgRef.current == null ||
      l1ref.current == null ||
      l2ref.current == null ||
      l3ref.current == null ||
      l4ref.current == null ||
      l5ref.current == null ||
      l6ref.current == null
    )
      return;

    const observer = new ResizeObserver(() => {
      if (
        svgRef.current == null ||
        l1ref.current == null ||
        l2ref.current == null ||
        l3ref.current == null ||
        l4ref.current == null ||
        l5ref.current == null ||
        l6ref.current == null
      )
        return;

      const l1 = l1ref.current.getBoundingClientRect();
      const l2 = l2ref.current.getBoundingClientRect();
      const l3 = l3ref.current.getBoundingClientRect();
      const l4 = l4ref.current.getBoundingClientRect();
      const l5 = l5ref.current.getBoundingClientRect();
      const l6 = l6ref.current.getBoundingClientRect();

      const l1p = getSvgPoint(
        svgRef.current,
        l1.left + l1.width,
        l1.top + l1.height / 2
      );
      const l2p = getSvgPoint(
        svgRef.current,
        l2.left + l2.width,
        l2.top + l2.height / 2
      );
      const l3p = getSvgPoint(
        svgRef.current,
        l3.left + l3.width,
        l3.top + l3.height / 2
      );
      const l4p = getSvgPoint(
        svgRef.current,
        l4.left, //
        l4.top + l4.height / 2
      );
      const l5p = getSvgPoint(
        svgRef.current,
        l5.left, //
        l5.top + l5.height / 2
      );
      const l6p = getSvgPoint(
        svgRef.current,
        l6.left, //
        l6.top + l6.height / 2
      );

      setTest({
        l1: { x: l1p.x, y: l1p.y },
        l2: { x: l2p.x, y: l2p.y },
        l3: { x: l3p.x, y: l3p.y },
        l4: { x: l4p.x, y: l4p.y },
        l5: { x: l5p.x, y: l5p.y },
        l6: { x: l6p.x, y: l6p.y },
      });

      setReady(true);
    });

    observer.observe(svgRef.current);
    observer.observe(l1ref.current);
    observer.observe(l2ref.current);
    observer.observe(l3ref.current);
    observer.observe(l4ref.current);
    observer.observe(l5ref.current);
    observer.observe(l6ref.current);

    return () => observer.disconnect();
  }, [svgRef, l1ref]);

  return (
    <Center w="full" position="relative">
      <Image
        src="/landing/SD2000-SEMDEX-HUMAN-COMPUTER-KNOWLEDGE-SYSTEM.png"
        width="1023"
        height="569"
        alt=""
        style={
          {
            //   border: "1px solid blue",
          }
        }
      />

      <Grid
        position="absolute"
        w="full"
        h="full"
        gridTemplateRows="33% auto 33%"
        gridTemplateColumns="35% 1fr 25%"
        style={
          {
            //   border: "1px dotted red",
          }
        }
      >
        <Center gridRow="1">
          <Label ref={l1ref} onClick={() => setModal("l1")}>
            SQLITE&nbsp;OR
            <br />
            POSTGRESQL
          </Label>
        </Center>
        <Center gridRow="2">
          <Label ref={l2ref} onClick={() => setModal("l2")}>
            FILESYSTEM
            <br />
            OR&nbsp;S3
          </Label>
        </Center>
        <Center gridRow="3">
          <Label ref={l3ref} onClick={() => setModal("l3")}>
            STATE&nbsp;OF&nbsp;THE
            <br />
            ART&nbsp;SECURITY
          </Label>
        </Center>
        <Center gridRow="1" gridColumn="3">
          <Label ref={l4ref} onClick={() => setModal("l4")}>
            CONTAINERISED
          </Label>
        </Center>
        <Center gridRow="2" gridColumn="3">
          <Label ref={l5ref} onClick={() => setModal("l5")}>
            HEADLESS OPTION
          </Label>
        </Center>
        <Center gridRow="3" gridColumn="3">
          <Label ref={l6ref} onClick={() => setModal("l6")}>
            OPENAPI
            <br />
            SPEC
          </Label>
        </Center>
      </Grid>

      <svg
        ref={svgRef}
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 920 569"
        preserveAspectRatio="none"
        style={{
          position: "absolute",
          top: 0,
          left: 0,
          width: "100%",
          height: "100%",
          pointerEvents: "none",
          opacity: ready ? 1 : 0,
          transition: "all 1000ms ease",
        }}
      >
        <line
          id="l1"
          x1={test.l1.x}
          y1={test.l1.y}
          x2={350}
          y2={155}
          stroke={token("colors.Shades.newspaper")}
          strokeWidth="2"
        />
        <line
          id="l2"
          x1={test.l2.x}
          y1={test.l2.y}
          x2={440}
          y2={170}
          stroke={token("colors.Shades.newspaper")}
          strokeWidth="2"
        />
        <line
          id="l3"
          x1={test.l3.x}
          y1={test.l3.y}
          x2={393}
          y2={395}
          stroke={token("colors.Shades.newspaper")}
          strokeWidth="2"
        />
        <line
          id="l4"
          x1={test.l4.x}
          y1={test.l4.y}
          x2={600}
          y2={75}
          stroke={token("colors.Shades.newspaper")}
          strokeWidth="2"
        />
        <line
          id="l5"
          x1={test.l5.x}
          y1={test.l5.y}
          x2={690}
          y2={235}
          stroke={token("colors.Shades.newspaper")}
          strokeWidth="2"
        />
        <line
          id="l6"
          x1={test.l6.x}
          y1={test.l6.y}
          x2={690}
          y2={350}
          stroke={token("colors.Shades.newspaper")}
          strokeWidth="2"
        />
      </svg>

      {modal && (
        <Box position="absolute" w="full" h="full">
          <Modal label={modal} onClick={() => setModal(null)} />
        </Box>
      )}
    </Center>
  );
}

const Label = forwardRef<
  HTMLParagraphElement,
  PropsWithChildren & {
    onClick: () => void;
  }
>(({ children, onClick }, ref) => (
  <Box ref={ref} onClick={onClick} p={{ base: "0.5", sm: "1" }}>
    {/* TODO: Use a Button for semantics */}
    <styled.button
      fontFamily="gorton"
      textAlign="center"
      color="Shades.newspaper"
      cursor="pointer"
      fontSize={{
        base: "2xs",
        sm: "xs",
        lg: "sm",
      }}
      width="min"
    >
      {children}
    </styled.button>
  </Box>
));

const cellStyle = { border: "1px solid currentColor", padding: "4px" };

function Modal({ label, onClick }: { label: label; onClick: () => void }) {
  return (
    <Box bgColor="Mono.ink/85" h="full" onClick={onClick}>
      <Center h="full">
        <table
          style={{
            width: "80%",
            backdropFilter: "blur(4px)",
          }}
        >
          <tbody>
            <tr>
              <td style={cellStyle}>
                <styled.aside
                  textAlign="end"
                  fontSize={{
                    base: "2xs",
                    sm: "xs",
                    md: "sm",
                  }}
                >
                  SPEC SHEET {labelCopy[label].n} of 6
                </styled.aside>
              </td>
            </tr>
            <tr>
              <td style={cellStyle}>
                <styled.h1
                  fontSize={{
                    base: "sm",
                    sm: "md",
                    md: "xl",
                    lg: "2xl",
                  }}
                  letterSpacing={8}
                  textAlign="center"
                  fontFamily="gorton"
                >
                  {labelCopy[label].h}
                </styled.h1>
              </td>
            </tr>
            <tr>
              <td style={cellStyle}>
                <styled.p
                  textAlign="center"
                  letterSpacing="widest"
                  fontSize={{
                    base: "sm",
                    sm: "md",
                    md: "lg",
                    lg: "xl",
                  }}
                  textWrap="balance"
                >
                  {labelCopy[label].p}
                </styled.p>
              </td>
            </tr>
          </tbody>
        </table>
      </Center>
    </Box>
  );
}
