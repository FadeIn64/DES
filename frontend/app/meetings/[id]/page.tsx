import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { format } from "date-fns"
import { ru } from "date-fns/locale"
import { Calendar, MapPin, Flag, ExternalLink, Radio, ArrowLeft } from "lucide-react"
import { LiveDriverStats } from "@/components/live-driver-stats"
import { ThemeToggle } from "@/components/theme-toggle"
import Link from "next/link"

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

interface DriverStats {
  position: number
  meeting_key: number
  session_key: number
  driver_number: number
  lap_number: number
  interval: number
  prediction_laps_to_overtake: number
  last_lap_duration: number
  pitsops: number
  last_pit_lap: number
  full_name: string
  abbreviation: string
  team_name: string
  color: string
}

async function getMeeting(id: string): Promise<Meeting> {
  try {
    const res = await fetch(`http://localhost:2112/meetings/${id}`, {
      cache: "no-store",
      next: { revalidate: 0 },
    })

    if (!res.ok) {
      console.error(`Failed to fetch meeting data: ${res.status} ${res.statusText}`)
      throw new Error(`Failed to fetch meeting data: ${res.status}`)
    }

    const data = await res.json()
    console.log("Meeting data:", data)
    return data
  } catch (error) {
    console.error("Error fetching meeting data:", error)
    throw new Error(`Failed to fetch meeting: ${error instanceof Error ? error.message : String(error)}`)
  }
}

async function getDriverStats(id: string): Promise<DriverStats[]> {
  try {
    const res = await fetch(`http://localhost:2112/meetings/${id}/drivers_stats`, {
      cache: "no-store",
      next: { revalidate: 0 },
    })

    if (!res.ok) {
      console.error(`Failed to fetch driver stats: ${res.status} ${res.statusText}`)
      return []
    }

    const data = await res.json()
    console.log("Driver stats data:", data)

    return data.sort((a: DriverStats, b: DriverStats) => a.position - b.position)
  } catch (error) {
    console.error("Error fetching driver stats:", error)
    return []
  }
}

function isRaceLive(meeting: Meeting): boolean {
  const now = new Date()
  const startDate = new Date(meeting.start_date)
  const endDate = new Date(meeting.end_date)

  return now >= startDate && now <= endDate && meeting.dashboard_link.trim() !== ""
}

export default async function MeetingPage({ params }: { params: { id: string } }) {
  try {
    const meeting = await getMeeting(params.id)

    if (!meeting || typeof meeting !== "object") {
      throw new Error("Invalid meeting data received")
    }

    // Fetch driver stats
    const driverStats = await getDriverStats(params.id)

    const startDate = new Date(meeting.start_date)
    const endDate = new Date(meeting.end_date)
    const isLive = isRaceLive(meeting)
    const hasDriverStats = driverStats.length > 0

    return (
        <div className="container mx-auto py-10">
          {/* Navigation Button */}
          <div className="mb-6 flex items-center justify-between">
            <Button asChild variant="outline" className="flex items-center gap-2">
              <Link href="/">
                <ArrowLeft className="w-4 h-4" />
                На главную
              </Link>
            </Button>
            <ThemeToggle />
          </div>

          <Card className="max-w-4xl mx-auto border-slate-200 dark:border-slate-700">
            <CardHeader className="bg-gradient-to-r from-red-50 to-red-100 dark:from-red-950/20 dark:to-red-900/20 border-b border-slate-200 dark:border-slate-700">
              <div className="flex items-start justify-between">
                <div className="space-y-2">
                  <div className="flex items-center gap-3">
                    <CardTitle className="text-3xl md:text-4xl">{meeting.name}</CardTitle>
                    <Badge variant="secondary" className="text-base">
                      {meeting.year}
                    </Badge>
                  </div>
                  <CardDescription className="text-xl">{meeting.description}</CardDescription>
                  {isLive && (
                      <div className="flex items-center gap-2 mt-3">
                        <Badge variant="destructive" className="animate-pulse">
                          <Radio className="w-3 h-3 mr-1" />
                          Гонка сейчас идет
                        </Badge>
                      </div>
                  )}
                </div>
                <div className="bg-red-600 text-white text-2xl font-bold w-16 h-16 flex items-center justify-center rounded-full">
                  <Flag className="w-8 h-8" />
                </div>
              </div>
            </CardHeader>

            <CardContent className="pt-6 space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-4">
                  <div className="flex items-center gap-2 text-lg">
                    <MapPin className="h-5 w-5 text-slate-500 dark:text-slate-400" />
                    <span className="font-medium">Трасса:</span> {meeting.circuit}
                  </div>
                  <div className="flex items-center gap-2 text-lg">
                    <MapPin className="h-5 w-5 text-slate-500 dark:text-slate-400" />
                    <span className="font-medium">Местоположение:</span> {meeting.location}
                  </div>
                  <div className="flex items-center gap-2 text-lg">
                    <Calendar className="h-5 w-5 text-slate-500 dark:text-slate-400" />
                    <span className="font-medium">Начало:</span>
                    {format(startDate, "d MMMM yyyy", { locale: ru })}
                  </div>
                  <div className="flex items-center gap-2 text-lg">
                    <Calendar className="h-5 w-5 text-slate-500 dark:text-slate-400" />
                    <span className="font-medium">Окончание:</span>
                    {format(endDate, "d MMMM yyyy", { locale: ru })}
                  </div>
                </div>

                <div className="space-y-4">
                  <div>
                    <h3 className="font-medium text-xl mb-2">Продолжительность</h3>
                    <p className="text-slate-700 dark:text-slate-300 text-lg">
                      {Math.ceil((endDate.getTime() - startDate.getTime()) / (1000 * 60 * 60 * 24))} дней
                    </p>
                  </div>

                  {isLive && (
                      <div className="bg-red-50 dark:bg-red-950/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
                        <h4 className="font-medium text-red-800 dark:text-red-200 mb-2 flex items-center gap-2">
                          <Radio className="w-4 h-4 animate-pulse" />
                          Прямая трансляция
                        </h4>
                        <p className="text-red-700 dark:text-red-300 text-base mb-3">
                          Гонка проходит прямо сейчас! Не пропустите захватывающие моменты.
                        </p>
                        <Button asChild className="w-full bg-red-600 hover:bg-red-700">
                          <a
                              href={meeting.dashboard_link}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="flex items-center justify-center gap-2"
                          >
                            <ExternalLink className="w-4 h-4" />
                            Следить в прямом эфире
                          </a>
                        </Button>
                      </div>
                  )}

                  {!isLive && meeting.dashboard_link.trim() !== "" && (
                      <div className="bg-slate-50 dark:bg-slate-800/50 border border-slate-200 dark:border-slate-700 rounded-lg p-4">
                        <h4 className="font-medium text-slate-800 dark:text-slate-200 mb-2">Дополнительная информация</h4>
                        <Button asChild variant="outline" className="w-full">
                          <a
                              href={meeting.dashboard_link}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="flex items-center justify-center gap-2"
                          >
                            <ExternalLink className="w-4 h-4" />
                            Перейти к дашборду
                          </a>
                        </Button>
                      </div>
                  )}
                </div>
              </div>
            </CardContent>

            {hasDriverStats && (
                <LiveDriverStats
                    meetingId={params.id}
                    initialData={driverStats}
                    isLive={meeting.dashboard_link.trim() !== ""}
                />
            )}
          </Card>
        </div>
    )
  } catch (error) {
    console.error("Error in MeetingPage:", error)
    throw error
  }
}
