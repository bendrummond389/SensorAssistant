import Image from 'next/image'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import SensorDataComponent from '@/components/websockets/SensorData'
import { ModeToggle } from '@/components/theme/ThemeToggle'

export default function Home() {
  return (
    <main className="">
      <div className="m-5">
        <ModeToggle />
      </div>
      <div className="ml-5">
        <SensorDataComponent />
      </div>
    </main>
  )
}
