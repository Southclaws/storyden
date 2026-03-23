import { NextResponse } from "next/server";

import { RequestError } from "@/api/common";
import { APIError } from "@/api/openapi-schema";
import { getSession } from "@/api/openapi-server/misc";

export async function GET() {
  try {
    const resp = await getSession({
      cache: "no-store",
    });

    return NextResponse.json(
      {
        client: resp.data.client ?? null,
      },
      {
        headers: {
          "Cache-Control": "no-store",
        },
      },
    );
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
