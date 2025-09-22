import logo from './logotype.png';
import './App.css';
import { useState } from 'react';




function App() {
  
  const [message, setMessage] = useState("");

  function logtoconsole(){
    console.log("\n\n\nSomebody wants to join the book club!\n\n\n")
    fetch("/api/").then(res => res.text()).then(data => {
      console.log(data);
      setMessage(data);
    })
  }


  const join_prompt = "Join The Infinite Library!"
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <button className="join_button" type="button" onClick={logtoconsole}>
          {join_prompt}
        </button>
        <br/>
        <br/>
        {message && (
          <p className="response">
            {message}
          </p>
        )

        }
      </header>
    </div>
  );
}

export default App;