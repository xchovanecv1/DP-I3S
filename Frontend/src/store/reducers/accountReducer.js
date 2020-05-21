import { ALL_ACCOUNTS_LOADED, ALL_ACCOUNTS_PENDING, ALL_ACCOUNTS_FAILED } from "../actions/account"

const initialState = {
  accountsPending: true,
  accounts: null,
  accountsError: null,
}

export default (state = initialState, action) => {
    switch (action.type) {
     case ALL_ACCOUNTS_LOADED:
      return {
        ...state,
        ...action,
      }
    case ALL_ACCOUNTS_PENDING:
      return {
        ...state,
        ...action,
      }
    case ALL_ACCOUNTS_FAILED:
      return {
        ...state,
        ...action,
      }
     default:
      return state
    }
   }