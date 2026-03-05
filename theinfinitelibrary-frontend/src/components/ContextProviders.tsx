import { createContext, ReactNode, useContext, useState } from "react"
import * as Types from '../types/types';

//User Context

export const UserCtx = createContext< Types.UserContext | null>(null)

export function UserContextProvider({children}: {children: ReactNode}){
  const [user, setUser] = useState<Types.User>({username: "", password: ""})
  return(
    <UserCtx.Provider value ={{user, setUser}}>
      {children}
    </UserCtx.Provider>
  )
}

export function useUserContext(){
  const ctx =  useContext(UserCtx)
  if(!ctx) throw new Error("useUserContext must be used within <UserContextProvider>")
  return ctx
}