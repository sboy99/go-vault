export function DemoBanner({ demo }: { demo: boolean }) {
  if (!demo) return null;
  return (
    <div className="mb-4 rounded-md border border-warning/30 bg-warning-muted px-3 py-2 text-sm text-warning">
      Showing demo fixtures — go-vault API is unreachable or returned an error.
    </div>
  );
}
