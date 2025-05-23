import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { format } from "date-fns"
import { ru } from "date-fns/locale"
import { Flag, User, Users } from "lucide-react"

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

export default async function DriverPage({ params }: { params: { id: string } }) {
  try {
    const driver = await getDriver(params.id)

    if (!driver || typeof driver !== "object") {
      throw new Error("Invalid driver data received")
    }

    // Fetch team data using the team_key from driver data
    const team = await getTeam(driver.team_key)

    const birthDate = new Date(driver.date_of_birth)

    return (
      <div className="container mx-auto py-10">
        <Card className="max-w-3xl mx-auto">
          <CardHeader className="bg-slate-50 border-b">
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="text-2xl md:text-3xl">{driver.full_name}</CardTitle>
                <CardDescription className="text-lg">
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
                <div className="flex items-center gap-2">
                  <Flag className="h-5 w-5 text-slate-500" />
                  <span className="font-medium">Страна:</span> {driver.country}
                </div>
                <div className="flex items-center gap-2">
                  <User className="h-5 w-5 text-slate-500" />
                  <span className="font-medium">Дата рождения:</span> {format(birthDate, "d MMMM yyyy", { locale: ru })}
                </div>
                <div className="flex items-center gap-2">
                  <Users className="h-5 w-5 text-slate-500" />
                  <span className="font-medium">Команда:</span>
                  <span style={{ color: team.color }} className="font-semibold">
                    {team.name}
                  </span>
                </div>
              </div>
              <div>
                <h3 className="font-medium text-lg mb-2">Описание</h3>
                <p className="text-slate-700 mb-4">{driver.description}</p>

                <h4 className="font-medium text-base mb-2">О команде</h4>
                <p className="text-slate-600 text-sm">{team.description}</p>
                <p className="text-slate-500 text-xs mt-1">Страна команды: {team.country}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  } catch (error) {
    console.error("Error in DriverPage:", error)
    throw error
  }
}
