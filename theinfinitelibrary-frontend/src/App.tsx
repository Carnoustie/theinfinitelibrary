import './style/App.css';
import { useState } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import {Home, Login, Signup, Loggedin, AddBook, SecurityInfo, ChatRoom} from './components/Components'
import * as Types from './types/types';
import { UserContextProvider } from './components/ContextProviders';

function App() {
  const [username, setUsername] = useState("")
  const [bookList, setBookList] = useState<Types.Book[]>([])
  const [previousSite, setPreviousSite] = useState("")
  const [chatRooms, setChatRooms] = useState<Types.ChatroomID[]>([])

  const join_prompt = "Join The Infinite Library!"
  return (
    <UserContextProvider>
      <BrowserRouter>
      <Routes>
        <Route path="/"
        element={
          <Home
            join_prompt={join_prompt}
          />
          }
        />
        <Route 
        path="/signup"
        element= {<Signup previousSite = {previousSite} setPreviousSite = {setPreviousSite}/>}
        />
        <Route 
        path="/login"
        element= {
            <Login
              bookList = {bookList}
              setBookList = {setBookList}
              previousSite = {previousSite}
              setPreviousSite = {setPreviousSite}
              chatrooms = {chatRooms}
              setChatrooms = {setChatRooms}
              />
        }
        />
        <Route
        path="/loggedin"
        element=  {<Loggedin username = {username} setUsername= {setUsername} booklist = {bookList} setBookList = {setBookList}/>}
        />
        <Route
        path="/addbook"
        element= {<AddBook username = {username} setUsername = {setUsername} booklist = {bookList} setBookList = {setBookList}/>}
        />
        <Route
        path="/securityinfo"
        element = {<SecurityInfo previousSite = {previousSite} setPreviousSite = {setPreviousSite} />}
        />
        <Route
        path="/chatroom/:chatId"
        element = {
            <ChatRoom username={username} setUsername={setUsername}/>
        }
        />
      </Routes>
      </BrowserRouter>
    </UserContextProvider>

  );
}

export default App;