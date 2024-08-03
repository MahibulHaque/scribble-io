'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'

import type { DrawOptions } from '@/types'
import { useCanvasStore } from '@/store/canvasStore'
import { useUserStore } from '@/store/userStore'
import { socket } from '@/lib/socket'
import { draw, drawWithDataURL } from '@/lib/utils'
import useDraw, { type DrawProps } from '@/hooks/useDraw'
import { Skeleton } from '@/components/ui/skeleton'
import UndoButton from '@/components/UndoButton'
import ClearButton from '@/components/ClearButton'

export default function DrawingCanvas() {
  const router = useRouter()
  const { roomId } = useParams()

  const containerRef = useRef<HTMLDivElement>(null)

  const [isCanvasLoading, setIsCanvasLoading] = useState(true)

  const strokeColor = useCanvasStore(state => state.strokeColor)
  const strokeWidth = useCanvasStore(state => state.strokeWidth)
  const dashGap = useCanvasStore(state => state.dashGap)
  const user = useUserStore(state => state.user)

  useEffect(() => {
    if (!user) {
      router.replace('/')
    }
  }, [router, user])

  const onDraw = useCallback(
    ({ ctx, currentPoint, prevPoint }: DrawProps) => {
      const drawOptions = {
        ctx,
        currentPoint,
        prevPoint,
        strokeColor,
        strokeWidth,
        dashGap,
      }
      draw(drawOptions)
      socket.emit('draw',  drawOptions, roomId )
    },
    [strokeColor, strokeWidth, dashGap, roomId]
  )

  const { canvasRef, onInteractStart, clear, undo } = useDraw(onDraw)

  useEffect(() => {
    const canvasElement = canvasRef.current
    const ctx = canvasElement?.getContext('2d')

    socket.emit('client-ready', roomId)

    socket.on('client-loaded', () => {
      setIsCanvasLoading(false)
    })

    socket.on('get-canvas-state', () => {
      const canvasState = canvasRef.current?.toDataURL()
      if (!canvasState) return

      socket.emit('send-canvas-state', JSON.stringify({ canvasState, roomId }))
    })

    socket.on('canvas-state-from-server', (canvasState: string) => {
      if (!ctx || !canvasElement) return

      drawWithDataURL(canvasState, ctx, canvasElement)
      setIsCanvasLoading(false)
    })

    socket.on('update-canvas-state', (drawOptions: DrawOptions) => {
      // const { drawOptions } = JSON.parse(data)
      console.log('Draw options', drawOptions)
      if (!ctx) return
      draw({ ...drawOptions, ctx })
    })

    socket.on('undo-canvas', (canvasState: string) => {
      if (!ctx || !canvasElement) return

      drawWithDataURL(canvasState, ctx, canvasElement)
    })

    return () => {
      socket.off('get-canvas-state')
      socket.off('canvas-state-from-server')
      socket.off('update-canvas-state')
      socket.off('undo-canvas')
    }
  }, [canvasRef, roomId])

  useEffect(() => {
    const setCanvasDimensions = () => {
      if (!containerRef.current || !canvasRef.current) return

      const { width, height } = containerRef.current?.getBoundingClientRect()

      canvasRef.current.width = width - 50
      canvasRef.current.height = height - 50
    }

    setCanvasDimensions()
  }, [canvasRef])

  const handleInteractStart = () => {
    const canvasElement = canvasRef.current
    if (!canvasElement) return

    socket.emit(
      'add-undo-point',
      JSON.stringify({
        roomId,
        undoPoint: canvasElement.toDataURL(),
      })
    )
    onInteractStart()
  }

  return (
    <div
      ref={containerRef}
      className='relative flex items-center justify-center w-full h-full'
    >
      {!isCanvasLoading && (
        <div className='absolute right-[25px] top-[25px] flex select-none rounded-none rounded-bl rounded-tr-[2.5px]'>
          <UndoButton undo={undo} />

          <ClearButton canvasRef={canvasRef} clear={clear} />
        </div>
      )}

      {isCanvasLoading && (
        <Skeleton className='absolute h-[calc(100%-50px)] w-[calc(100%-50px)]' />
      )}

      <canvas
        id='canvas'
        ref={canvasRef}
        onMouseDown={handleInteractStart}
        onTouchStart={handleInteractStart}
        width={0}
        height={0}
        className='bg-white border rounded touch-none'
      />
    </div>
  )
}
