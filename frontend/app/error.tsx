"use client"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { AlertCircle } from "lucide-react"
import { useEffect } from "react"

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  useEffect(() => {
    console.error(error)
  }, [error])

  return (
    <div className="container mx-auto py-10">
      <Card className="max-w-md mx-auto">
        <CardHeader className="bg-red-50 text-red-700">
          <CardTitle className="flex items-center gap-2">
            <AlertCircle className="h-5 w-5" />
            Ошибка загрузки данных
          </CardTitle>
        </CardHeader>
        <CardContent className="pt-6">
          <p>Не удалось загрузить список гонок. Пожалуйста, попробуйте снова позже.</p>
        </CardContent>
        <CardFooter>
          <Button onClick={reset}>Попробовать снова</Button>
        </CardFooter>
      </Card>
    </div>
  )
}
