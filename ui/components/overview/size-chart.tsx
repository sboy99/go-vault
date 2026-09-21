"use client";

import { useId } from "react";
import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import type { DailyBackupPoint } from "@/lib/api/types";
import { formatBytes } from "@/lib/format";

export function SizeChart({ series }: { series: DailyBackupPoint[] }) {
  const fadeId = useId().replace(/:/g, "");
  const data = series.map((p) => ({
    date: p.date.slice(5),
    size: p.size_bytes,
  }));

  return (
    <div className="h-44 w-full">
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={data} margin={{ top: 8, right: 8, left: 0, bottom: 0 }}>
          <defs>
            <linearGradient id={fadeId} x1="0" y1="0" x2="0" y2="1">
              <stop
                offset="0%"
                stopColor="var(--color-chart-1)"
                stopOpacity={0.4}
              />
              <stop
                offset="100%"
                stopColor="var(--color-chart-1)"
                stopOpacity={0}
              />
            </linearGradient>
          </defs>
          <CartesianGrid
            stroke="var(--color-chart-grid)"
            strokeDasharray="3 3"
            vertical={false}
          />
          <XAxis
            dataKey="date"
            tick={{ fontSize: 11, fill: "var(--color-foreground-muted)" }}
            axisLine={{ stroke: "var(--color-edge)" }}
            tickLine={false}
            minTickGap={24}
          />
          <YAxis
            width={56}
            tick={{ fontSize: 11, fill: "var(--color-foreground-muted)" }}
            axisLine={false}
            tickLine={false}
            tickFormatter={(v: number) => formatBytes(v)}
          />
          <Tooltip
            contentStyle={{
              background: "var(--color-surface-raised)",
              border: "1px solid var(--color-edge)",
              borderRadius: 4,
              color: "var(--color-foreground)",
              fontSize: 12,
            }}
            formatter={(value) => [formatBytes(Number(value ?? 0)), "Size"]}
          />
          <Area
            type="monotone"
            dataKey="size"
            stroke="var(--color-chart-1)"
            fill={`url(#${fadeId})`}
            strokeWidth={1.5}
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
