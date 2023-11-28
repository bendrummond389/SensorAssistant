'use client'

import React from 'react'
import { SensorData } from '@/types/SensorData'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

const SensorDataComponent: React.FC = () => {
  const [sensorData, setSensorData] = React.useState<SensorData[]>([])
  const [ws, setWs] = React.useState<WebSocket | null>(null)
  const goWebsocketServerUrl = process.env.NEXT_PUBLIC_WEBSOCKET_SERVER_URL

  React.useEffect(() => {
    const webSocket = new WebSocket('ws://localhost:8080/ws')
    setWs(webSocket)

    webSocket.onmessage = (event) => {
      const data: SensorData[] = JSON.parse(event.data)
      console.log('Data received:', data)
      setSensorData(data)
    }

    webSocket.onerror = (error: Event) => {
      console.error('WebSocket Error:', error)
    }

    return () => {
      webSocket.close()
    }
  }, [])

  return (
    <div>
      {sensorData && sensorData.length > 0 ? (
        sensorData.map((data, index) => (
          <div key={index}>
            <Card className="max-w-[500px]">
              <CardHeader>
                <CardTitle>{data.sensor_name}</CardTitle>
                <CardDescription>Id: {data.sensor_id}</CardDescription>
              </CardHeader>
              <CardContent>
                <p>
                  {data.value}: {data.units}
                </p>
              </CardContent>
            </Card>
          </div>
        ))
      ) : (
        <div>No sensor data available</div>
      )}
    </div>
  )
}

export default SensorDataComponent
