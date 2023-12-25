import { useEffect, useState } from "react";
import "./App.css";
import { Start_client, SendMessage } from "../wailsjs/go/main/App.js";
import { EventsOn } from "../wailsjs/runtime";

function App() {
  const [msgVal, setMsgVal] = useState<string>("");
  const [fromMsg, setFromMsg] = useState<string>("");
  const [toMsg, setToMsg] = useState<string>("");
  const [messagePane, setMessagePane] = useState<any[]>([]);

  useEffect(() => {
    EventsOn("clientStarted", (message: string) => {
      if (message && message.length > 0) {
        setMessagePane((prevMessagePane) => [...prevMessagePane, message]);
      }
    });
    EventsOn("recieveMessage", (message) => {
      if (message && message.length > 0) {
        setMessagePane((prevMessagePane) => [...prevMessagePane, message]);
      }
    });
  }, []);

  const SendChat = () => {
    setMessagePane((prev) => [...prev, msgVal]);
    SendMessage(msgVal);
    setMsgVal("");
  };

  const startClient = () => {
    Start_client(fromMsg, toMsg);
  };

  return (
    <div className="container">
      <div className="head-container">
        <h1>Reply</h1>

        <div className="input-container-from-to">
          <input
            type="text"
            value={fromMsg}
            onChange={(e) => setFromMsg(e.target.value)}
            placeholder="From"
          />
          <input
            type="text"
            value={toMsg}
            onChange={(e) => setToMsg(e.target.value)}
            placeholder="To"
          />
          <button type="submit" onClick={startClient}>
            Start
          </button>
        </div>
      </div>

      <div
        style={{
          minWidth: "400px",
          minHeight: "200px",
          maxHeight: "200px",
          overflow: "auto",
          display: "flex",
          flexDirection: "column-reverse",
        }}
      >
        <ul
          style={{
            listStyleType: "none",
            textAlign: "left", // Add this line to align list items to the left
          }}
        >
          {messagePane.map((item, index) => (
            <li key={index}>{item}</li>
          ))}
        </ul>
      </div>
      <div className="input-container-bottom">
        <input
          type="text"
          value={msgVal}
          onChange={(e) => setMsgVal(e.target.value)}
          placeholder="Message"
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              SendChat();
            }
          }}
        />
        <button type="submit" onClick={SendChat}>
          Send
        </button>
      </div>
    </div>
  );
}

export default App;
