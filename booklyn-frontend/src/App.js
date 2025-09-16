import logo from './latest_logo.png';
import './App.css';

function App() {

  const join_prompt = "Join the Booklyn community!"
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <button className="join_button" type="button">
          {join_prompt}
        </button>
      </header>
    </div>
  );
}

export default App;
