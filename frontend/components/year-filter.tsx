"use client"

import { Button } from "@/components/ui/button"
import { useRouter, useSearchParams } from "next/navigation"

interface YearFilterProps {
  availableYears: number[]
  selectedYear: number | null
}

export function YearFilter({ availableYears, selectedYear }: YearFilterProps) {
  const router = useRouter()
  const searchParams = useSearchParams()

  const handleYearChange = (year: number | null) => {
    const params = new URLSearchParams(searchParams.toString())

    if (year) {
      params.set("year", year.toString())
    } else {
      params.delete("year")
    }

    router.push(`/?${params.toString()}`)
  }

  return (
      <div className="flex flex-wrap gap-2">
        <Button
            variant={selectedYear === null ? "default" : "outline"}
            size="default"
            className="text-base"
            onClick={() => handleYearChange(null)}
        >
          Все годы
        </Button>

        {availableYears.map((year) => (
            <Button
                key={year}
                variant={selectedYear === year ? "default" : "outline"}
                size="default"
                className="text-base"
                onClick={() => handleYearChange(year)}
            >
              {year}
            </Button>
        ))}
      </div>
  )
}
