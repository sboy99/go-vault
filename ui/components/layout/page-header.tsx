export function PageHeader({
  title,
  description,
  actions,
}: {
  title: string;
  description?: string;
  actions?: React.ReactNode;
}) {
  return (
    <div className="mb-3 flex flex-col gap-1.5 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <h1 className="m-0 font-heading text-xl font-semibold text-foreground">
          {title}
        </h1>
        {description ? (
          <p className="mt-0.5 mb-0 max-w-xl text-sm text-foreground-muted">
            {description}
          </p>
        ) : null}
      </div>
      {actions ? (
        <div className="flex flex-wrap items-center gap-2">{actions}</div>
      ) : null}
    </div>
  );
}
