import React, { useState, useEffect } from "react";
import { Col, Row, Input, Button, Card, InputGroup } from "reactstrap";
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
    input.trim().length && socket.emit("message", input);
    setInput("");
  };

  return (
    <>
      <Row>
        <Col xs="6">
          <Card style={{ height: "300px", overflow: "auto" }}>
            {messages.map((msg) => (
              <p>{msg}</p>
            ))}
          </Card>
        </Col>
      </Row>

      <Row>
        <Col xs="6">
          <InputGroup>
            <Input
              style={{
                borderTopRightRadius: "0px",
                borderBottomRightRadius: "0px",
              }}
              className="input-sm"
              type="text"
              name="name"
              value={input}
              onChange={(e) => setInput(e.target.value)}
            />

            <Button
              style={{
                borderTopLeftRadius: "0px",
                borderBottomLeftRadius: "0px",
              }}
              onClick={handleSubmit}
            >
              GO
            </Button>
          </InputGroup>
        </Col>
      </Row>
    </>
  );
};
export default Chat;
