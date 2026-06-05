import { NextResponse } from "next/server";

import { RequestError } from "@/api/common";
import {
  APIError,
  ClientInfo,
  NetworkHeadersSample,
} from "@/api/openapi-schema";
import { adminSettingsGet } from "@/api/openapi-server/admin";
import { getSession } from "@/api/openapi-server/misc";

const ProblemJSONMediaType = "application/problem+json";

export async function GET() {
  try {
    const [sessionResp, adminResp] = await Promise.all([
      getSession({
        cache: "no-store",
      }),
      adminSettingsGet({
        cache: "no-store",
      }),
    ]);

    const payload: {
      client: ClientInfo | null;
      headers: NetworkHeadersSample | null;
    } = {
      client: sessionResp.data.client ?? null,
      headers: adminResp.data.headers ?? null,
    };

    return NextResponse.json(payload, {
      headers: {
        "Cache-Control": "no-store",
      },
    });
  } catch (err) {
    if (err instanceof RequestError) {
      const payload: APIError = {
        trace_id: err.problem?.trace_id ?? crypto.randomUUID(),
        type: err.problem?.type ?? "about:blank",
        title: err.problem?.title ?? "Upstream request failed",
        detail: err.problem?.detail ?? err.message,
        metadata: {
          code: "client_ip_test_upstream_request_failed",
          upstream_status: err.status,
        },
      };

      return NextResponse.json(payload, {
        headers: {
          "Content-Type": ProblemJSONMediaType,
        },
        status: err.status,
      });
    }

    const payload: APIError = {
      trace_id: crypto.randomUUID(),
      type: "about:blank",
      title: "Internal Server Error",
      detail: "Failed to fetch client IP test data.",
      metadata: {
        code: "client_ip_test_request_failed",
      },
    };

    return NextResponse.json(payload, {
      headers: {
        "Content-Type": ProblemJSONMediaType,
      },
      status: 500,
    });
  }
}
