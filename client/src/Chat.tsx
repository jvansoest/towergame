import { useEffect, useState } from "react";
import socket from "./socket";

const Chat = () => {
  const [input, setInput] = useState("");
  const [messages, setMessages] = useState<string[]>([]);

  useEffect(() => {
    const onChatUpdate = (data: Record<string, unknown>) => {
      setMessages(data.chat as string[]);
    };
    socket.on("chatUpdate", onChatUpdate);
    return () => {
      socket.off("chatUpdate", onChatUpdate);
    };
  }, []);

  const handleSubmit = () => {
    if (input.trim().length) {
      socket.emit("message", { text: input });
    }
    setInput("");
  };

  return (
    <div className="chat">
      <div className="chat__log">
        {messages.map((msg, i) => (
          <p className="chat__msg" key={i}>
            {msg}
          </p>
        ))}
      </div>

      <form
        className="chat__form"
        onSubmit={(e) => {
          e.preventDefault();
          handleSubmit();
        }}
      >
        <input
          className="chat__input"
          type="text"
          name="message"
          value={input}
          onChange={(e) => setInput(e.target.value)}
        />
        <button className="chat__send" type="submit">
          GO
        </button>
      </form>
    </div>
  );
};

export default Chat;
