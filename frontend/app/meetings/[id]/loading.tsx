import { Card, CardContent, CardHeader } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"

export default function Loading() {
  return (
    <div className="container mx-auto py-10">
      <Card className="max-w-4xl mx-auto">
        <CardHeader className="bg-gradient-to-r from-red-50 to-red-100 border-b">
          <div className="flex items-start justify-between">
            <div className="space-y-2 flex-1">
              <div className="flex items-center gap-3">
                <Skeleton className="h-8 w-64" />
                <Skeleton className="h-6 w-16" />
              </div>
              <Skeleton className="h-5 w-80" />
              <Skeleton className="h-6 w-40" />
            </div>
            <Skeleton className="w-16 h-16 rounded-full" />
          </div>
        </CardHeader>

        <CardContent className="pt-6 space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-4">
              <Skeleton className="h-5 w-48" />
              <Skeleton className="h-5 w-56" />
              <Skeleton className="h-5 w-64" />
              <Skeleton className="h-5 w-60" />
            </div>

            <div className="space-y-4">
              <div>
                <Skeleton className="h-6 w-32 mb-2" />
                <Skeleton className="h-4 w-20" />
              </div>
              <Skeleton className="h-32 w-full rounded-lg" />
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
