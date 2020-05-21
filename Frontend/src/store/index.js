/*
import { createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';
import rootReducer from './reducers/rootReducer';


export default function configureStore() {
 return createStore(
  rootReducer,
  applyMiddleware(thunk)
 );
}
*/

import { createStore } from 'easy-peasy';
import { model as AuthModel} from "./actions/auth"
import { model as ImportModel} from "./actions/import"
import { model as AccountsModel} from "./actions/account"
import { model as PropertiesModel} from "./actions/props"
import { model as RoomsModel} from "./actions/rooms"
import { model as CardsModel} from "./actions/cards"

export default function configureStore() {
    const store = createStore(
        {
            auth: AuthModel,
            import: ImportModel,
            accounts: AccountsModel,
            props: PropertiesModel,
            rooms: RoomsModel,
            cards: CardsModel,
        }
    ); 
    return store; 

}