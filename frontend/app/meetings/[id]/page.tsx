import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { format } from "date-fns"
import { ru } from "date-fns/locale"
import { Calendar, MapPin, Flag, ExternalLink, Radio } from "lucide-react"

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

    const startDate = new Date(meeting.start_date)
    const endDate = new Date(meeting.end_date)
    const isLive = isRaceLive(meeting)

    return (
      <div className="container mx-auto py-10">
        <Card className="max-w-4xl mx-auto">
          <CardHeader className="bg-gradient-to-r from-red-50 to-red-100 border-b">
            <div className="flex items-start justify-between">
              <div className="space-y-2">
                <div className="flex items-center gap-3">
                  <CardTitle className="text-2xl md:text-3xl">{meeting.name}</CardTitle>
                  <Badge variant="secondary" className="text-sm">
                    {meeting.year}
                  </Badge>
                </div>
                <CardDescription className="text-lg">{meeting.description}</CardDescription>
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
                <div className="flex items-center gap-2">
                  <MapPin className="h-5 w-5 text-slate-500" />
                  <span className="font-medium">Трасса:</span> {meeting.circuit}
                </div>
                <div className="flex items-center gap-2">
                  <MapPin className="h-5 w-5 text-slate-500" />
                  <span className="font-medium">Местоположение:</span> {meeting.location}
                </div>
                <div className="flex items-center gap-2">
                  <Calendar className="h-5 w-5 text-slate-500" />
                  <span className="font-medium">Начало:</span>
                  {format(startDate, "d MMMM yyyy, HH:mm", { locale: ru })}
                </div>
                <div className="flex items-center gap-2">
                  <Calendar className="h-5 w-5 text-slate-500" />
                  <span className="font-medium">Окончание:</span>
                  {format(endDate, "d MMMM yyyy, HH:mm", { locale: ru })}
                </div>
              </div>

              <div className="space-y-4">
                <div>
                  <h3 className="font-medium text-lg mb-2">Продолжительность</h3>
                  <p className="text-slate-700">
                    {Math.ceil((endDate.getTime() - startDate.getTime()) / (1000 * 60 * 60 * 24))} дней
                  </p>
                </div>

                {isLive && (
                  <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                    <h4 className="font-medium text-red-800 mb-2 flex items-center gap-2">
                      <Radio className="w-4 h-4 animate-pulse" />
                      Прямая трансляция
                    </h4>
                    <p className="text-red-700 text-sm mb-3">
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
                  <div className="bg-slate-50 border border-slate-200 rounded-lg p-4">
                    <h4 className="font-medium text-slate-800 mb-2">Дополнительная информация</h4>
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
        </Card>
      </div>
    )
  } catch (error) {
    console.error("Error in MeetingPage:", error)
    throw error
  }
}
