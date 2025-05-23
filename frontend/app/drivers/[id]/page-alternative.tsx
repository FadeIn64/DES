import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { format } from "date-fns"
import { ru } from "date-fns/locale"
import { Flag, User } from "lucide-react"

interface Driver {
  driver_number: number
  team_key: number
  full_name: string
  abbreviation: string
  country: string
  date_of_birth: string
  description: string
}

// Alternative approach using a local API route
async function getDriver(id: string): Promise<Driver> {
  // Use a relative URL to our own API route that proxies the request
  const res = await fetch(`/drivers/${id}/api-route`, {
    cache: "no-store",
  })

  if (!res.ok) {
    throw new Error(`Failed to fetch driver data: ${res.status}`)
  }

  return res.json()
}

export default async function DriverPage({ params }: { params: { id: string } }) {
  const driver = await getDriver(params.id)
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
            <div className="bg-slate-800 text-white text-4xl font-bold w-16 h-16 flex items-center justify-center rounded-full">
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
                <span className="font-medium">Команда ID:</span> {driver.team_key}
              </div>
            </div>
            <div>
              <h3 className="font-medium text-lg mb-2">Описание</h3>
              <p className="text-slate-700">{driver.description}</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
