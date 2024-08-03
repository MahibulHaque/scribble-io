const io = require('socket.io-client')

const SERVER =
  process.env.NODE_ENV === 'production'
    ? 'https://scribble-d2zqkosdka-uc.a.run.app'
    : 'http://localhost:3001'

// export const socket = io(SERVER, {
//   transports: ['websocket'],
// })
export const socket = io(SERVER, { transports: ['websocket'] })
