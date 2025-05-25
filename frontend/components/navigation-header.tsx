import { Button } from "@/components/ui/button"
import { ArrowLeft, Home } from "lucide-react"
import { ThemeToggle } from "@/components/theme-toggle"
import Link from "next/link"

interface NavigationHeaderProps {
  showBackButton?: boolean
  backHref?: string
  backLabel?: string
  title?: string
}

export function NavigationHeader({
  showBackButton = true,
  backHref = "/",
  backLabel = "На главную",
  title,
}: NavigationHeaderProps) {
  return (
    <div className="mb-6">
      <div className="flex items-center justify-between">
        {showBackButton && (
          <Button asChild variant="outline" className="flex items-center gap-2">
            <Link href={backHref}>
              <ArrowLeft className="w-4 h-4" />
              {backLabel}
            </Link>
          </Button>
        )}
        {title && <h1 className="text-2xl font-bold text-slate-800 dark:text-slate-200">{title}</h1>}
        <div className="flex items-center gap-2">
          <ThemeToggle />
          <Button asChild variant="ghost" size="sm">
            <Link href="/" className="flex items-center gap-2">
              <Home className="w-4 h-4" />
              Главная
            </Link>
          </Button>
        </div>
      </div>
    </div>
  )
}
