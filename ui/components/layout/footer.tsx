const LICENSE_URL =
  "https://github.com/sboy99/go-vault/blob/main/LICENSE";
const GITHUB_URL = "https://github.com/sboy99/go-vault";
const DOCKER_URL =
  "https://hub.docker.com/repository/docker/sboy99/go-vault/general";

const linkClass =
  "text-[13px] text-foreground-muted no-underline transition-colors duration-200 ease-out hover:text-brand";

export function Footer() {
  return (
    <footer className="flex flex-wrap items-center justify-between gap-2 border-t border-dashed border-edge py-4">
      <a
        href={LICENSE_URL}
        target="_blank"
        rel="noopener noreferrer"
        className={linkClass}
      >
        MIT License
      </a>

      <div className="flex items-center gap-3">
        <a
          href={GITHUB_URL}
          target="_blank"
          rel="noopener noreferrer"
          className={linkClass}
        >
          GitHub
        </a>
        <a
          href={DOCKER_URL}
          target="_blank"
          rel="noopener noreferrer"
          className={linkClass}
        >
          Docker
        </a>
      </div>
    </footer>
  );
}
