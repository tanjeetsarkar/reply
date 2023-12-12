import { useState } from "react";
import "./App.css";
function App() {
  const [msgVal, setMsgVal] = useState<any>();
  const [fromMsg, setFromMsg] = useState<any>();
  const [toMsg, setToMsg] = useState<any>();
  const [messagePane, setMessagePane] = useState<any[]>([]);

  const sendMessage = () => {
    setMessagePane((prevMessagePane) => [...prevMessagePane, msgVal]);
    setMsgVal("");
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
