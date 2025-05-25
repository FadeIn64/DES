import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { format } from "date-fns"
import { ru } from "date-fns/locale"
import { Calendar, MapPin, Flag, Radio, Home } from "lucide-react"
import Link from "next/link"
import { YearFilter } from "@/components/year-filter"
import { SortFilter } from "@/components/sort-filter"
import { ThemeToggle } from "@/components/theme-toggle"
import { getCountryFlag } from "@/utils/country-flags"

interface Meeting {
  meeting_key: number
  name: string
  description: string
  circuit: string
  location: string
  start_date: string
  end_date: string
  year: number
  dashboard_link: string
}

async function getAllMeetings(): Promise<Meeting[]> {
  try {
    const res = await fetch(`http://localhost:2112/meetings`, {
      cache: "no-store",
      next: { revalidate: 0 },
    })

    if (!res.ok) {
      console.error(`Failed to fetch meetings data: ${res.status} ${res.statusText}`)
      throw new Error(`Failed to fetch meetings data: ${res.status}`)
    }

    const data = await res.json()
    console.log("Meetings data:", data)

    // Return data without sorting - sorting will be handled by component
    return data
  } catch (error) {
    console.error("Error fetching meetings data:", error)
    throw new Error(`Failed to fetch meetings: ${error instanceof Error ? error.message : String(error)}`)
  }
}

function isRaceLive(meeting: Meeting): boolean {
  const now = new Date()
  const startDate = new Date(meeting.start_date)
  const endDate = new Date(meeting.end_date)

  return now >= startDate && now <= endDate && meeting.dashboard_link.trim() !== ""
}

function MeetingCard({ meeting }: { meeting: Meeting }) {
  const startDate = new Date(meeting.start_date)
  const endDate = new Date(meeting.end_date)
  const isLive = isRaceLive(meeting)
  const countryFlag = getCountryFlag(meeting.location)

  return (
    <Link href={`/meetings/${meeting.meeting_key}`}>
      <Card className="h-full hover:shadow-lg dark:hover:shadow-xl transition-shadow cursor-pointer group border-slate-200 dark:border-slate-700">
        <CardHeader className="pb-3">
          <div className="flex items-start justify-between">
            <div className="space-y-1 flex-1">
              <div className="flex items-center gap-2 flex-wrap">
                <CardTitle className="text-lg group-hover:text-red-600 dark:group-hover:text-red-400 transition-colors">
                  {meeting.name}
                </CardTitle>
                <Badge variant="secondary" className="text-xs">
                  {meeting.year}
                </Badge>
                {isLive && (
                  <Badge variant="destructive" className="animate-pulse text-xs">
                    <Radio className="w-3 h-3 mr-1" />
                    LIVE
                  </Badge>
                )}
              </div>
              <CardDescription className="text-sm line-clamp-2">{meeting.description}</CardDescription>
            </div>
            <div className="bg-red-600 text-white text-lg font-bold w-12 h-12 flex items-center justify-center rounded-full ml-3">
              <Flag className="w-6 h-6" />
            </div>
          </div>
        </CardHeader>

        <CardContent className="pt-0">
          <div className="space-y-2 text-sm">
            <div className="flex items-center gap-2 text-slate-600 dark:text-slate-400">
              <MapPin className="h-4 w-4" />
              <span className="truncate">{meeting.circuit}</span>
            </div>
            <div className="flex items-center gap-2 text-slate-600 dark:text-slate-400">
              <span className="text-lg">{countryFlag}</span>
              <span className="truncate">{meeting.location}</span>
            </div>
            <div className="flex items-center gap-2 text-slate-600 dark:text-slate-400">
              <Calendar className="h-4 w-4" />
              <span>{format(startDate, "d MMM yyyy", { locale: ru })}</span>
              <span className="text-slate-400 dark:text-slate-500">—</span>
              <span>{format(endDate, "d MMM yyyy", { locale: ru })}</span>
            </div>
          </div>
        </CardContent>
      </Card>
    </Link>
  )
}

export default async function HomePage({
  searchParams,
}: {
  searchParams: { year?: string; sort?: string }
}) {
  try {
    const meetings = await getAllMeetings()

    if (!meetings || !Array.isArray(meetings)) {
      throw new Error("Invalid meetings data received")
    }

    // Get unique years for filter
    const availableYears = [...new Set(meetings.map((meeting) => meeting.year))].sort((a, b) => b - a)

    // Filter by year if specified
    const selectedYear = searchParams.year ? Number.parseInt(searchParams.year) : null
    let filteredMeetings = selectedYear ? meetings.filter((meeting) => meeting.year === selectedYear) : meetings

    // Sort meetings based on sort parameter
    const sortOrder = searchParams.sort || "desc" // default to descending (newest first)
    filteredMeetings = filteredMeetings.sort((a: Meeting, b: Meeting) => {
      const dateA = new Date(a.start_date).getTime()
      const dateB = new Date(b.start_date).getTime()

      return sortOrder === "asc" ? dateA - dateB : dateB - dateA
    })

    return (
      <div className="container mx-auto py-10">
        {/* Header with theme toggle */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-3xl md:text-4xl font-bold mb-2">Гонки Формулы-1</h1>
            <p className="text-slate-600 dark:text-slate-400 mb-6">
              Все гонки чемпионата мира по Формуле-1. Выберите год для фильтрации и настройте сортировку.
            </p>
          </div>
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

        <div className="mb-8">
          <div className="flex flex-col sm:flex-row gap-4 items-start">
            <div>
              <label className="text-sm font-medium text-slate-700 dark:text-slate-300 mb-2 block">
                Фильтр по году:
              </label>
              <YearFilter availableYears={availableYears} selectedYear={selectedYear} />
            </div>
            <div>
              <label className="text-sm font-medium text-slate-700 dark:text-slate-300 mb-2 block">
                Сортировка по дате:
              </label>
              <SortFilter currentSort={sortOrder} />
            </div>
          </div>
        </div>

        {filteredMeetings.length === 0 ? (
          <Card className="max-w-md mx-auto">
            <CardContent className="pt-6 text-center">
              <Flag className="w-12 h-12 text-slate-400 dark:text-slate-500 mx-auto mb-4" />
              <h3 className="text-lg font-medium mb-2">Гонки не найдены</h3>
              <p className="text-slate-600 dark:text-slate-400">
                {selectedYear ? `Нет гонок за ${selectedYear} год` : "Нет доступных гонок"}
              </p>
            </CardContent>
          </Card>
        ) : (
          <>
            <div className="mb-4 text-sm text-slate-600 dark:text-slate-400">
              Найдено гонок: {filteredMeetings.length}
              {selectedYear && ` за ${selectedYear} год`}
              {sortOrder === "asc" ? " (сначала старые)" : " (сначала новые)"}
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredMeetings.map((meeting) => (
                <MeetingCard key={meeting.meeting_key} meeting={meeting} />
              ))}
            </div>
          </>
        )}
      </div>
    )
  } catch (error) {
    console.error("Error in HomePage:", error)
    throw error
  }
}
