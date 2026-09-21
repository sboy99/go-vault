import { cn } from "@/lib/cn";

export function Input({
  className,
  ...props
}: React.InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      className={cn(
        "h-9 w-full rounded-md border border-control bg-field px-2.5 text-xs text-foreground placeholder:text-foreground-muted transition-colors duration-200 ease-out hover:border-control focus-visible:border-brand-border",
        className,
      )}
      {...props}
    />
  );
}

export function Select({
  className,
  children,
  ...props
}: React.SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <select
      className={cn(
        "h-9 rounded-none border border-dashed border-edge bg-transparent px-2.5 text-xs text-foreground transition-colors duration-200 ease-out hover:border-edge",
        className,
      )}
      {...props}
    >
      {children}
    </select>
  );
}
