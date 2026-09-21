export const FAVICON_STORAGE_KEY = "govault-favicon";

const FAVICON_SVG = (stroke: string) =>
  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="${stroke}" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M3 12a9 3 0 0 0 5 2.69"/><path d="M21 9.3V5"/><path d="M3 5v14a9 3 0 0 0 6.47 2.88"/><path d="M12 12v4h4"/><path d="M13 20a5 5 0 0 0 9-3 4.5 4.5 0 0 0-4.5-4.5c-1.33 0-2.54.54-3.41 1.41L12 16"/></svg>`;

export function faviconDataUrl(stroke: string): string {
  return `data:image/svg+xml,${encodeURIComponent(FAVICON_SVG(stroke))}`;
}

function isHexColor(value: string): boolean {
  return /^#[0-9a-fA-F]{6}$/.test(value);
}

export function applyFavicon(stroke: string) {
  if (!isHexColor(stroke)) return;

  try {
    localStorage.setItem(FAVICON_STORAGE_KEY, stroke);
  } catch {
    /* ignore quota / private mode */
  }

  const href = faviconDataUrl(stroke);
  const links = document.querySelectorAll<HTMLLinkElement>(
    'link[rel="icon"], link[rel="shortcut icon"]',
  );

  if (links.length === 0) {
    const link = document.createElement("link");
    link.rel = "icon";
    link.type = "image/svg+xml";
    link.href = href;
    document.head.appendChild(link);
    return;
  }

  for (const link of links) {
    link.type = "image/svg+xml";
    link.href = href;
  }
}

export function readStoredFavicon(): string | null {
  try {
    const stored = localStorage.getItem(FAVICON_STORAGE_KEY);
    return stored && isHexColor(stored) ? stored : null;
  } catch {
    return null;
  }
}
