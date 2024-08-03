'use client'

import { useEffect } from 'react'

import { useMembersStore } from '@/store/membersStore'
import { socket } from '@/lib/socket'
import { ScrollArea } from '@/components/ui/scroll-area'
import { useToast } from '@/hooks/useToast'

export default function MemberList() {
  const { toast } = useToast()
  const [members, setMembers] = useMembersStore(state => [
    state.members,
    state.setMembers,
  ])

  useEffect(() => {
    socket.on('update-members', (data: string) => {
      const members = JSON.parse(data)
      if (members) {
        setMembers(members)
      } else {
        setMembers([])
      }
    })

    socket.on('send-notification', (data: string) => {
      const { title, message } = JSON.parse(data)
      toast({
        title: title,
        description: message,
      })
    })

    return () => {
      socket.off('update-members')
      socket.off('send-notification')
    }
  }, [setMembers, toast])

  return (
    <div className='my-6 select-none'>
      <h2 className='pb-2.5 font-medium'>Members</h2>

      <ScrollArea className='h-48'>
        <ul className='flex flex-col gap-1 rounded-md px-1'>
          {members.map(({ id, username }) => (
            <li key={id}>{username}</li>
          ))}
        </ul>
      </ScrollArea>
    </div>
  )
}
