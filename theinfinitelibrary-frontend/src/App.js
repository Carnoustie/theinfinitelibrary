import logo from './latest_logo.png';
import './App.css';

function logtoconsole(){
  console.log("\n\n\nSomebody wants to join the book club!\n\n\n")
  fetch("http://localhost:8080/").then(res => res.text()).then(data => console.log(data))
}

function App() {

  const join_prompt = "Join the Booklyn community!"
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <button className="join_button" type="button" onClick={logtoconsole}>
          {join_prompt}
        </button>
      </header>
    </div>
  );
}

export default App;
