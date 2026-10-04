import { parseAsInteger, parseAsString, useQueryStates } from "nuqs";
import { useState } from "react";

import { handle } from "@/api/client";
import { RequestError } from "@/api/common";
import {
  adminOAuthRegistrationApprovalBulkSubmit,
  adminOAuthRegistrationApprovalSubmit,
  useAdminOAuthRegistrationApprovalGet,
  useAdminOAuthRegistrationApprovalList,
} from "@/api/openapi-client/admin";
import { deriveError } from "@/utils/error";

export function useRegistrationApprovals(enabled: boolean) {
  const [filters, setFilters] = useQueryStates({
    approvals_page: parseAsInteger.withDefault(1),
    tab: parseAsString,
    verification_code: parseAsString.withDefault(""),
  });
  const [submitting, setSubmitting] = useState<{
    code: string;
    approved: boolean;
  }>();
  const [bulkSubmitting, setBulkSubmitting] = useState(false);
  const [bulkUpdated, setBulkUpdated] = useState<number>();
  const [outcome, setOutcome] = useState<"approved" | "denied">();
  const [submitError, setSubmitError] = useState<string>();
  const queue = useAdminOAuthRegistrationApprovalList(
    { page: filters.approvals_page.toString() },
    { swr: { enabled, refreshInterval: 30_000 } },
  );
  const queuedReview = queue.data?.registrations.find(
    (item) => item.verification_code === filters.verification_code,
  );
  const lookup = useAdminOAuthRegistrationApprovalGet(
    { verification_code: filters.verification_code },
    {
      swr: {
        enabled: enabled && !!filters.verification_code && !queuedReview,
        shouldRetryOnError: false,
        keepPreviousData: false,
      },
    },
  );
  const review = filters.verification_code
    ? (queuedReview ?? (lookup.error ? undefined : lookup.data))
    : undefined;
  const registrations = review
    ? [
        review,
        ...(queue.data?.registrations ?? []).filter(
          (item) => item.verification_code !== review.verification_code,
        ),
      ]
    : (queue.data?.registrations ?? []);

  async function decide(code: string, approved: boolean) {
    if (submitting || bulkSubmitting) return;
    setSubmitting({ code, approved });
    setOutcome(undefined);
    setBulkUpdated(undefined);
    setSubmitError(undefined);
    await handle(
      async () => {
        await adminOAuthRegistrationApprovalSubmit({
          verification_code: code,
          approved,
        });
        setOutcome(approved ? "approved" : "denied");
        if (filters.verification_code === code) {
          await setFilters({ verification_code: null, tab: "reviews" });
        }
        if (
          queue.data?.registrations.length === 1 &&
          filters.approvals_page > 1
        ) {
          await setFilters({ approvals_page: filters.approvals_page - 1 });
        }
        await queue.mutate();
      },
      {
        errorToast: false,
        onError: async (error) => {
          setSubmitError(deriveError(error));
        },
        cleanup: async () => {
          setSubmitting(undefined);
        },
      },
    );
  }

  async function decideAll(approved: boolean, createdBefore: string) {
    if (submitting || bulkSubmitting) return;
    setBulkSubmitting(true);
    setSubmitError(undefined);
    setOutcome(undefined);
    setBulkUpdated(undefined);
    await handle(
      async () => {
        const result = await adminOAuthRegistrationApprovalBulkSubmit({
          approved,
          created_before: createdBefore,
        });
        setBulkUpdated(result.updated);
        setOutcome(approved ? "approved" : "denied");
        await setFilters({
          verification_code: null,
          approvals_page: 1,
          tab: "reviews",
        });
        await queue.mutate();
      },
      {
        errorToast: false,
        onError: async (error) => {
          setSubmitError(deriveError(error));
        },
        cleanup: async () => {
          setBulkSubmitting(false);
        },
      },
    );
  }

  return {
    filters,
    setFilters,
    queue,
    lookup,
    review,
    registrations,
    submitting,
    bulkSubmitting,
    bulkUpdated,
    decideAll,
    outcome,
    submitError,
    dismissOutcome: () => setOutcome(undefined),
    dismissSubmitError: () => setSubmitError(undefined),
    lookupErrorMessage:
      lookup.error instanceof RequestError && lookup.error.status === 400
        ? "This code is invalid, expired, or already reviewed. Check the code or ask the agent to start a new request."
        : deriveError(lookup.error),
    decide,
  };
}
