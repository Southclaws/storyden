import { useEffect, useRef, useState } from "react";

import { getPluginGetLogsKey } from "@/api/openapi-client/plugins";
import { Plugin } from "@/api/openapi-schema";
import { API_ADDRESS } from "@/config";
import { Box, styled } from "@/styled-system/jsx";

type Props = {
  plugin: Plugin;
};

export function PluginLogViewer({ plugin }: Props) {
  const [logs, setLogs] = useState<string[]>([]);
  const logsEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const [streamURL] = getPluginGetLogsKey(plugin.id);
    const url = `${API_ADDRESS}/api${streamURL}`;
    const eventSource = new EventSource(url, {
      withCredentials: true,
    });

    eventSource.onmessage = (event) => {
      setLogs((prev) => [...prev, event.data]);
    };

    eventSource.addEventListener("end", () => {
      eventSource.close();
    });

    eventSource.onerror = () => {
      eventSource.close();
    };

    return () => {
      eventSource.close();
    };
  }, [plugin.id]);

  useEffect(() => {
    logsEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [logs]);

  return (
    <Box w="full" fontSize="sm" h="96" overflowY="auto">
      <styled.pre h="full" overflowX="scroll">
        {logs.map((log, index) => (
          <styled.p maxW="full" minW="0" key={index}>
            {log}
          </styled.p>
        ))}
        <div ref={logsEndRef} />
      </styled.pre>
    </Box>
  );
}
