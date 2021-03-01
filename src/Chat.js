import React, { useState, useEffect } from "react";

/*
  NETWORK
*/
import socket from "./socket.js";

const Chat = () => {
  const [input, setInput] = useState("");
  const [messages, setMessages] = useState([]);

  const setter = (data) => {
    setMessages(data.chat);
  };

  useEffect(() => {
    socket.on("chatUpdate", setter);
  });

  const handleSubmit = (e) => {
    socket.emit("message", input);
  };

  return (
    <>
      <input
        type="text"
        name="name"
        onChange={(e) => setInput(e.target.value)}
      />
      <button onClick={handleSubmit}>GO</button>
      {messages.map((msg) => (
        <p>{msg}</p>
      ))}
    </>
  );
};
export default Chat;
