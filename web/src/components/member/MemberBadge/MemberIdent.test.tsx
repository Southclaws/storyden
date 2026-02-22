import { render } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { AccountRoleRef, ProfileReference } from "@/api/openapi-schema";
import { token } from "@/styled-system/tokens";

import { MemberIdent, MemberName } from "./MemberIdent";

function role(overrides: Partial<AccountRoleRef> = {}): AccountRoleRef {
  return {
    id: "role_custom",
    name: "Custom",
    colour: "#00aaff",
    badge: false,
    default: false,
    meta: {},
    ...overrides,
  };
}

function profileWithRoles(roles: AccountRoleRef[]): ProfileReference {
  return {
    id: "acc_1",
    handle: "alice",
    name: "Alice",
    joined: new Date().toISOString(),
    roles,
  };
}

describe("MemberIdent decorations", () => {
  it("applies role colour only when metadata.coloured is true", () => {
    const withColour = render(
      <MemberName
        profile={profileWithRoles([
          role({
            colour: "#123456",
            meta: { coloured: true },
          }),
        ])}
        name="full-horizontal"
        size="md"
      />,
    );

    const withColourRoot = withColour.container.querySelector(
      ".member-name__show-horizontal",
    ) as HTMLDivElement;
    expect(withColourRoot.style.getPropertyValue("--colors-color-palette")).toBe(
      "#123456",
    );

    const withoutColour = render(
      <MemberName
        profile={profileWithRoles([
          role({
            colour: "#654321",
            meta: { coloured: false },
          }),
        ])}
        name="full-horizontal"
        size="md"
      />,
    );

    const withoutColourRoot = withoutColour.container.querySelector(
      ".member-name__show-horizontal",
    ) as HTMLDivElement;
    expect(
      withoutColourRoot.style.getPropertyValue("--colors-color-palette"),
    ).toBe(token("colors.fg.default"));
  });

  it("uses bold and italic metadata for decoration vars", () => {
    const { container } = render(
      <MemberName
        profile={profileWithRoles([
          role({
            meta: { coloured: true, bold: true, italic: true },
          }),
        ])}
        name="handle"
        size="md"
      />,
    );

    const root = container.querySelector(
      ".member-name__show-handle",
    ) as HTMLDivElement;
    expect(root.style.getPropertyValue("--decoration-font-style")).toBe("italic");
    expect(root.style.getPropertyValue("--decoration-font-weight")).toBe(
      token("fontWeights.semibold"),
    );
  });

  it("uses default roles when they have decoration metadata", () => {
    const { container } = render(
      <MemberName
        profile={profileWithRoles([
          role({
            id: "00000000000000000a00",
            default: true,
            colour: "#ff0000",
            meta: { coloured: true },
          }),
          role({
            id: "role_custom_2",
            default: false,
            colour: "#00ff00",
            meta: { coloured: true },
          }),
        ])}
        name="full-vertical"
        size="md"
      />,
    );

    const root = container.querySelector(
      ".member-name__show-vertical",
    ) as HTMLDivElement;
    expect(root.style.getPropertyValue("--colors-color-palette")).toBe("#ff0000");
  });

  it("falls back to subtle default for handle variant when no custom roles", () => {
    const { container } = render(
      <MemberName
        profile={profileWithRoles([])}
        name="handle"
        size="md"
      />,
    );

    const root = container.querySelector(
      ".member-name__show-handle",
    ) as HTMLDivElement;
    expect(root.style.getPropertyValue("--colors-color-palette")).toBe(
      token("colors.fg.subtle"),
    );
  });
});

describe("MemberIdent common permutations", () => {
  it("renders handle + avatar variant used in compact lists", () => {
    const profile = profileWithRoles([]);
    const { getByText, getByAltText, queryByText } = render(
      <MemberIdent
        profile={profile}
        size="sm"
        name="handle"
        avatar="visible"
        showRoles="hidden"
      />,
    );

    expect(getByAltText(`${profile.handle}'s avatar`)).toBeInTheDocument();
    expect(getByText(`@${profile.handle}`)).toBeInTheDocument();
    expect(queryByText(profile.name)).not.toBeInTheDocument();
  });

  it("renders full-vertical + avatar variant used on profile/header contexts", () => {
    const profile = profileWithRoles([]);
    const { getByText, getByAltText } = render(
      <MemberIdent
        profile={profile}
        size="lg"
        name="full-vertical"
        avatar="visible"
        showRoles="all"
      />,
    );

    expect(getByAltText(`${profile.handle}'s avatar`)).toBeInTheDocument();
    expect(getByText(profile.name)).toBeInTheDocument();
    expect(getByText(`@${profile.handle}`)).toBeInTheDocument();
  });

  it("renders full-horizontal without avatar when requested", () => {
    const profile = profileWithRoles([]);
    const { getByText, queryByAltText } = render(
      <MemberIdent
        profile={profile}
        size="md"
        name="full-horizontal"
        avatar="hidden"
        showRoles="hidden"
      />,
    );

    expect(queryByAltText(`${profile.handle}'s avatar`)).not.toBeInTheDocument();
    expect(getByText(profile.name)).toBeInTheDocument();
    expect(getByText(`@${profile.handle}`)).toBeInTheDocument();
  });
});
