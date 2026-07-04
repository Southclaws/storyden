"use client";

import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { parseAsInteger, useQueryState } from "nuqs";

import type { Account, ProfileReference } from "@/api/openapi-schema";
import { InvitedByFilter } from "@/components/library/members/MemberFilters/InvitedByFilter";
import { JoinedDateFilter } from "@/components/library/members/MemberFilters/JoinedDateFilter";
import { RoleFilter } from "@/components/library/members/MemberFilters/RoleFilter";
import { SortMenu } from "@/components/library/members/MemberFilters/SortMenu";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { MemberIdent } from "@/components/member/MemberBadge/MemberIdent";
import { EmptyState } from "@/components/site/EmptyState";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { Timestamp } from "@/components/site/Timestamp";
import { Unready } from "@/components/site/Unready";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Group } from "@/components/ui/group";
import { CheckIcon } from "@/components/ui/icons/Check";
import { MembersIcon } from "@/components/ui/icons/Members";
import { LinkButton } from "@/components/ui/link-button";
import * as Menu from "@/components/ui/menu";
import {
  Box,
  CardBox,
  Flex,
  Grid,
  HStack,
  LStack,
  VStack,
  WStack,
  styled,
} from "@/styled-system/jsx";

import { SearchInput } from "../SearchInput";

import { Props, useAdminMemberIndexScreen } from "./useAdminMemberIndexScreen";

export function AdminMemberIndexScreen(props: Props) {
  const { data, error, filters, setFilters } = useAdminMemberIndexScreen(props);
  const searchParams = useSearchParams();
  const [page] = useQueryState(
    "page",
    parseAsInteger.withDefault(props.page ?? 1),
  );
  const currentParams = Object.fromEntries(searchParams.entries());

  if (!data) {
    return <Unready error={error} />;
  }

  async function updateSingle(name: "admin" | "suspended", value: string) {
    if (value === "any") {
      await setFilters({ [name]: null } as never);
      return;
    }

    await setFilters({ [name]: value === "true" } as never);
  }

  async function updateKind(value: string) {
    await setFilters({ kind: value === "any" ? null : value });
  }

  return (
    <VStack alignItems="stretch" gap="4" w="full">
      <Group>
        <SearchInput
          index="/m"
          initialQuery={props.query}
          placeholder="Search members by name, handle, ID, or email"
        />

        <LinkButton size="md" variant="subtle" bg="bg.warning" href="/m">
          <MembersIcon />
        </LinkButton>
      </Group>

      <VStack w="full" gap="2" alignItems="stretch">
        <Flex
          w="full"
          gap="2"
          flexDir={{
            base: "column",
            md: "row",
          }}
        >
          <Box flex="1" minW="0">
            <RoleFilter />
          </Box>
          <Box flex="1" minW="0">
            <InvitedByFilter />
          </Box>
        </Flex>

        <Flex w="full" gap="2" flexDir="row" flexWrap="wrap" alignItems="start">
          <JoinedDateFilter />
          <SortMenu />
          <BooleanFilterMenu
            label={
              filters.admin === true
                ? "Admins only"
                : filters.admin === false
                  ? "Non-admins"
                  : "Admin status"
            }
            options={[
              { value: "any", label: "Any" },
              { value: "true", label: "Admins only" },
              { value: "false", label: "Non-admins" },
            ]}
            selected={String(filters.admin ?? "any")}
            onSelect={(value) => updateSingle("admin", value)}
          />
          <BooleanFilterMenu
            label={
              filters.suspended === true
                ? "Suspended only"
                : filters.suspended === false
                  ? "Active only"
                  : "Suspended state"
            }
            options={[
              { value: "any", label: "Any" },
              { value: "true", label: "Suspended" },
              { value: "false", label: "Active" },
            ]}
            selected={String(filters.suspended ?? "any")}
            onSelect={(value) => updateSingle("suspended", value)}
          />
          <BooleanFilterMenu
            label={
              filters.kind === "human"
                ? "Humans"
                : filters.kind === "bot"
                  ? "Bots"
                  : "Account type"
            }
            options={[
              { value: "any", label: "Any" },
              { value: "human", label: "Humans" },
              { value: "bot", label: "Bots" },
            ]}
            selected={filters.kind ?? "any"}
            onSelect={updateKind}
          />
        </Flex>
      </VStack>

      {data.accounts.length === 0 ? (
        <EmptyState>No accounts matched the current admin filters.</EmptyState>
      ) : (
        <VStack gap="4" alignItems="stretch">
          {data.accounts.map((account) => {
            const authServices = dedupe(account.auth_services);

            return (
              <CardBox key={account.id}>
                <VStack alignItems="stretch" gap="3">
                  <WStack
                    justifyContent="space-between"
                    alignItems="start"
                    gap="3"
                    flexWrap="wrap"
                  >
                    <HStack
                      gap="2"
                      flexWrap="wrap"
                      minW="0"
                      alignItems="center"
                    >
                      <Link
                        href={`/m/${account.handle}`}
                        style={{ minWidth: 0 }}
                      >
                        <MemberIdent
                          profile={asProfileReference(account)}
                          name="full-horizontal"
                          size="md"
                        />
                      </Link>
                      {account.admin && (
                        <Badge colorPalette="orange">admin</Badge>
                      )}
                      {account.suspended ? (
                        <Badge colorPalette="red">suspended</Badge>
                      ) : (
                        <Badge colorPalette="green">active</Badge>
                      )}
                      <Badge variant="outline">{account.verified_status}</Badge>
                    </HStack>

                    <Timestamp
                      created={account.joined}
                      color="fg.subtle"
                      large
                    />
                  </WStack>

                  <styled.code color="fg.muted" fontSize="xs" fontFamily="mono">
                    {account.id}
                  </styled.code>

                  <Grid
                    gridTemplateColumns={{
                      base: "1fr",
                      lg: "repeat(2, minmax(0, 1fr))",
                    }}
                    gap="3"
                  >
                    <VStack alignItems="stretch" gap="3">
                      <InfoBlock label="Emails">
                        <VStack alignItems="stretch" gap="1.5">
                          {account.email_addresses.length === 0 ? (
                            <styled.span color="fg.muted" fontSize="sm">
                              No email addresses
                            </styled.span>
                          ) : (
                            account.email_addresses.map((email) => (
                              <HStack
                                key={email.id}
                                justifyContent="space-between"
                                gap="2"
                                flexWrap="wrap"
                              >
                                <styled.span fontFamily="mono" fontSize="sm">
                                  {email.email_address}
                                </styled.span>
                                <Badge variant="outline">
                                  {email.verified ? "verified" : "unverified"}
                                </Badge>
                              </HStack>
                            ))
                          )}
                        </VStack>
                      </InfoBlock>

                      <InfoBlock label="Roles">
                        <HStack gap="2" flexWrap="wrap">
                          {account.roles.length === 0 ? (
                            <styled.span color="fg.muted" fontSize="sm">
                              No roles
                            </styled.span>
                          ) : (
                            account.roles.map((role) => (
                              <Badge key={role.id} variant="subtle">
                                {role.name}
                              </Badge>
                            ))
                          )}
                        </HStack>
                      </InfoBlock>
                    </VStack>

                    <VStack alignItems="stretch" gap="3">
                      <InfoBlock label="Auth services">
                        <HStack gap="2" flexWrap="wrap">
                          {authServices.length === 0 ? (
                            <styled.span color="fg.muted" fontSize="sm">
                              None
                            </styled.span>
                          ) : (
                            <>
                              {authServices.slice(0, 5).map((service) => (
                                <Badge key={service} variant="outline">
                                  {service}
                                </Badge>
                              ))}
                              {authServices.length > 5 && (
                                <Badge variant="subtle" colorPalette="gray">
                                  +{authServices.length - 5} more
                                </Badge>
                              )}
                            </>
                          )}
                        </HStack>
                      </InfoBlock>

                      <InfoBlock label="Invitation">
                        {account.invited_by ? (
                          <MemberBadge
                            profile={account.invited_by}
                            name="handle"
                            size="sm"
                            avatar="hidden"
                          />
                        ) : (
                          <styled.span color="fg.subtle" fontStyle="italic">
                            n/a
                          </styled.span>
                        )}
                      </InfoBlock>
                    </VStack>
                  </Grid>
                </VStack>
              </CardBox>
            );
          })}
        </VStack>
      )}

      <PaginationControls
        path="/m"
        params={currentParams}
        currentPage={page}
        totalPages={data.total_pages}
        pageSize={data.page_size}
      />
    </VStack>
  );
}

function BooleanFilterMenu(props: {
  label: string;
  selected: string;
  options: { value: string; label: string }[];
  onSelect: (value: string) => Promise<void>;
}) {
  return (
    <Menu.Root positioning={{ placement: "bottom-start" }} lazyMount>
      <Menu.Trigger asChild>
        <Button variant="subtle" size="sm">
          {props.label}
        </Button>
      </Menu.Trigger>
      <Menu.Positioner>
        <Menu.Content minW="44">
          <Menu.ItemGroup id={props.label}>
            {props.options.map((option) => (
              <Menu.Item
                key={option.value}
                value={option.value}
                onClick={() => props.onSelect(option.value)}
              >
                <HStack justify="space-between" w="full">
                  <span>{option.label}</span>
                  {props.selected === option.value && <CheckIcon />}
                </HStack>
              </Menu.Item>
            ))}
          </Menu.ItemGroup>
        </Menu.Content>
      </Menu.Positioner>
    </Menu.Root>
  );
}

function InfoBlock(props: React.PropsWithChildren<{ label: string }>) {
  return (
    <LStack gap="2" alignItems="stretch">
      <styled.span fontSize="sm" color="fg.subtle" fontWeight="medium">
        {props.label}
      </styled.span>
      {props.children}
    </LStack>
  );
}

function dedupe(values: string[]): string[] {
  return [...new Set(values)];
}

function asProfileReference(account: Account): ProfileReference {
  return {
    id: account.id,
    handle: account.handle,
    joined: account.joined,
    name: account.name,
    roles: account.roles,
    signature: account.signature,
    suspended: account.suspended,
  };
}
