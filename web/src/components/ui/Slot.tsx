import React from "react";

export function Slot({ children, ...props }: { children?: React.ReactNode }) {
  if (React.Children.count(children) > 1) {
    throw new Error("Only one child allowed");
  }

  if (React.isValidElement(children)) {
    return React.cloneElement(children, {
      ...props,
      ...(children.props as object),
    });
  }

  return null;
}
