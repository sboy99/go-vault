import { cn } from "@/lib/cn";

type BadgeTone = "success" | "warning" | "error" | "info" | "neutral" | "running";

const tones: Record<BadgeTone, string> = {
  success: "bg-success/10 text-success",
  warning: "bg-warning/10 text-warning",
  error: "bg-destructive/10 text-destructive",
  info: "bg-info/10 text-info",
  neutral: "bg-surface-control text-foreground-muted",
  running: "bg-info/10 text-info",
};

export function Badge({
  tone = "neutral",
  className,
  children,
}: {
  tone?: BadgeTone;
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-md px-1.5 py-1 text-[10px]/3 font-medium",
        tones[tone],
        className,
      )}
    >
      {children}
    </span>
  );
}
