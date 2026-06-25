"use client";

//
// MIT License
//
// Copyright (c) 2024-2026 pqoqubbw
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// NOTE: This file is adapted from the "Grip" icon in the "lucide-animated" site
// by dmytro: https://github.com/pqoqubbw https://pqoqubbw.dev/ nice work!
//
import type { Variants } from "framer-motion";
import { motion } from "framer-motion";
import {
  forwardRef,
  useCallback,
  useEffect,
  useImperativeHandle,
  useRef,
  useState,
} from "react";

import { Box, BoxProps } from "@/styled-system/jsx";

export interface GripIconHandle {
  startAnimation: () => void;
  stopAnimation: () => void;
}

interface GripProps extends BoxProps {
  size?: number;
}

const CIRCLES = [
  { cx: 19, cy: 5 },
  { cx: 19, cy: 12 },
  { cx: 12, cy: 5 },
  { cx: 19, cy: 19 },
  { cx: 12, cy: 12 },
  { cx: 5, cy: 5 },
  { cx: 12, cy: 19 },
  { cx: 5, cy: 12 },
  { cx: 5, cy: 19 },
];

const VARIANTS: Variants = {
  normal: {
    opacity: 1,
    transition: { duration: 0.25 },
  },
  animate: (index: number) => ({
    opacity: [1, 0.3, 0.3, 1],
    transition: {
      delay: index * 0.07,
      duration: 1.1,
      repeat: Infinity,
      repeatDelay: 0.15,
      times: [0, 0.2, 0.8, 1],
    },
  }),
};

const RobotActivityIcon = forwardRef<GripIconHandle, GripProps>(
  ({ className, size = 28, ...props }, ref) => {
    const [isAnimating, setIsAnimating] = useState(false);
    const isControlledRef = useRef(false);

    const startAnimation = useCallback(() => {
      setIsAnimating(true);
    }, []);

    const stopAnimation = useCallback(() => {
      setIsAnimating(false);
    }, []);

    useImperativeHandle(ref, () => {
      isControlledRef.current = true;
      return { startAnimation, stopAnimation };
    });

    useEffect(() => {
      if (isControlledRef.current) {
        return;
      }

      startAnimation();
    }, []);

    return (
      <Box
        display="inline-flex"
        alignItems="center"
        justifyContent="center"
        className={className}
        {...props}
      >
        <svg
          fill="none"
          height={size}
          stroke="currentColor"
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="2"
          viewBox="0 0 24 24"
          width={size}
          xmlns="http://www.w3.org/2000/svg"
        >
          {CIRCLES.map((circle, index) => (
            <motion.circle
              animate={isAnimating ? "animate" : "normal"}
              custom={index}
              cx={circle.cx}
              cy={circle.cy}
              initial="normal"
              key={`${circle.cx}-${circle.cy}`}
              r="1"
              variants={VARIANTS}
            />
          ))}
        </svg>
      </Box>
    );
  },
);

RobotActivityIcon.displayName = "RobotActivityIcon";

export { RobotActivityIcon };
