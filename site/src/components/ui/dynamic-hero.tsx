import { useEffect, useRef, useCallback, useState, type MouseEventHandler } from "react";
import { Play } from "lucide-react";
import { cn } from "@/lib/utils";

type Rgb = { r: number; g: number; b: number };

export interface HeroNavItem {
  id: string;
  label: string;
  href?: string;
  target?: string;
  onClick?: MouseEventHandler<HTMLElement>;
}

export interface HeroSectionProps {
  heading?: string;
  tagline?: string;
  buttonText?: string;
  /** When set, the call-to-action renders as a link instead of a button. */
  buttonHref?: string;
  imageUrl?: string;
  imageAlt?: string;
  videoUrl?: string;
  /** Pass an empty array to hide the built-in nav (e.g. when the page has its own). */
  navItems?: HeroNavItem[];
  className?: string;
}

// Parses 'rgb(r, g, b)' or 'rgba(r, g, b, a)' to {r, g, b}.
const parseRgbColor = (colorString: string): Rgb | null => {
  const match = colorString.match(/rgba?\((\d+),\s*(\d+),\s*(\d+)(?:,\s*[\d.]+)?\)/);
  return match
    ? { r: parseInt(match[1], 10), g: parseInt(match[2], 10), b: parseInt(match[3], 10) }
    : null;
};

const defaultNavItems: HeroNavItem[] = [
  { id: "home", label: "Home", href: "#" },
  { id: "about", label: "About", href: "#about-section" },
  { id: "pricing", label: "Pricing", href: "#pricing" },
  { id: "get-started-nav", label: "Get Started", href: "#get-started" },
];

const navItemClass =
  "py-2 px-3 sm:px-4 rounded-md text-muted-foreground hover:text-foreground hover:bg-accent/10 dark:hover:bg-accent/20 focus:outline-none focus:ring-2 focus:ring-ring transition-colors duration-150 ease-in-out whitespace-nowrap";

const ctaClass =
  "inline-flex py-2 px-4 rounded-xl border border-foreground/50 hover:border-foreground/80 text-foreground transition-colors focus:outline-none focus:ring-2 focus:ring-ring";

const HeroSection = ({
  heading = "Something you really want",
  tagline = "You can't live without this product. I'm sure of it.",
  buttonText = "Get Started",
  buttonHref,
  imageUrl,
  imageAlt = "",
  videoUrl,
  navItems = defaultNavItems,
  className,
}: HeroSectionProps) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const targetRef = useRef<HTMLElement | null>(null);
  const mousePosRef = useRef<{ x: number | null; y: number | null }>({ x: null, y: null });
  const ctxRef = useRef<CanvasRenderingContext2D | null>(null);
  const animationFrameIdRef = useRef<number | null>(null);
  const videoRef = useRef<HTMLVideoElement>(null);
  const [showVideo, setShowVideo] = useState(false);
  const strokeRef = useRef<Rgb>({ r: 128, g: 128, b: 128 });

  // Resolve --foreground to RGB for the canvas, and follow theme changes.
  useEffect(() => {
    const probe = document.createElement("div");
    probe.style.display = "none";
    document.body.appendChild(probe);

    const update = () => {
      probe.style.color = "var(--foreground)";
      const parsed = parseRgbColor(getComputedStyle(probe).color);
      const root = document.documentElement;
      const isDark = root.classList.contains("dark") || root.dataset.theme === "dark";
      strokeRef.current = parsed ?? (isDark ? { r: 250, g: 250, b: 250 } : { r: 10, g: 10, b: 10 });
    };
    update();

    const observer = new MutationObserver(update);
    observer.observe(document.documentElement, { attributes: true, attributeFilter: ["class", "data-theme"] });
    return () => {
      observer.disconnect();
      probe.remove();
    };
  }, []);

  const drawArrow = useCallback(() => {
    const targetEl = targetRef.current;
    const ctx = ctxRef.current;
    const { x: x0, y: y0 } = mousePosRef.current;
    if (!targetEl || !ctx || x0 === null || y0 === null) return;

    const rect = targetEl.getBoundingClientRect();
    const cx = rect.left + rect.width / 2;
    const cy = rect.top + rect.height / 2;

    const a = Math.atan2(cy - y0, cx - x0);
    const x1 = cx - Math.cos(a) * (rect.width / 2 + 12);
    const y1 = cy - Math.sin(a) * (rect.height / 2 + 12);

    const midX = (x0 + x1) / 2;
    const midY = (y0 + y1) / 2;
    const offset = Math.min(200, Math.hypot(x1 - x0, y1 - y0) * 0.5);
    const t = Math.max(-1, Math.min(1, (y0 - y1) / 200));
    const controlX = midX;
    const controlY = midY + offset * t;

    const r = Math.hypot(x1 - x0, y1 - y0);
    const opacity = Math.min(1, (r - Math.max(rect.width, rect.height) / 2) / 500);

    const c = strokeRef.current;
    ctx.strokeStyle = `rgba(${c.r}, ${c.g}, ${c.b}, ${opacity})`;
    ctx.lineWidth = 2;

    // Curve
    ctx.save();
    ctx.beginPath();
    ctx.moveTo(x0, y0);
    ctx.quadraticCurveTo(controlX, controlY, x1, y1);
    ctx.setLineDash([10, 5]);
    ctx.stroke();
    ctx.restore();

    // Arrowhead
    const angle = Math.atan2(y1 - controlY, x1 - controlX);
    const headLength = 10 * (ctx.lineWidth / 1.5);
    ctx.beginPath();
    ctx.moveTo(x1, y1);
    ctx.lineTo(x1 - headLength * Math.cos(angle - Math.PI / 6), y1 - headLength * Math.sin(angle - Math.PI / 6));
    ctx.moveTo(x1, y1);
    ctx.lineTo(x1 - headLength * Math.cos(angle + Math.PI / 6), y1 - headLength * Math.sin(angle + Math.PI / 6));
    ctx.stroke();
  }, []);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas || !targetRef.current) return;
    // Decorative motion: skip it entirely for users who asked for less motion.
    if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) return;

    const ctx = canvas.getContext("2d");
    ctxRef.current = ctx;

    const updateCanvasSize = () => {
      canvas.width = window.innerWidth;
      canvas.height = window.innerHeight;
    };
    const handleMouseMove = (e: MouseEvent) => {
      mousePosRef.current = { x: e.clientX, y: e.clientY };
    };

    window.addEventListener("resize", updateCanvasSize);
    window.addEventListener("mousemove", handleMouseMove);
    updateCanvasSize();

    const animateLoop = () => {
      if (ctx) {
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        drawArrow();
      }
      animationFrameIdRef.current = requestAnimationFrame(animateLoop);
    };
    animateLoop();

    return () => {
      window.removeEventListener("resize", updateCanvasSize);
      window.removeEventListener("mousemove", handleMouseMove);
      if (animationFrameIdRef.current) cancelAnimationFrame(animationFrameIdRef.current);
    };
  }, [drawArrow]);

  useEffect(() => {
    const video = videoRef.current;
    if (!video || !videoUrl) return;

    const handleVideoEnd = () => {
      setShowVideo(false);
      video.currentTime = 0;
    };

    if (showVideo) {
      video.play().catch((error) => {
        console.error("HeroSection: Error playing video:", error);
        setShowVideo(false);
      });
      video.addEventListener("ended", handleVideoEnd);
    } else {
      video.pause();
    }
    return () => video.removeEventListener("ended", handleVideoEnd);
  }, [showVideo, videoUrl]);

  return (
    <div className={cn("bg-background text-foreground min-h-screen flex flex-col", className)}>
      {navItems.length > 0 && (
        <nav className="w-full max-w-screen-md mx-auto flex flex-wrap justify-center sm:justify-between items-center px-4 sm:px-8 py-4 text-sm">
          {navItems.map((item) =>
            item.href ? (
              <a
                key={item.id}
                href={item.href}
                target={item.target}
                rel={item.target === "_blank" ? "noopener noreferrer" : undefined}
                className={navItemClass}
                onClick={item.onClick}
              >
                {item.label}
              </a>
            ) : (
              <button key={item.id} type="button" className={navItemClass} onClick={item.onClick}>
                {item.label}
              </button>
            ),
          )}
        </nav>
      )}

      <div className="flex-grow flex flex-col items-center justify-center">
        <div className="mt-12 sm:mt-16 lg:mt-24 flex flex-col items-center">
          <h1 className="text-3xl sm:text-4xl lg:text-5xl font-medium text-center px-4">{heading}</h1>
          <p className="mt-3 block text-muted-foreground text-center text-base sm:text-lg px-4 max-w-xl">{tagline}</p>
        </div>

        <div className="mt-8 flex justify-center">
          {buttonHref ? (
            <a ref={(el) => { targetRef.current = el; }} href={buttonHref} className={ctaClass}>
              {buttonText}
            </a>
          ) : (
            <button ref={(el) => { targetRef.current = el; }} type="button" className={ctaClass}>
              {buttonText}
            </button>
          )}
        </div>

        <div className="mt-12 lg:mt-16 w-full max-w-screen-sm mx-auto overflow-hidden px-4 sm:px-2">
          <div className="bg-border rounded-[2rem] p-[0.25rem]">
            <div className="relative h-64 sm:h-72 md:h-80 lg:h-96 rounded-[1.75rem] bg-card flex items-center justify-center overflow-hidden">
              {imageUrl && (
                <img
                  src={imageUrl}
                  alt={imageAlt}
                  className={`absolute inset-0 w-full h-full object-cover transition-opacity duration-300 ${
                    showVideo ? "opacity-0 pointer-events-none" : "opacity-100"
                  }`}
                />
              )}
              {videoUrl && (
                <video
                  ref={videoRef}
                  src={videoUrl}
                  muted
                  playsInline
                  className={`w-full h-full object-cover transition-opacity duration-300 ${
                    showVideo ? "opacity-100" : "opacity-0 pointer-events-none"
                  }`}
                />
              )}
              {!showVideo && videoUrl && imageUrl && (
                <button
                  type="button"
                  onClick={() => setShowVideo(true)}
                  className="absolute bottom-3 left-3 sm:bottom-4 sm:left-4 z-20 p-2 sm:p-3 bg-accent/30 hover:bg-accent/50 text-accent-foreground backdrop-blur-sm rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-ring"
                  aria-label="Play video"
                >
                  <Play className="w-4 h-4 sm:w-5 sm:h-5" fill="currentColor" />
                </button>
              )}
              {!imageUrl && !videoUrl && <div className="text-muted-foreground italic">Card Content Area</div>}
            </div>
          </div>
        </div>
      </div>
      <div className="h-12 sm:h-16 md:h-24" />
      <canvas ref={canvasRef} className="fixed inset-0 pointer-events-none z-10" aria-hidden="true" />
    </div>
  );
};

export { HeroSection };
