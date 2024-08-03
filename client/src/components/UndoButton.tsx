'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import { useHotkeys } from 'react-hotkeys-hook'
import { Loader2 } from 'lucide-react'

import { socket } from '@/lib/socket'
import { cn, isMacOS } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { Kbd } from '@/components/ui/kbd'

interface UndoButtonProps {
  undo: (lastUndoPoint: string) => void
}

export default function UndoButton({ undo }: UndoButtonProps) {
  const { roomId } = useParams()

  const [isLoading, setIsLoading] = useState(false)

  const isMac = isMacOS()
  const hotKey = `${isMac ? 'Meta' : 'ctrl'} + z`

  const undoCanvas = () => {
    setIsLoading(true)
    socket.emit('get-last-undo-point', roomId)
  }

  useHotkeys(hotKey, undoCanvas)

  useEffect(() => {
    // This socket does undo function
    socket.on('last-undo-point-from-server', (lastUndoPoint: string) => {
      undo(lastUndoPoint)
      socket.emit(
        'undo',
        JSON.stringify({
          canvasState: lastUndoPoint,
          roomId,
        })
      )

      socket.emit('delete-last-undo-point', roomId)
      setIsLoading(false)

      return () => {
        socket.off('last-undo-point-from-server')
      }
    })
  }, [roomId, undo])

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            variant='outline'
            className='w-16 p-0 border-0 border-b border-l rounded-none rounded-bl-md focus-within:z-10'
            disabled={isLoading}
            onClick={undoCanvas}
          >
            {isLoading ? <Loader2 className='w-4 h-4 animate-spin' /> : 'Undo'}
          </Button>
        </TooltipTrigger>

        <TooltipContent className='flex gap-1'>
          <Kbd className={cn({ 'text-xs': isMac })}>{isMac ? '⌘' : 'Ctrl'}</Kbd>
          <Kbd>Z</Kbd>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  )
}
