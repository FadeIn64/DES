"use client"

import { useState, useEffect } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Trophy, Medal, Award, RefreshCw } from "lucide-react"
import Link from "next/link"

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

interface LiveDriverStatsProps {
  meetingId: string
  initialData: DriverStats[]
  isLive: boolean
}

async function fetchDriverStats(meetingId: string): Promise<DriverStats[]> {
  try {
    const res = await fetch(`http://localhost:2112/meetings/${meetingId}/drivers_stats`, {
      cache: "no-store",
    })

    if (!res.ok) {
      console.error(`Failed to fetch driver stats: ${res.status} ${res.statusText}`)
      return []
    }

    const data = await res.json()
    return data.sort((a: DriverStats, b: DriverStats) => a.position - b.position)
  } catch (error) {
    console.error("Error fetching driver stats:", error)
    return []
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
      return "text-yellow-600 font-bold"
    case 2:
      return "text-gray-500 font-bold"
    case 3:
      return "text-amber-600 font-bold"
    default:
      return "text-slate-700 dark:text-slate-300"
  }
}

export function LiveDriverStats({ meetingId, initialData, isLive }: LiveDriverStatsProps) {
  const [driverStats, setDriverStats] = useState<DriverStats[]>(initialData)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [lastUpdated, setLastUpdated] = useState<Date>(new Date())

  useEffect(() => {
    if (!isLive) return

    const interval = setInterval(async () => {
      setIsRefreshing(true)
      try {
        const newData = await fetchDriverStats(meetingId)
        if (newData.length > 0) {
          setDriverStats(newData)
          setLastUpdated(new Date())
        }
      } catch (error) {
        console.error("Error updating driver stats:", error)
      } finally {
        setIsRefreshing(false)
      }
    }, 5000) // Update every 5 seconds

    return () => clearInterval(interval)
  }, [meetingId, isLive])

  if (driverStats.length === 0) {
    return null
  }

  return (
      <div className="border-t border-slate-200 dark:border-slate-700">
        <div className="flex items-center justify-between p-6 pb-3">
          <div>
            <h3 className="text-2xl font-semibold">Позиции гонщиков</h3>
            <p className="text-base text-slate-600 dark:text-slate-400">
              {isLive ? "Данные обновляются в реальном времени" : "Текущие позиции в гонке"}
            </p>
          </div>
          {isLive && (
              <div className="flex items-center gap-2 text-base text-slate-500 dark:text-slate-400">
                <RefreshCw className={`w-4 h-4 ${isRefreshing ? "animate-spin" : ""}`} />
                <span>
              Обновлено:{" "}
                  {lastUpdated.toLocaleTimeString("ru-RU", {
                    hour: "2-digit",
                    minute: "2-digit",
                    second: "2-digit",
                  })}
            </span>
              </div>
          )}
        </div>

        <div className="px-6 pb-6">
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="border-slate-200 dark:border-slate-700">
                  <TableHead className="w-20 text-base">Позиция</TableHead>
                  <TableHead className="w-20 text-base">Номер</TableHead>
                  <TableHead className="text-base">Гонщик</TableHead>
                  <TableHead className="text-base">Команда</TableHead>
                  <TableHead className="text-right text-base">Круг</TableHead>
                  <TableHead className="text-right text-base">Интервал</TableHead>
                  <TableHead className="text-right text-base">Пит-стопы</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {driverStats.map((driver) => (
                    <TableRow
                        key={driver.driver_number}
                        className={`hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors border-slate-200 dark:border-slate-700 ${isRefreshing ? "opacity-75" : ""}`}
                    >
                      <TableCell>
                        <div className="flex items-center gap-2">
                          {getPositionIcon(driver.position)}
                          <span className={getPositionStyle(driver.position)}>{driver.position}</span>
                        </div>
                      </TableCell>
                      <TableCell>
                        <Link href={`/drivers/${driver.driver_number}`}>
                          <div
                              className="w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-bold cursor-pointer hover:scale-110 transition-transform"
                              style={{ backgroundColor: driver.color }}
                          >
                            {driver.driver_number}
                          </div>
                        </Link>
                      </TableCell>
                      <TableCell>
                        <Link
                            href={`/drivers/${driver.driver_number}`}
                            className="hover:text-red-600 dark:hover:text-red-400 transition-colors"
                        >
                          <div className="cursor-pointer">
                            <div className="font-medium text-lg">{driver.full_name}</div>
                            <div className="text-base text-slate-500 dark:text-slate-400">{driver.abbreviation}</div>
                          </div>
                        </Link>
                      </TableCell>
                      <TableCell>
                        <Link href={`/drivers/${driver.driver_number}`} className="hover:opacity-80 transition-opacity">
                      <span style={{ color: driver.color }} className="font-medium cursor-pointer text-lg">
                        {driver.team_name}
                      </span>
                        </Link>
                      </TableCell>
                      <TableCell className="text-right">{driver.lap_number}</TableCell>
                      <TableCell className="text-right">
                        {driver.interval > 0 ? `+${driver.interval.toFixed(3)}` : driver.interval.toFixed(3)}
                      </TableCell>
                      <TableCell className="text-right">
                        <div>
                          <div>{driver.pitsops}</div>
                          {driver.last_pit_lap > 0 && (
                              <div className="text-sm text-slate-500 dark:text-slate-400">Круг {driver.last_pit_lap}</div>
                          )}
                        </div>
                      </TableCell>
                    </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </div>
      </div>
  )
}
