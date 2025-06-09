import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { format } from "date-fns"
import { ru } from "date-fns/locale"
import { Flag, User, Users, ArrowLeft, Trophy, Medal, Award, Calendar, MapPin, Radio } from "lucide-react"
import Link from "next/link"
import { ThemeToggle } from "@/components/theme-toggle"
import { getCountryFlag } from "@/utils/country-flags"

interface Driver {
  driver_number: number
  team_key: number
  full_name: string
  abbreviation: string
  country: string
  date_of_birth: string
  description: string
}

interface Team {
  team_key: number
  name: string
  description: string
  color: string
  country: string
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

async function getDriver(id: string): Promise<Driver> {
  try {
    const res = await fetch(`http://localhost:2112/drivers/${id}`, {
      cache: "no-store",
      next: { revalidate: 0 },
    })

    if (!res.ok) {
      console.error(`Failed to fetch driver data: ${res.status} ${res.statusText}`)
      throw new Error(`Failed to fetch driver data: ${res.status}`)
    }

    const data = await res.json()
    console.log("Driver data:", data)
    return data
  } catch (error) {
    console.error("Error fetching driver data:", error)
    throw new Error(`Failed to fetch driver: ${error instanceof Error ? error.message : String(error)}`)
  }
}

async function getTeam(teamId: number): Promise<Team> {
  try {
    const res = await fetch(`http://localhost:2112/teams/${teamId}`, {
      cache: "no-store",
      next: { revalidate: 0 },
    })

    if (!res.ok) {
      console.error(`Failed to fetch team data: ${res.status} ${res.statusText}`)
      throw new Error(`Failed to fetch team data: ${res.status}`)
    }

    const data = await res.json()
    console.log("Team data:", data)
    return data
  } catch (error) {
    console.error("Error fetching team data:", error)
    throw new Error(`Failed to fetch team: ${error instanceof Error ? error.message : String(error)}`)
  }
}

async function getDriverStats(id: string): Promise<DriverStats[]> {
  try {
    const res = await fetch(`http://localhost:2112/drivers/${id}/stats`, {
      cache: "no-store",
      next: { revalidate: 0 },
    })

    if (!res.ok) {
      console.error(`Failed to fetch driver stats: ${res.status} ${res.statusText}`)
      return []
    }

    const data = await res.json()
    console.log("Driver stats data:", data)
    return data
  } catch (error) {
    console.error("Error fetching driver stats:", error)
    return []
  }
}

async function getMeeting(meetingKey: number): Promise<Meeting | null> {
  try {
    const res = await fetch(`http://localhost:2112/meetings/${meetingKey}`, {
      cache: "no-store",
      next: { revalidate: 0 },
    })

    if (!res.ok) {
      console.error(`Failed to fetch meeting data: ${res.status} ${res.statusText}`)
      return null
    }

    const data = await res.json()
    return data
  } catch (error) {
    console.error("Error fetching meeting data:", error)
    return null
  }
}

function getPositionIcon(position: number) {
  switch (position) {
    case 1:
      return <Trophy className="w-5 h-5 text-yellow-500" />
    case 2:
      return <Medal className="w-5 h-5 text-gray-400" />
    case 3:
      return <Award className="w-5 h-5 text-amber-600" />
    default:
      return null
  }
}

function getPositionStyle(position: number) {
  switch (position) {
    case 1:
      return "text-yellow-600 font-bold bg-yellow-50 dark:bg-yellow-950/20 border-yellow-200 dark:border-yellow-800"
    case 2:
      return "text-gray-600 font-bold bg-gray-50 dark:bg-gray-950/20 border-gray-200 dark:border-gray-800"
    case 3:
      return "text-amber-600 font-bold bg-amber-50 dark:bg-amber-950/20 border-amber-200 dark:border-amber-800"
    default:
      return "text-slate-700 dark:text-slate-300 bg-slate-50 dark:bg-slate-800/50 border-slate-200 dark:border-slate-700"
  }
}

function isRaceLive(meeting: Meeting): boolean {
  return meeting.dashboard_link.trim() !== ""
}

export default async function DriverPage({ params }: { params: { id: string } }) {
  try {
    const driver = await getDriver(params.id)

    if (!driver || typeof driver !== "object") {
      throw new Error("Invalid driver data received")
    }

    // Fetch team data and driver stats
    const [team, driverStats] = await Promise.all([getTeam(driver.team_key), getDriverStats(params.id)])

    // Fetch meeting data for each race result
    const raceResults = await Promise.all(
        driverStats.map(async (stat) => {
          const meeting = await getMeeting(stat.meeting_key)
          return {
            ...stat,
            meeting,
          }
        }),
    )

    // Sort by meeting end date (latest first - inverted)
    const sortedResults = raceResults
        .filter((result) => result.meeting !== null)
        .sort((a, b) => {
          if (!a.meeting || !b.meeting) return 0
          return new Date(b.meeting.end_date).getTime() - new Date(a.meeting.end_date).getTime()
        })

    const birthDate = new Date(driver.date_of_birth)

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

          <div className="space-y-6">
            {/* Driver Info Card */}
            <Card className="max-w-3xl mx-auto border-slate-200 dark:border-slate-700">
              <CardHeader className="bg-slate-50 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-700">
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle className="text-3xl md:text-4xl">{driver.full_name}</CardTitle>
                    <CardDescription className="text-xl">
                      Номер: {driver.driver_number} | Аббревиатура: {driver.abbreviation}
                    </CardDescription>
                  </div>
                  <div
                      className="text-white text-4xl font-bold w-16 h-16 flex items-center justify-center rounded-full"
                      style={{ backgroundColor: team.color }}
                  >
                    {driver.driver_number}
                  </div>
                </div>
              </CardHeader>
              <CardContent className="pt-6 space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-4">
                    <div className="flex items-center gap-2 text-lg">
                      <Flag className="h-5 w-5 text-slate-500 dark:text-slate-400" />
                      <span className="font-medium">Страна:</span> {driver.country}
                    </div>
                    <div className="flex items-center gap-2 text-lg">
                      <User className="h-5 w-5 text-slate-500 dark:text-slate-400" />
                      <span className="font-medium">Дата рождения:</span>{" "}
                      {format(birthDate, "d MMMM yyyy", { locale: ru })}
                    </div>
                    <div className="flex items-center gap-2 text-lg">
                      <Users className="h-5 w-5 text-slate-500 dark:text-slate-400" />
                      <span className="font-medium">Команда:</span>
                      <span style={{ color: team.color }} className="font-semibold">
                      {team.name}
                    </span>
                    </div>
                  </div>
                  <div>
                    <h3 className="font-medium text-xl mb-2">Описание</h3>
                    <p className="text-slate-700 dark:text-slate-300 mb-4 text-lg">{driver.description}</p>

                    <h4 className="font-medium text-lg mb-2">О команде</h4>
                    <p className="text-slate-600 dark:text-slate-400 text-lg">{team.description}</p>
                    <p className="text-slate-500 dark:text-slate-500 text-sm mt-1">Страна команды: {team.country}</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Race Results Card */}
            {sortedResults.length > 0 && (
                <Card className="max-w-3xl mx-auto border-slate-200 dark:border-slate-700">
                  <CardHeader>
                    <CardTitle className="text-2xl">Результаты последних гонок</CardTitle>
                    <CardDescription>
                      Результаты выступлений в обратном хронологическом порядке ({sortedResults.length} гонок)
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-3">
                      {sortedResults.map((result) => {
                        if (!result.meeting) return null

                        const isLive = isRaceLive(result.meeting)
                        const countryFlag = getCountryFlag(result.meeting.location)

                        return (
                            <Link key={result.meeting_key} href={`/meetings/${result.meeting_key}`}>
                              <div
                                  className={`p-4 rounded-lg border-2 cursor-pointer hover:shadow-md dark:hover:shadow-lg transition-all ${getPositionStyle(result.position)}`}
                              >
                                <div className="flex items-center justify-between">
                                  <div className="flex-1">
                                    <div className="flex items-center gap-3 mb-2">
                                      <div className="flex items-center gap-2">
                                        {getPositionIcon(result.position)}
                                        <span className="text-2xl font-bold">P{result.position}</span>
                                      </div>
                                      {isLive && (
                                          <Badge variant="destructive" className="animate-pulse">
                                            <Radio className="w-3 h-3 mr-1" />
                                            LIVE
                                          </Badge>
                                      )}
                                    </div>

                                    <h4 className="font-semibold text-xl mb-1">{result.meeting.name}</h4>

                                    <div className="flex items-center gap-4 text-base text-slate-600 dark:text-slate-400">
                                      <div className="flex items-center gap-1">
                                        <MapPin className="w-4 h-4" />
                                        <span>{result.meeting.circuit}</span>
                                      </div>
                                      <div className="flex items-center gap-1">
                                        <span className="text-base">{countryFlag}</span>
                                        <span>{result.meeting.location}</span>
                                      </div>
                                      <div className="flex items-center gap-1">
                                        <Calendar className="w-4 h-4" />
                                        <span>{format(new Date(result.meeting.end_date), "d MMM yyyy", { locale: ru })}</span>
                                      </div>
                                    </div>

                                    <div className="flex items-center gap-4 text-sm text-slate-500 dark:text-slate-500 mt-2">
                                      <span>Кругов: {result.lap_number}</span>
                                      <span>Пит-стопов: {result.pitsops}</span>
                                    </div>
                                  </div>
                                </div>
                              </div>
                            </Link>
                        )
                      })}
                    </div>
                  </CardContent>
                </Card>
            )}

            {sortedResults.length === 0 && (
                <Card className="max-w-3xl mx-auto border-slate-200 dark:border-slate-700">
                  <CardContent className="pt-6 text-center">
                    <Trophy className="w-12 h-12 text-slate-400 dark:text-slate-500 mx-auto mb-4" />
                    <h3 className="text-lg font-medium mb-2">Нет данных о результатах</h3>
                    <p className="text-slate-600 dark:text-slate-400">Результаты гонок для этого гонщика пока недоступны</p>
                  </CardContent>
                </Card>
            )}
          </div>
        </div>
    )
  } catch (error) {
    console.error("Error in DriverPage:", error)
    throw error
  }
}
