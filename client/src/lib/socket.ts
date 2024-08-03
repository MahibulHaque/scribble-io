const io = require('socket.io-client')

const SERVER =
  process.env.NODE_ENV === 'production'
    ? process.env.SOCKET_SERVER_URL
    : 'http://localhost:3001'

// export const socket = io(SERVER, {
//   transports: ['websocket'],
// })
export const socket = io(SERVER, { transports: ['websocket'] })
