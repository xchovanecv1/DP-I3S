import { LOGIN_SUCCESSFUL } from "../actions/auth"

const initialState = {
  token: null,
}

export default (state = initialState, action) => {
    switch (action.type) {
     case LOGIN_SUCCESSFUL:
      return {
        ...state,
        token: action.token,
      }
     default:
      return state
    }
   }