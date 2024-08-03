import { nanoid } from 'nanoid'
import ThemeMenuButton from '@/components/ThemeMenuButton'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Separator } from '@radix-ui/react-dropdown-menu'
import JoinRoomButtoon from '@/components/JoinRoomButton'
import CreateRoomForm from '@/components/CreateRoomForm'

export const dynamic = 'force-dynamic'

export default function Home() {
  const roomId = nanoid()

  return (
    <main className='flex h-screen flex-col items-center justify-between pb-5 pt-[13vh]'>
      <ThemeMenuButton className='fixed right-[5vw] top-5 flex-1 md:right-5' />
      <Card className='w-[90vw] max-w-[400px]'>
        <CardHeader>
          <CardTitle>Scribble</CardTitle>
          <CardDescription>
            Draw on the same canvas with your friends in real-time.
          </CardDescription>
        </CardHeader>
        <CardContent className='flex flex-col space-y-4'>
          <CreateRoomForm roomId={roomId} />
          <div className='flex items-center space-x-2'>
            <Separator />
            <span className='text-xs text-muted-foreground text-center w-full'>OR</span>
            <Separator />
          </div>
          <JoinRoomButtoon />
        </CardContent>
      </Card>
    </main>
  )
}
