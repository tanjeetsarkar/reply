import { useEffect, useState } from "react";
import "./App.css";
import { Start_client, SendMessage,GetContacts,AddContact } from "../wailsjs/go/main/App.js";
import { EventsOn } from "../wailsjs/runtime";

function App() {
  const [msgVal, setMsgVal] = useState<string>("");
  const [fromMsg, setFromMsg] = useState<string>("");
  const [toMsg, setToMsg] = useState<string>("");
  const [messagePane, setMessagePane] = useState<any[]>([]);
  const [contactList, setContactList] = useState<any[]>([]);

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

  useEffect(() => {
    // const contactsObj = [
    //   {
    //     name: "sum",
    //     uid : "14443"
    //   },
    //   {
    //     name: "tan",
    //     uid: "13232"
    //   },
    //   {
    //     name: "cat",
    //     uid: "12321"
    //   },
    //   {
    //     name: "dog",
    //     uid: "12356"
    //   },
    //   {
    //     name: "bird",
    //     uid: "12345"
    //   },
    //   {
    //     name: "fish",
    //     uid: "12346"
    //   },
    //   {
    //     name: "lion",
    //     uid: "12347"
    //   },
    //   {
    //     name: "tiger",
    //     uid: "12348"
    //   },
    //   {
    //     name: "elephant",
    //     uid: "12349"
    //   },
    //   {
    //     name: "snake",
    //     uid: "12350"
    //   },
    //   {
    //     name: "rabbit",
    //     uid: "12351"
    //   },
    //   {
    //     name: "monkey",
    //     uid: "12352"
    //   },
    //   {
    //     name: "horse",
    //     uid: "12353"
    //   },
    //   {
    //     name: "cow",
    //     uid: "12354"
    //   },
    //   {
    //     name: "sheep",
    //     uid: "12355"
    //   },
    //   {
    //     name: "goat",
    //     uid: "12356"
    //   },
    //   {
    //     name: "chicken",
    //     uid: "12357"
    //   },
    //   {
    //     name: "pig",
    //     uid: "12358"
    //   },
    //   {
    //     name: "duck",
    //     uid: "12359"
    //   },
    //   {
    //     name: "goose",
    //     uid: "12360"
    //   },
    //   {
    //     name: "turkey",
    //     uid: "12361"
    //   },
    //   {
    //     name: "donkey",
    //     uid: "12362"
    //   },
    //   {
    //     name: "camel",
    //     uid: "12363"
    //   },
    //   {
    //     name: "kangaroo",
    //     uid: "12364"
    //   },
    //   {
    //     name: "penguin",
    //     uid: "12365"
    //   },
    //   {
    //     name: "bear",
    //     uid: "12366"
    //   },
    //   {
    //     name: "fox",
    //     uid: "12367"
    //   },
    //   {
    //     name: "wolf",
    //     uid: "12368"
    //   },
    //   {
    //     name: "whale",
    //     uid: "12369"
    //   },
    //   {
    //     name: "dolphin",
    //     uid: "12370"
    //   },
    //   {
    //     name: "shark",
    //     uid: "12371"
    //   },
    //   {
    //     name: "crocodile",
    //     uid: "12372"
    //   },
    // ]
    GetContacts().then((res) => {
      console.log("response",res)
      setContactList(res);
    })
    // check if elements in contctsObj is in contactList, if not add it to contactList
    // for (let i = 0; i < contactsObj.length; i++) { 
    //   if (!contactList.includes(contactsObj[i])) {
    //     AddContact(contactsObj[i].name, contactsObj[i].uid);
    //     setContactList((prev) => [...prev, contactsObj[i]]);
    //   }
    // }
  },[])

  const SendChat = () => {
    setMessagePane((prev) => [...prev, msgVal]);
    SendMessage(msgVal, fromMsg, toMsg);
    setMsgVal("");
  };

  const startClient = () => {
    Start_client(fromMsg);
  };

  return (
    <div className="container">
      <div className="head-container">

        <div className="input-container-from-to">
          <input
            type="text"
            value={fromMsg}
            onChange={(e) => setFromMsg(e.target.value)}
            placeholder="From"
          />
          <select
            style={{

            }}
            name="to"
            id="to"
            onChange={(e) => setToMsg(e.target.value)}
            value={toMsg}
          >
            <option value="">To</option>
            {contactList && contactList.map((contact) => (
              <option value={contact.name}>{contact.name}</option>
            ))}
          </select>
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
