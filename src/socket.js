import socketIOClient from "socket.io-client";
const socket = socketIOClient("http://localhost:7777");
export default socket;
