"use client"

import { Button } from "@/components/ui/button"
import { ArrowUp, ArrowDown } from "lucide-react"
import { useRouter, useSearchParams } from "next/navigation"

interface SortFilterProps {
    currentSort: string
}

export function SortFilter({ currentSort }: SortFilterProps) {
    const router = useRouter()
    const searchParams = useSearchParams()

    const handleSortChange = (sortOrder: string) => {
        const params = new URLSearchParams(searchParams.toString())

        if (sortOrder === "desc") {
            // Default sort, remove parameter
            params.delete("sort")
        } else {
            params.set("sort", sortOrder)
        }

        router.push(`/?${params.toString()}`)
    }

    return (
        <div className="flex gap-2">
            <Button
                variant={currentSort === "desc" ? "default" : "outline"}
                size="default"
                onClick={() => handleSortChange("desc")}
                className="flex items-center gap-2 text-base"
            >
                <ArrowDown className="w-4 h-4" />
                Новые первыми
            </Button>

            <Button
                variant={currentSort === "asc" ? "default" : "outline"}
                size="default"
                onClick={() => handleSortChange("asc")}
                className="flex items-center gap-2 text-base"
            >
                <ArrowUp className="w-4 h-4" />
                Старые первыми
            </Button>
        </div>
    )
}
