import { Card, CardContent, CardHeader } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"

export default function Loading() {
  return (
    <div className="container mx-auto py-10">
      <Card className="max-w-3xl mx-auto">
        <CardHeader className="bg-slate-50 border-b">
          <div className="flex items-center justify-between">
            <div className="space-y-2">
              <Skeleton className="h-8 w-64" />
              <Skeleton className="h-4 w-48" />
            </div>
            <Skeleton className="w-16 h-16 rounded-full" />
          </div>
        </CardHeader>
        <CardContent className="pt-6 space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-4">
              <Skeleton className="h-5 w-40" />
              <Skeleton className="h-5 w-56" />
              <Skeleton className="h-5 w-48" />
            </div>
            <div>
              <Skeleton className="h-5 w-24 mb-2" />
              <Skeleton className="h-20 w-full mb-4" />
              <Skeleton className="h-4 w-20 mb-2" />
              <Skeleton className="h-16 w-full" />
              <Skeleton className="h-3 w-32 mt-1" />
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
