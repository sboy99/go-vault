import type { Metadata } from "next";
import { Fira_Code } from "next/font/google";
import { AppShell } from "@/components/layout/app-shell";
import "./globals.css";

const firaCode = Fira_Code({
  subsets: ["latin"],
  variable: "--font-fira-code",
  display: "swap",
});

export const metadata: Metadata = {
  metadataBase: new URL(
    (process.env.NEXT_PUBLIC_APP_URL || "http://localhost:3000").replace(
      /\/$/,
      "",
    ),
  ),
  title: {
    default: "go-vault",
    template: "%s · go-vault",
  },
  description: "PostgreSQL backup control plane",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html
      lang="en"
      suppressHydrationWarning
      data-theme="dark"
      data-primary="emerald"
      data-neutral="neutral"
      className={`dark h-full ${firaCode.variable}`}
    >
      <head>
        <script
          dangerouslySetInnerHTML={{
            __html: `(function(){try{var p=localStorage.getItem("govault-primary");var n=localStorage.getItem("govault-neutral");var t=localStorage.getItem("govault-theme");var f=localStorage.getItem("govault-favicon");var P=["black","red","orange","amber","yellow","lime","green","emerald","teal","cyan","sky","blue","indigo","violet","purple","fuchsia","pink","rose"];var N=["slate","gray","zinc","neutral","stone","taupe","mauve","mist","olive"];var C={black:"#a1a1aa",red:"#ef4444",orange:"#f97316",amber:"#f59e0b",yellow:"#eab308",lime:"#84cc16",green:"#22c55e",emerald:"#10b981",teal:"#14b8a6",cyan:"#06b6d4",sky:"#0ea5e9",blue:"#3b82f6",indigo:"#6366f1",violet:"#8b5cf6",purple:"#a855f7",fuchsia:"#d946ef",pink:"#ec4899",rose:"#f43f5e"};var root=document.documentElement;if(p&&P.indexOf(p)!==-1)root.dataset.primary=p;else p="emerald";if(n&&N.indexOf(n)!==-1)root.dataset.neutral=n;if(t==="light"||t==="dark"){root.dataset.theme=t;root.classList.toggle("dark",t==="dark");}var stroke=(f&&/^#[0-9a-fA-F]{6}$/.test(f))?f:(C[p]||"#10b981");try{localStorage.setItem("govault-favicon",stroke);}catch(e){}var svg='<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="'+stroke+'" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M3 12a9 3 0 0 0 5 2.69"/><path d="M21 9.3V5"/><path d="M3 5v14a9 3 0 0 0 6.47 2.88"/><path d="M12 12v4h4"/><path d="M13 20a5 5 0 0 0 9-3 4.5 4.5 0 0 0-4.5-4.5c-1.33 0-2.54.54-3.41 1.41L12 16"/></svg>';var href="data:image/svg+xml,"+encodeURIComponent(svg);function setIcon(){var links=document.querySelectorAll('link[rel="icon"],link[rel="shortcut icon"]');if(!links.length){var link=document.createElement("link");link.rel="icon";link.type="image/svg+xml";link.href=href;document.head.appendChild(link);return;}for(var i=0;i<links.length;i++){links[i].type="image/svg+xml";links[i].href=href;}}setIcon();document.addEventListener("DOMContentLoaded",setIcon);setTimeout(setIcon,0);setTimeout(setIcon,50);}catch(e){}})();`,
          }}
        />
      </head>
      <body className="min-h-full bg-background font-sans text-foreground antialiased">
        <AppShell>{children}</AppShell>
      </body>
    </html>
  );
}
