import { NextResponse } from "next/server"

// This is a mock API route that simulates the external API
// You can use this if the localhost API is not accessible from the server
export async function GET(request: Request, { params }: { params: { id: string } }) {
  try {
    // Try to fetch from the actual API first
    const response = await fetch(`http://localhost:2112/drivers/${params.id}`)

    if (response.ok) {
      const data = await response.json()
      return NextResponse.json(data)
    }

    // If the actual API fails, return mock data for driver 81
    if (params.id === "81") {
      return NextResponse.json({
        driver_number: 81,
        team_key: 4,
        full_name: "Oscar Piastri",
        abbreviation: "PIA",
        country: "Australia",
        date_of_birth: "2001-04-06T00:00:00Z",
        description: "Перспективный новичок, показывает стабильные результаты",
      })
    }

    return NextResponse.json({ error: "Driver not found" }, { status: 404 })
  } catch (error) {
    console.error("API route error:", error)
    return NextResponse.json({ error: "Failed to fetch driver data" }, { status: 500 })
  }
}
