import { useState } from "react";
import "./App.css";
function App() {
  const [msgVal, setMsgVal] = useState<string>('');
  const [fromMsg, setFromMsg] = useState<string>('');
  const [toMsg, setToMsg] = useState<string>('');
  const [messagePane, setMessagePane] = useState<any[]>([]);
  const [ready, setReady] = useState<boolean>(false);
  const [ws, setWs] = useState<WebSocket>();

  const sendMessage = () => {
    if (ready && ws){
      ws.send(msgVal);
    }
    setMessagePane((prevMessagePane) => [...prevMessagePane, msgVal]);
    setMsgVal("");
  };

  const handleWsStart = () => {
    var from = fromMsg
    var to = toMsg
    var ws = new WebSocket(
      "ws://localhost:5000/ws?from=" +
        encodeURIComponent(from) +
        "&to=" +
        encodeURIComponent(to)
    );
    setWs(ws);
    ws.onopen = () => {
      setMessagePane((prevM) => [...prevM, `Connected`])
      setReady(true);
    };
    ws.onmessage = (e) => {
      setMessagePane((prevMessagePane) => [...prevMessagePane, e.data]);
    };
    ws.onclose = () => {
      setReady(false);
      setMessagePane((prevMessagePane) => [
        ...prevMessagePane,
        "Connection closed",
      ]);
    };
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
          <button type="submit" onClick={handleWsStart}>
            Start
          </button>
        </div>
      </div>

      <div
        style={{
          border: "1px solid black",
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
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              sendMessage();
            }
          }}
          placeholder="Message"
        />
        <button type="submit" onClick={sendMessage}>
          Send
        </button>
      </div>
    </div>
  );
}

export default App;
