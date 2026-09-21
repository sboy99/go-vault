import { cn } from "@/lib/cn";

export function Card({
  className,
  variant = "default",
  children,
}: {
  className?: string;
  variant?: "default" | "flat";
  children: React.ReactNode;
}) {
  return (
    <div
      className={cn(
        "p-3",
        variant === "default" && "rounded-md border border-edge bg-surface",
        variant === "flat" && "rounded-none bg-transparent shadow-none",
        className,
      )}
    >
      {children}
    </div>
  );
}

export function CardTitle({
  className,
  children,
}: {
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <div
      className={cn(
        "font-heading text-[15px] font-semibold text-foreground",
        className,
      )}
    >
      {children}
    </div>
  );
}
