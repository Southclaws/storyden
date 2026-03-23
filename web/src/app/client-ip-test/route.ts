import { NextResponse } from "next/server";

import { RequestError } from "@/api/common";
import {
  APIError,
  ClientInfo,
  NetworkHeadersSample,
} from "@/api/openapi-schema";
import { adminSettingsGet } from "@/api/openapi-server/admin";
import { getSession } from "@/api/openapi-server/misc";

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
        error: "client_ip_test_upstream_request_failed",
        message: err.message,
        metadata: {
          status: err.status,
        },
        suggested: "Check API/frontend connectivity and try again.",
      };

      return NextResponse.json(payload, {
        status: err.status,
      });
    }

    const payload: APIError = {
      error: "client_ip_test_request_failed",
      message: "Failed to fetch client IP test data.",
      suggested: "Save settings and retry the test.",
    };

    return NextResponse.json(payload, {
      status: 500,
    });
  }
}
