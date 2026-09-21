import { cn } from "@/lib/cn";

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary" | "ghost" | "danger";
};

export function Button({
  className,
  variant = "primary",
  type = "button",
  children,
  ...props
}: ButtonProps) {
  return (
    <button
      type={type}
      className={cn(
        "inline-flex items-center justify-center gap-2 border px-2.5 py-1.5 text-sm font-medium transition-[background-color,border-color,color] duration-200 ease-out",
        variant === "primary" &&
          "rounded border-transparent bg-brand-button text-background hover:bg-brand-button-hover",
        variant === "secondary" &&
          "rounded-none border-dashed border-edge bg-transparent text-foreground hover:text-brand",
        variant === "ghost" &&
          "rounded border-transparent bg-transparent text-foreground-light hover:bg-surface-control hover:text-foreground",
        variant === "danger" &&
          "rounded border-destructive/30 bg-destructive text-background hover:bg-destructive-hover",
        className,
      )}
      {...props}
    >
      {children}
    </button>
  );
}
