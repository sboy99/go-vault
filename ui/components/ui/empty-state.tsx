import { cn } from "@/lib/cn";

export function EmptyState({
  title,
  description,
  action,
  className,
}: {
  title: string;
  description?: string;
  action?: React.ReactNode;
  className?: string;
}) {
  return (
    <div
      className={cn(
        "flex flex-col items-center justify-center gap-2 rounded-none border border-dashed border-edge bg-transparent px-6 py-10 text-center",
        className,
      )}
    >
      <p className="font-heading font-medium text-foreground">{title}</p>
      {description ? (
        <p className="max-w-sm text-sm text-foreground-muted">{description}</p>
      ) : null}
      {action}
    </div>
  );
}
