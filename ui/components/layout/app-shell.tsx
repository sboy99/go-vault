import { Footer } from "@/components/layout/footer";
import { Topbar } from "@/components/layout/topbar";
import { ToastProvider } from "@/components/ui/toast";

export function AppShell({ children }: { children: React.ReactNode }) {
  return (
    <ToastProvider>
      <div className="min-h-screen bg-background text-foreground">
        <div className="mx-auto flex min-h-screen w-full max-w-5xl flex-col px-3 md:px-4">
          <Topbar />
          <main className="flex-1 py-3">{children}</main>
          <Footer />
        </div>
      </div>
    </ToastProvider>
  );
}
