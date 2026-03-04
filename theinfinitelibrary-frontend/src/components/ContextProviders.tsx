import { createContext, ReactNode } from "react"
import * as Types from '../types/types';

//User Context

export const UserContext = createContext< Types.User| null>(null)

export function UserContextProvider({children, userValue}: {children: ReactNode, userValue: Types.User}){
  return(
    <UserContext.Provider value ={userValue}>
      {children}
    </UserContext.Provider>
  )
}