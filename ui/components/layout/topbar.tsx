"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import {
  CalendarClock,
  DatabaseBackup,
  HardDrive,
  LayoutDashboard,
  ListTodo,
  Menu,
  Settings,
} from "lucide-react";
import { cn } from "@/lib/cn";
import { ModeToggle } from "@/components/theme/mode-toggle";
import { ThemePicker } from "@/components/theme/theme-picker";

const nav = [
  { href: "/", label: "Overview", icon: LayoutDashboard },
  { href: "/backups", label: "Backups", icon: HardDrive },
  { href: "/jobs", label: "Jobs", icon: ListTodo },
  { href: "/schedule", label: "Schedule", icon: CalendarClock },
  { href: "/settings", label: "Settings", icon: Settings },
] as const;

export function Topbar() {
  const pathname = usePathname();
  const [open, setOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;

    const onPointer = (e: MouseEvent) => {
      if (!menuRef.current?.contains(e.target as Node)) setOpen(false);
    };
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") setOpen(false);
    };

    document.addEventListener("mousedown", onPointer);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("mousedown", onPointer);
      document.removeEventListener("keydown", onKey);
    };
  }, [open]);

  return (
    <header className="relative z-40 flex items-center gap-2 border-b border-dashed border-edge py-2 backdrop-blur-xs">
      <div className="relative" ref={menuRef}>
        <button
          type="button"
          aria-label="Menu"
          aria-expanded={open}
          onClick={() => setOpen((v) => !v)}
          className="inline-flex size-8 items-center justify-center rounded-none border border-dashed border-edge bg-transparent text-foreground transition-colors duration-200 ease-out hover:text-brand"
        >
          <Menu className="size-4" strokeWidth={1.75} />
        </button>

        {open ? (
          <nav
            aria-label="Main"
            className="absolute top-full left-0 z-30 mt-1.5 min-w-[11rem] rounded-none border border-dashed border-edge bg-surface p-0"
          >
            <div className="divide-y divide-dashed divide-edge">
              {nav.map((item) => {
                const active =
                  item.href === "/"
                    ? pathname === "/"
                    : pathname === item.href ||
                      pathname.startsWith(`${item.href}/`);
                const Icon = item.icon;
                return (
                  <Link
                    key={item.href}
                    href={item.href}
                    aria-current={active ? "page" : undefined}
                    onClick={() => setOpen(false)}
                    className={cn(
                      "flex items-center gap-2 px-2.5 py-1.5 text-[13px] no-underline transition-colors duration-200 ease-out",
                      active
                        ? "font-medium text-brand"
                        : "text-foreground-light hover:text-foreground",
                    )}
                  >
                    <Icon className="size-3.5 shrink-0" strokeWidth={1.75} />
                    {item.label}
                  </Link>
                );
              })}
            </div>
          </nav>
        ) : null}
      </div>

      <Link
        href="/"
        className="flex items-center gap-2 font-heading text-[15px] font-semibold text-foreground no-underline"
      >
        <DatabaseBackup
          aria-hidden
          className="size-5 shrink-0 text-brand"
          strokeWidth={1.75}
        />
        go-vault
      </Link>

      <div className="ml-auto flex items-center gap-1.5">
        <ThemePicker />
        <ModeToggle />
      </div>
    </header>
  );
}
